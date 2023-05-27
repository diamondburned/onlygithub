package about

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"libdb.so/onlygithub/frontend"
)

func New() http.Handler {
	r := chi.NewRouter()
	r.Use(frontend.EnforceHomepageVisibility)
	r.Get("/", get)
	return r
}

func get(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	site := frontend.SiteFromRequest(r)
	owner := frontend.OwnerFromRequest(r)

	about(r, site, owner).Render(ctx, w)
}
