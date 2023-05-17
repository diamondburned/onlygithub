package frontend

import (
	"context"
	"embed"
	"net/http"

	"github.com/diamondburned/hrt"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"
	"libdb.so/onlygithub"
	"libdb.so/onlygithub/frontend/layouts"
	"libdb.so/onlygithub/internal/auth"
	"libdb.so/onlygithub/internal/gh"
)

//go:embed static
var staticFS embed.FS

// StaticHandler returns a handler that serves the /static folder. The files in
// /static will be served from root.
func StaticHandler() http.Handler {
	return http.FileServer(http.FS(staticFS))
}

// ServeJS creates a new handler that serves the given JS blob.
func ServeJS(blob string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/javascript")
		w.Write([]byte(blob))
	}
}

type ctxKey int

const (
	_ ctxKey = iota
	siteKey
	ownerKey
	sessionKey
)

// Session is the session data.
type Session struct {
	*oauth2.Token
	GitHub *gh.Client
	Me     *onlygithub.User
}

// Deps is the dependencies for the frontend.
type Deps struct {
	Tiers       onlygithub.TierService
	Users       onlygithub.UserService
	Posts       onlygithub.PostService
	Images      onlygithub.ImageService
	Config      onlygithub.ConfigService
	GitHubOAuth *auth.GitHubAuthorizer
}

func (d *Deps) RenderingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		site, err := d.Config.SiteConfig(r.Context())
		if err != nil {
			layouts.RenderError(w, r, errors.Wrap(err, "failed to get site config"))
			return
		}

		owner, err := d.Users.Owner(r.Context())
		if err != nil {
			if errors.Is(err, onlygithub.ErrNotFound) {
				layouts.RenderError(w, r, hrt.NewHTTPError(
					http.StatusInternalServerError,
					"missing owner, please head to /login to add your account and use `onlyserve --mkowner`"))
			} else {
				layouts.RenderError(w, r, errors.Wrap(err, "failed to get owner"))
			}
			return
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, siteKey, site)
		ctx = context.WithValue(ctx, ownerKey, owner)

		oauth := auth.TokenFromRequest(r)
		if oauth != nil {
			source := d.GitHubOAuth.Config.TokenSource(r.Context(), oauth)
			client := gh.NewClient(r.Context(), source)

			me, err := client.Me()
			if err != nil {
				layouts.RenderError(w, r, errors.Wrap(err, "failed to get user"))
				return
			}

			user, err := d.Users.User(r.Context(), me.ID)
			if err != nil {
				if errors.Is(err, onlygithub.ErrNotFound) {
					// Create the user.
					if err = d.Users.UpdateUser(r.Context(), me); err != nil {
						layouts.RenderError(w, r, errors.Wrap(err, "failed to create user"))
						return
					}
					user = me
				} else {
					layouts.RenderError(w, r, errors.Wrap(err, "failed to get user"))
					return
				}
			}

			ctx = context.WithValue(ctx, sessionKey, &Session{
				Token:  oauth,
				GitHub: client,
				Me:     user,
			})
		}

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func SiteFromRequest(r *http.Request) *onlygithub.SiteConfig {
	return r.Context().Value(siteKey).(*onlygithub.SiteConfig)
}

func OwnerFromRequest(r *http.Request) *onlygithub.User {
	return r.Context().Value(ownerKey).(*onlygithub.User)
}

func SessionFromRequest(r *http.Request) *Session {
	v, _ := r.Context().Value(sessionKey).(*Session)
	return v
}

// OwnerOnly is a middleware that checks if the user is the owner.
func OwnerOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session := SessionFromRequest(r)
		if session == nil || !session.Me.IsOwner {
			layouts.RenderError(w, r, onlygithub.ErrUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func LoggedInOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session := SessionFromRequest(r)
		if session == nil {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}

		next.ServeHTTP(w, r)
	})
}
