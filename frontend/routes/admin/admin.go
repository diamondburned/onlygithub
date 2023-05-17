package admin

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/hashicorp/go-hclog"
	"libdb.so/onlygithub"
	"libdb.so/onlygithub/frontend"
	"libdb.so/onlygithub/frontend/layouts"
)

func New(tiers onlygithub.TierService) http.Handler {
	h := &handler{tiers}

	r := chi.NewRouter()
	r.Use(frontend.OwnerOnly)

	r.Route("/tiers", func(r chi.Router) {
		r.Get("/refresh", h.refreshTiers)
	})

	return r
}

func adminOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session := frontend.SessionFromRequest(r)
		if session == nil || !session.Me.IsOwner {
			layouts.RenderError(w, r, onlygithub.ErrUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}

type handler struct {
	tiers onlygithub.TierService
}

func (h *handler) refreshTiers(w http.ResponseWriter, r *http.Request) {
	session := frontend.SessionFromRequest(r)

	tiers, err := session.GitHub.Tiers(r.Context(), 100).All()
	if err != nil {
		layouts.RenderError(w, r, err)
		return
	}

	if err := h.tiers.UpdateTiers(r.Context(), tiers); err != nil {
		layouts.RenderError(w, r, err)
		return
	}

	log := hclog.FromContext(r.Context())
	log.Info("refreshed tiers", "count", len(tiers), "tiers", tiers)
}
