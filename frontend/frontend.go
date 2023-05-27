package frontend

import (
	"context"
	"io"
	"io/fs"
	"net/http"
	"strings"

	"github.com/diamondburned/hrt"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"
	"libdb.so/onlygithub"
	"libdb.so/onlygithub/frontend/layouts"
	"libdb.so/onlygithub/internal/auth"
	"libdb.so/onlygithub/internal/gh"
)

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

			ctx = context.WithValue(ctx, sessionKey, Session{
				Token:  oauth,
				GitHub: client,
				Me:     user,
			})
		}

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// SiteFromRequest returns the site config from the request.
func SiteFromRequest(r *http.Request) *onlygithub.SiteConfig {
	return r.Context().Value(siteKey).(*onlygithub.SiteConfig)
}

// OwnerFromRequest returns the owner from the request.
func OwnerFromRequest(r *http.Request) *onlygithub.User {
	return r.Context().Value(ownerKey).(*onlygithub.User)
}

// SessionFromRequest returns the session from the request.
func SessionFromRequest(r *http.Request) Session {
	v, _ := r.Context().Value(sessionKey).(Session)
	return v
}

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

// TryFiles returns a handler that tries to serve the files from the filesystems
// in order. If none of the filesystems have the file, it will serve the
// NotFoundHandler.
func TryFiles(fses ...fs.FS) http.Handler {
	errNotFound := errors.New("file not found")

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		name := r.URL.Path
		if name == "/" {
			name = "."
		} else {
			name = strings.TrimPrefix(name, "/")
		}

		for _, fs := range fses {
			f, err := fs.Open(name)
			if err != nil {
				continue
			}
			defer f.Close()

			fi, err := f.Stat()
			if err != nil {
				layouts.RenderError(w, r, errors.Wrap(err, "failed to stat file"))
				return
			}

			if fi.IsDir() {
				continue
			}

			if seeker, ok := f.(io.ReadSeeker); ok {
				http.ServeContent(w, r, fi.Name(), fi.ModTime(), seeker)
			} else {
				io.Copy(w, f)
			}

			return
		}

		layouts.RenderError(w, r, errNotFound)
	})
}
