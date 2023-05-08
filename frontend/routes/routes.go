package routes

import (
	"net/http"

	"github.com/diamondburned/hrt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"libdb.so/onlygithub/frontend"
	"libdb.so/onlygithub/frontend/layouts"
	"libdb.so/onlygithub/frontend/routes/index"
)

// New returns a new page router.
func New(d frontend.Deps) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(d.GitHubOAuth.Use())

	r.Mount("/static", frontend.StaticHandler())
	r.Route("/login", func(r chi.Router) {
		r.Mount("/github", d.GitHubOAuth)
		r.Handle("/", redirectHandler("/login/github"))
	})
	r.HandleFunc("/logout", unimplemented)

	r.Group(func(r chi.Router) {
		r.Use(d.RenderingMiddleware)
		r.NotFound(func(w http.ResponseWriter, r *http.Request) {
			layouts.RenderError(w, r, hrt.NewHTTPError(http.StatusNotFound, "page not found"))
		})

		r.Get("/", index.Handle)
	})

	return r
}

func redirectHandler(path string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, path, http.StatusFound)
	})
}

func unimplemented(w http.ResponseWriter, r *http.Request) {
	layouts.RenderError(w, r, hrt.NewHTTPError(http.StatusNotImplemented, "not implemented"))
}
