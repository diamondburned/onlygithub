package settings

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"libdb.so/onlygithub"
	"libdb.so/onlygithub/frontend"
	"libdb.so/onlygithub/frontend/layouts"
)

type handler struct {
	config onlygithub.ConfigService
}

func New(config onlygithub.ConfigService) http.Handler {
	h := &handler{config}

	r := chi.NewRouter()
	r.Use(frontend.LoggedInOnly)
	r.Get("/", h.get)

	return r
}

func (h handler) get(w http.ResponseWriter, r *http.Request) {
	site := frontend.SiteFromRequest(r)
	owner := frontend.OwnerFromRequest(r)
	session := frontend.SessionFromRequest(r)

	cfg, err := h.config.UserConfig(r.Context(), session.Me.ID)
	if err != nil {
		layouts.RenderError(w, r, err)
		return
	}

	data := settingsData{
		Me:       session.Me,
		MeConfig: cfg,
	}

	if session.Me.IsOwner {
		siteCfg, err := h.config.SiteConfig(r.Context())
		if err != nil {
			layouts.RenderError(w, r, err)
			return
		}

		data.SiteConfig = siteCfg
	}

	settings(r, site, owner, data).Render(r.Context(), w)
}
