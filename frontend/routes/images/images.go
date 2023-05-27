package images

import (
	"fmt"
	"io"
	"mime"
	"net/http"
	"path/filepath"

	"github.com/diamondburned/hrt"
	"github.com/go-chi/chi/v5"
	"github.com/rs/xid"
	"libdb.so/onlygithub"
	"libdb.so/onlygithub/frontend"
	"libdb.so/onlygithub/frontend/layouts"
	"libdb.so/onlygithub/internal/auth"
)

// ImageService extends the onlygithub.ImageService with image data.
type ImageService interface {
	onlygithub.ImageService
}

func New(images ImageService, oauth *auth.OAuthMiddleware) http.Handler {
	r := chi.NewRouter()
	r.Get("/{id}", (&handler{images}).get)

	return r
}

type handler struct {
	isrv ImageService
}

func (h *handler) get(w http.ResponseWriter, r *http.Request) {
	id, err := xid.FromString(chi.URLParam(r, "id"))
	if err != nil {
		layouts.RenderError(w, r, hrt.WrapHTTPError(http.StatusBadRequest, err))
		return
	}

	image, err := h.isrv.Image(r.Context(), id)
	if err != nil {
		layouts.RenderError(w, r, err)
		return
	}

	session := frontend.SessionFromRequest(r)
	if !image.IsVisibleTo(session.Me) {
		layouts.RenderError(w, r, hrt.WrapHTTPError(http.StatusNotFound, err))
		return
	}

	mimeType := mime.TypeByExtension(filepath.Ext(image.Filename))
	if mimeType == "" {
		mimeType = "application/octet-stream"
	}

	w.Header().Set("Content-Type", mimeType)
	w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename=%q`, image.Filename))
	w.Header().Set("Cache-Control", "private, max-age=31536000, immutable")

	data, err := h.isrv.ImageData(r.Context(), id)
	if err != nil {
		layouts.RenderError(w, r, err)
		return
	}
	defer data.Close()

	if _, err := io.Copy(w, data); err != nil {
		layouts.RenderError(w, r, err)
		return
	}
}
