package index

import (
	"net/http"

	"libdb.so/onlygithub/frontend"
)

func Handle(w http.ResponseWriter, r *http.Request) {
	site := frontend.SiteFromRequest(r)
	owner := frontend.OwnerFromRequest(r)

	var opts indexOpts

	session := frontend.SessionFromRequest(r)
	if session != nil {
		opts.Me = session.Me
	}

	index(r, site, owner, opts).Render(r.Context(), w)
}
