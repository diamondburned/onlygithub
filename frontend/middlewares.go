package frontend

import (
	"net/http"

	"github.com/pkg/errors"
	"libdb.so/onlygithub"
	"libdb.so/onlygithub/frontend/layouts"
)

// OwnerOnly is a middleware that checks if the user is the owner.
func OwnerOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session := SessionFromRequest(r)
		if session.Me == nil || !session.Me.IsOwner {
			layouts.RenderError(w, r, onlygithub.ErrUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// LoggedInOnly is a middleware that checks if the user is logged in.
func LoggedInOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session := SessionFromRequest(r)
		if session.Me == nil {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// ParseForm is a middleware that parses the form.
func ParseForm(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			layouts.RenderError(w, r, errors.Wrap(err, "failed to parse form"))
			return
		}

		next.ServeHTTP(w, r)
	})
}

const defaultMaxMemory = 1 << 20 // 1 MB

// ParseMultipartForm is a middleware that parses the multipart form.
func ParseMultipartForm(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseMultipartForm(defaultMaxMemory); err != nil {
			layouts.RenderError(w, r, errors.Wrap(err, "failed to parse multipart form"))
			return
		}

		next.ServeHTTP(w, r)
	})
}

// EnforceHomepageVisibility is a middleware that enforces the homepage
// visibility settings. If the user is not allowed to view the homepage, then
// they will be redirected to the login page.
func EnforceHomepageVisibility(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		site := SiteFromRequest(r)
		owner := OwnerFromRequest(r)
		session := SessionFromRequest(r)

		switch site.HomepageVisibility {
		case onlygithub.NotVisible:
			if session.Me == nil || session.Me.ID != owner.ID {
				kickToLogin(w, r)
				return
			}
		case onlygithub.VisibleToSponsors:
			if session.Me == nil || session.Me.Sponsorship == nil {
				kickToLogin(w, r)
				return
			}
		case onlygithub.VisibleToPrivate:
			if session.Me == nil {
				kickToLogin(w, r)
				return
			}
		case onlygithub.VisibleToPublic:
			// Nothing.
		}

		next.ServeHTTP(w, r)
	})
}
