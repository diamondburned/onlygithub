package create

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"image"
	"mime"
	"mime/multipart"
	"net/http"
	"strconv"
	"strings"

	_ "embed"
	"image/jpeg"
	_ "image/jpeg"
	_ "image/png"

	"github.com/diamondburned/hrt"
	"github.com/go-chi/chi/v5"
	"github.com/hashicorp/go-hclog"
	"github.com/pkg/errors"
	"golang.org/x/image/draw"
	"libdb.so/onlygithub"
	"libdb.so/onlygithub/frontend"
	"libdb.so/onlygithub/frontend/layouts"
)

const maxMemory = 2 << 20 // 2 MiB

// Services is a collection of services used by the create page.
type Services struct {
	Images onlygithub.ImageService
	Posts  onlygithub.PostService
	Tiers  onlygithub.TierService
}

func New(services Services) http.Handler {
	r := chi.NewRouter()
	r.Use(frontend.OwnerOnly)

	r.Get("/", services.get)
	r.Post("/", services.post)

	return r
}

func (h *Services) get(w http.ResponseWriter, r *http.Request) {
	site := frontend.SiteFromRequest(r)
	owner := frontend.OwnerFromRequest(r)

	tiers, err := h.Tiers.Tiers(r.Context())
	if err != nil {
		layouts.RenderError(w, r, hrt.WrapHTTPError(500, err))
		return
	}

	create(r, site, owner, createData{
		Tiers: tiers,
	}).Render(r.Context(), w)
}

func (h *Services) post(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(maxMemory); err != nil {
		layouts.RenderError(w, r, hrt.WrapHTTPError(http.StatusBadRequest, err))
		return
	}

	minCost, err := strconv.Atoi(r.FormValue("minimum-cost"))
	if err != nil {
		layouts.RenderError(w, r, hrt.WrapHTTPError(
			http.StatusBadRequest, errors.Wrap(err, "mininum cost is not a valid number")))
		return
	}

	visibility := onlygithub.Visibility(r.FormValue("visibility"))
	if err := visibility.Validate(); err != nil {
		layouts.RenderError(w, r, hrt.WrapHTTPError(
			http.StatusBadRequest, errors.Wrap(err, "invalid visibility")))
		return
	}

	content := r.FormValue("content")
	imageFiles := r.MultipartForm.File["images"]
	if content == "" && len(imageFiles) == 0 {
		layouts.RenderError(w, r, hrt.WrapHTTPError(
			http.StatusBadRequest, errors.New("content or images must be provided")))
		return
	}

	imageIDs, err := h.uploadImages(
		r.Context(),
		visibility, onlygithub.Cents(minCost), imageFiles)
	if err != nil {
		layouts.RenderError(w, r, err)
		return
	}

	create := onlygithub.CreatePostRequest{
		Visibility:  visibility,
		MinimumCost: onlygithub.Cents(minCost),
		Markdown:    content,
		AssetIDs:    imageIDs,
	}
	if r.FormValue("comments") == "" {
		create.AllowComments = new(bool)
	}
	if r.FormValue("reactions") == "" {
		create.AllowReactions = new(bool)
	}

	post, err := h.Posts.CreatePost(r.Context(), create)
	if err != nil {
		layouts.RenderError(w, r, hrt.WrapHTTPError(500, err))
		return
	}

	postURL := fmt.Sprintf("/posts/%s", post.ID)
	http.Redirect(w, r, postURL, http.StatusFound)
}

func (h *Services) uploadImages(
	ctx context.Context, visibility onlygithub.Visibility, minimumCost onlygithub.Cents,
	headers []*multipart.FileHeader,
) ([]onlygithub.ID, error) {
	var assetIDs []onlygithub.ID

	fail := true
	defer func() {
		if !fail {
			return
		}

		log := hclog.FromContext(ctx)
		for _, id := range assetIDs {
			if err := h.Images.DeleteImage(ctx, id); err != nil {
				log.Warn("failed to delete image after failure", "id", id, "err", err)
			} else {
				log.Info("deleted image after failure", "id", id)
			}
		}
	}()

	for _, header := range headers {
		if header.Filename == "" {
			return nil, hrt.NewHTTPError(400, "image is missing filename")
		}

		assetID, err := h.uploadImage(ctx, visibility, minimumCost, header)
		if err != nil {
			return nil, err
		}

		assetIDs = append(assetIDs, assetID)
	}

	fail = false
	return assetIDs, nil
}

func (h *Services) uploadImage(
	ctx context.Context, visibility onlygithub.Visibility, minimumCost onlygithub.Cents,
	header *multipart.FileHeader,
) (onlygithub.ID, error) {
	data := onlygithub.UploadImageRequest{
		Filename:    header.Filename,
		Visibility:  visibility,
		MinimumCost: minimumCost,
	}

	contentType := header.Header.Get("Content-Type")
	if contentType == "" {
		contentType = mime.TypeByExtension(header.Filename)
	}
	if queryableImageTypes[contentType] {
		query, err := queryImage(header)
		if err != nil {
			log := hclog.FromContext(ctx)
			log.Info("couldn't query image", "file", header.Filename, "err", err)
		} else {
			data.PreviewURL = query.PreviewURL
			data.Width = query.Width
			data.Height = query.Height
		}
	}

	f, err := header.Open()
	if err != nil {
		return onlygithub.NullID, hrt.WrapHTTPError(500, errors.Wrap(err, "failed to open image"))
	}
	defer f.Close() // runs after the function returns

	asset, err := h.Images.UploadImage(ctx, data, f)
	if err != nil {
		return onlygithub.NullID, hrt.WrapHTTPError(500, errors.Wrap(err, "failed to upload image"))
	}

	return asset.ID, nil
}

var queryableImageTypes = map[string]bool{
	"image/png":  true,
	"image/jpeg": true,
}

type imageQuery struct {
	PreviewURL string
	Width      int
	Height     int
}

func queryImage(header *multipart.FileHeader) (imageQuery, error) {
	f, err := header.Open()
	if err != nil {
		return imageQuery{}, errors.Wrap(err, "failed to open image")
	}
	defer f.Close()

	img, _, err := image.Decode(f)
	if err != nil {
		return imageQuery{}, errors.Wrap(err, "failed to decode image")
	}

	previewWidth, previewHeight := clampSize(img.Bounds().Dx(), img.Bounds().Dy(), 8)
	preview := image.NewRGBA(image.Rect(0, 0, previewWidth, previewHeight))
	draw.CatmullRom.Scale(preview, preview.Bounds(), img, img.Bounds(), draw.Over, nil)

	var previewJPEG bytes.Buffer
	if err := jpeg.Encode(&previewJPEG, preview, &jpeg.Options{Quality: 80}); err != nil {
		return imageQuery{}, errors.Wrap(err, "failed to encode preview")
	}

	var previewData strings.Builder
	previewData.WriteString("data:image/jpeg;base64,")
	previewData.WriteString(base64.StdEncoding.EncodeToString(previewJPEG.Bytes()))

	return imageQuery{
		PreviewURL: previewData.String(),
		Width:      img.Bounds().Dx(),
		Height:     img.Bounds().Dy(),
	}, nil
}

func clampSize(w, h, max int) (int, int) {
	if w > max || h > max {
		if w > h {
			h = h * max / w
			w = max
		} else {
			w = w * max / h
			h = max
		}
	}
	return w, h
}
