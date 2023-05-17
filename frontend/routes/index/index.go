package index

import (
	"net/http"

	"libdb.so/onlygithub"
	"libdb.so/onlygithub/frontend"
	"libdb.so/onlygithub/frontend/layouts"
)

func GET(w http.ResponseWriter, r *http.Request) {
	site := frontend.SiteFromRequest(r)
	owner := frontend.OwnerFromRequest(r)

	var opts indexOpts

	session := frontend.SessionFromRequest(r)
	if session != nil {
		opts.Me = session.Me
	}

	switch site.HomepageVisibility {
	case onlygithub.NotVisible:
		if opts.Me == nil || opts.Me.ID != owner.ID {
			layouts.RenderError(w, r, onlygithub.ErrUnauthorized)
			return
		}
	case onlygithub.VisibleToSponsors:
		if opts.Me == nil || opts.Me.Sponsorship == nil {
			layouts.RenderError(w, r, onlygithub.ErrUnauthorized)
			return
		}
	case onlygithub.VisibleToPrivate:
		if opts.Me == nil {
			layouts.RenderError(w, r, onlygithub.ErrUnauthorized)
			return
		}
	case onlygithub.VisibleToPublic:
		// Nothing.
	}

	index(r, site, owner, opts).Render(r.Context(), w)
}
