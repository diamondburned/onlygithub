package auth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/diamondburned/hrt"
	"github.com/go-chi/chi/v5"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"
	"libdb.so/onlygithub"
	"libdb.so/onlygithub/internal/api"
)

type ctxKey uint8

const (
	_ ctxKey = iota
	oauthTokenCtxKey
)

func generateToken(prefix string) string {
	var buf [24]byte
	_, err := rand.Read(buf[:])
	if err != nil {
		panic(err)
	}
	return prefix + "-" + base64.URLEncoding.EncodeToString(buf[:])
}

// OAuthConfig is the configuration for any OAuth provider.
type OAuthConfig struct {
	ID           string // Client ID
	Secret       string // Client Secret
	RootEndpoint string // Root endpoint for mounting router
}

// OAuthAuthorizer is an interface for authorizing OAuth requests.
// It supplies middlewares and handlers with the ability to authorize
// requests.
type OAuthAuthorizer struct {
	chi.Router
	Config   *oauth2.Config
	Provider string

	tokens onlygithub.OAuthTokenService
}

// NewOAuthAuthorizer returns a new OAuthAuthorizer.
func NewOAuthAuthorizer(provider string, config *oauth2.Config, tokenService onlygithub.OAuthTokenService) *OAuthAuthorizer {
	a := &OAuthAuthorizer{
		Config:   config,
		tokens:   tokenService,
		Provider: provider,
	}

	a.Router = chi.NewRouter()
	a.Router.Get("/callback", a.handleCallback)
	a.Router.Get("/logout", a.handleLogout)
	a.Router.Get("/", a.handleLogin)

	return a
}

// handleLogin is the route that redirects the user to the OAuth provider
// to authorize the application.
func (a *OAuthAuthorizer) handleLogin(w http.ResponseWriter, r *http.Request) {
	ticket := generateToken(a.Provider + "-state")
	http.SetCookie(w, &http.Cookie{
		Name:    a.Provider + "-state",
		Value:   ticket,
		Expires: time.Now().Add(30 * time.Minute),
	})

	// Preserve the origin so we can redirect back to it after
	// the OAuth flow is complete.
	origin := r.Header.Get("Origin")
	if origin != "" {
		http.SetCookie(w, &http.Cookie{
			Name:  a.Provider + "-origin",
			Value: origin,
		})
	}

	config := *a.Config
	config.RedirectURL = "http://" + r.Host + r.URL.Path + "/callback"
	log.Println("set callback url to", config.RedirectURL)

	redirect := config.AuthCodeURL(ticket)
	http.Redirect(w, r, redirect, http.StatusFound)
}

func (a *OAuthAuthorizer) handleLogout(w http.ResponseWriter, r *http.Request) {
	if err := a.Logout(w, r); err != nil {
		api.RespondError(w, r, hrt.WrapHTTPError(http.StatusInternalServerError, err))
		return
	}
}

// handleCallback is the route that the OAuth provider redirects the user
// to after authorizing the application.
func (a *OAuthAuthorizer) handleCallback(w http.ResponseWriter, r *http.Request) {
	stateCookie, err := r.Cookie(a.Provider + "-state")
	if err != nil {
		api.RespondError(w, r, hrt.WrapHTTPError(http.StatusBadRequest, err))
		return
	}

	if stateCookie.Value != r.FormValue("state") {
		api.RespondError(w, r, hrt.NewHTTPError(http.StatusBadRequest, "invalid state"))
		return
	}

	token, err := a.Config.Exchange(r.Context(), r.FormValue("code"))
	if err != nil {
		api.RespondError(w, r, hrt.WrapHTTPError(http.StatusBadRequest, errors.Wrap(err, "code exchange failed")))
		return
	}

	ourToken := generateToken(a.Provider + "-token")
	http.SetCookie(w, &http.Cookie{
		Name:     a.Provider + "-token",
		Value:    ourToken,
		Path:     "/",
		SameSite: http.SameSiteStrictMode,
	})

	if err = a.tokens.SaveToken(r.Context(), ourToken, a.Provider, token); err != nil {
		api.RespondError(w, r, hrt.WrapHTTPError(http.StatusInternalServerError, errors.Wrap(err, "failed to save token")))
		return
	}

	originCookie, err := r.Cookie(a.Provider + "-origin")
	if err == nil {
		// Clear this cookie; we don't need it anymore.
		http.SetCookie(w, &http.Cookie{
			Name:   a.Provider + "-origin",
			Value:  "",
			MaxAge: -1,
		})
		http.Redirect(w, r, originCookie.Value, http.StatusFound)
	} else {
		http.Redirect(w, r, "/", http.StatusFound)
	}
}

// handleLogout is the route that logs the user out.
func (a *OAuthAuthorizer) Logout(w http.ResponseWriter, r *http.Request) error {
	cookie, err := r.Cookie(a.Provider + "-token")
	if err != nil {
		return hrt.WrapHTTPError(http.StatusBadRequest, err)
	}

	if err = a.tokens.DeleteToken(r.Context(), cookie.Value, a.Provider); err != nil {
		return hrt.WrapHTTPError(
			http.StatusInternalServerError,
			errors.Wrap(err, "failed to delete token"))
	}

	http.SetCookie(w, &http.Cookie{
		Name:     a.Provider + "-token",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		SameSite: http.SameSiteStrictMode,
	})
	return nil
}

// FromRequest returns the OAuth token from the request.
func (a *OAuthAuthorizer) FromRequest(r *http.Request) (*oauth2.Token, error) {
	cookie, err := r.Cookie(a.Provider + "-token")
	if err != nil {
		return nil, fmt.Errorf("no such cookie %s-token", a.Provider)
	}
	return a.tokens.RetrieveToken(r.Context(), cookie.Value, a.Provider)
}

// OAuthMiddleware contains middlewares for authorizing requests.
type OAuthMiddleware struct {
	*OAuthAuthorizer
	Prefix string
}

// Middleware creates an OAuthMiddleware with the given prefix.
func (a *OAuthAuthorizer) Middleware(prefix string) *OAuthMiddleware {
	return &OAuthMiddleware{
		OAuthAuthorizer: a,
		Prefix:          prefix,
	}
}

// Use returns a middleware that authorizes requests. It is required for
// TokenFromRequest to work.
func (a *OAuthMiddleware) Use() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token, err := a.FromRequest(r)
			if err == nil {
				ctx := context.WithValue(r.Context(), oauthTokenCtxKey, token)
				r = r.WithContext(ctx)
			}
			next.ServeHTTP(w, r)
		})
	}
}

// Require is similar to Use, except it 401s the user if they are not
// authorized.
func (a *OAuthMiddleware) Require() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := TokenFromRequest(r)
			if token == nil {
				api.RespondError(w, r, hrt.NewHTTPError(http.StatusUnauthorized, "not authorized"))
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

// RequireOrRedirect is similar to Use, except it redirects the user to the
// given URL if they are not authorized. Use this to establish a login wall.
func (a *OAuthMiddleware) RequireOrRedirect() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := TokenFromRequest(r)
			if token == nil {
				http.Redirect(w, r, a.Prefix, http.StatusFound)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

// TokenFromRequest returns the OAuth token from the request.
func TokenFromRequest(r *http.Request) *oauth2.Token {
	token, _ := r.Context().Value(oauthTokenCtxKey).(*oauth2.Token)
	return token
}
