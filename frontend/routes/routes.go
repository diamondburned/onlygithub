package routes

import (
	"net/http"

	"github.com/diamondburned/hrt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"libdb.so/onlygithub/dist"
	"libdb.so/onlygithub/frontend"
	"libdb.so/onlygithub/frontend/layouts"
	"libdb.so/onlygithub/frontend/routes/about"
	"libdb.so/onlygithub/frontend/routes/admin"
	"libdb.so/onlygithub/frontend/routes/create"
	"libdb.so/onlygithub/frontend/routes/images"
	"libdb.so/onlygithub/frontend/routes/index"
	"libdb.so/onlygithub/frontend/routes/settings"
)

// New returns a new page router.
func New(d frontend.Deps) http.Handler {
	oauthMiddleware := d.GitHubOAuth.Middleware("/login")

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.CleanPath)
	r.Use(middleware.Recoverer)
	r.Use(oauthMiddleware.Use())

	r.Mount("/static", frontend.TryFiles(dist.StaticFS))
	r.Route("/login", func(r chi.Router) {
		r.Mount("/github", d.GitHubOAuth)
		r.Handle("/", redirectHandler("/login/github"))
	})
	r.Post("/logout", func(w http.ResponseWriter, r *http.Request) {
		if err := d.GitHubOAuth.Logout(w, r); err != nil {
			layouts.RenderError(w, r, err)
			return
		}
		http.Redirect(w, r, "/", http.StatusFound)
	})

	r.Group(func(r chi.Router) {
		r.Use(d.RenderingMiddleware)
		r.NotFound(func(w http.ResponseWriter, r *http.Request) {
			layouts.RenderError(w, r, hrt.NewHTTPError(http.StatusNotFound, "page not found"))
		})

		r.Mount("/create", create.New(create.Services{
			Images: d.Images,
			Posts:  d.Posts,
			Tiers:  d.Tiers,
		}))
		r.Mount("/settings", settings.New(settings.Services{
			Config: d.Config,
			Images: d.Images,
		}))
		r.Mount("/images", images.New(d.Images, oauthMiddleware))
		r.Mount("/admin", admin.New(d.Tiers))
		r.Mount("/about", about.New())
		r.Mount("/membership", unimplemented)
		r.Mount("/", index.New(index.Services{
			Posts: d.Posts,
		}))
	})

	return r
}

func redirectHandler(path string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, path, http.StatusFound)
	})
}

var unimplemented = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	layouts.RenderError(w, r, hrt.NewHTTPError(http.StatusNotImplemented, "not implemented"))
})
