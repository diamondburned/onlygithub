package create

import (
	"context"
	"fmt"
	"mime/multipart"
	"net/http"
	"strconv"

	_ "embed"

	"github.com/diamondburned/hrt"
	"github.com/go-chi/chi/v5"
	"github.com/hashicorp/go-hclog"
	"github.com/pkg/errors"
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

	minCost, err := strconv.ParseFloat(r.FormValue("minimum-cost"), 64)
	if err != nil {
		layouts.RenderError(w, r, hrt.WrapHTTPError(
			http.StatusBadRequest, errors.Wrap(err, "mininum cost is not a valid number")))
		return
	}

	minCost *= 100 // convert to cents

	visibility := onlygithub.Visibility(r.FormValue("visibility"))
	if err := visibility.Validate(); err != nil {
		layouts.RenderError(w, r, hrt.WrapHTTPError(
			http.StatusBadRequest, errors.Wrap(err, "invalid visibility")))
		return
	}

	assetIDs, err := uploadImages(
		r.Context(), h.Images,
		visibility, onlygithub.Cents(minCost), r.MultipartForm.File["images"])
	if err != nil {
		layouts.RenderError(w, r, err)
		return
	}

	create := onlygithub.CreatePostRequest{
		Visibility:  visibility,
		MinimumCost: onlygithub.Cents(minCost),
		Markdown:    r.FormValue("content"),
		AssetIDs:    assetIDs,
	}
	if r.FormValue("disable-comments") != "" {
		create.AllowComments = new(bool)
	}
	if r.FormValue("disable-reactions") != "" {
		create.AllowReactions = new(bool)
	}

	post, err := h.Posts.CreatePost(r.Context(), create)
	if err != nil {
		layouts.RenderError(w, r, hrt.WrapHTTPError(500, err))
		return
	}

	postURL := fmt.Sprintf("/posts/%d", post.ID)
	http.Redirect(w, r, postURL, http.StatusFound)
}

func uploadImages(
	ctx context.Context,
	isrv onlygithub.ImageService,
	visibility onlygithub.Visibility, minimumCost onlygithub.Cents,
	headers []*multipart.FileHeader) ([]onlygithub.ID, error) {

	var assetIDs []onlygithub.ID

	fail := true
	defer func() {
		if !fail {
			return
		}

		log := hclog.FromContext(ctx)
		for _, id := range assetIDs {
			if err := isrv.DeleteImage(ctx, id); err != nil {
				log.Warn("failed to delete image after failure", "id", id, "err", err)
			}
		}
	}()

	for _, header := range headers {
		if header.Filename == "" {
			return nil, hrt.NewHTTPError(400, "image is missing filename")
		}

		f, err := header.Open()
		if err != nil {
			return nil, hrt.WrapHTTPError(500, errors.Wrap(err, "failed to open image"))
		}
		defer f.Close() // runs after the function returns

		data := onlygithub.UploadImageRequest{
			Filename:    header.Filename,
			Visibility:  visibility,
			MinimumCost: minimumCost,
		}

		asset, err := isrv.UploadImage(ctx, data, f)
		if err != nil {
			return nil, hrt.WrapHTTPError(500, errors.Wrap(err, "failed to upload image"))
		}

		assetIDs = append(assetIDs, asset.ID)
	}

	fail = false
	return assetIDs, nil
}
