package index

import (
	"net/http"

	"libdb.so/onlygithub"
	"libdb.so/onlygithub/frontend"
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
			kickToLogin(w, r)
			return
		}
	case onlygithub.VisibleToSponsors:
		if opts.Me == nil || opts.Me.Sponsorship == nil {
			kickToLogin(w, r)
			return
		}
	case onlygithub.VisibleToPrivate:
		if opts.Me == nil {
			kickToLogin(w, r)
			return
		}
	case onlygithub.VisibleToPublic:
		// Nothing.
	}

	index(r, site, owner, opts).Render(r.Context(), w)
}

func kickToLogin(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/login", http.StatusFound)
}
