package auth

import (
	"context"

	"github.com/pkg/errors"
	"golang.org/x/oauth2"
	"libdb.so/onlygithub"
	"libdb.so/onlygithub/internal/gh"

	githuboauth "golang.org/x/oauth2/github"
)

// GitHubAuthorizer is an OAuthAuthorizer for GitHub.
type GitHubAuthorizer struct {
	OAuthAuthorizer
}

// GitHubProvider is the string used to identify GitHub tokens in
// OAuthTokenService.
const GitHubProvider = "github"

// GitHubConfig is the configuration for GitHub OAuth.
type GitHubConfig struct {
	ID          string
	Secret      string
	RedirectURL string
}

// GitHubTokenService is an OAuthTokenService for GitHub.
type GitHubTokenService interface {
	OAuthTokenService
	onlygithub.UserService
}

// NewGitHubAuthorizer returns a new GitHubAuthorizer.
func NewGitHubAuthorizer(cfg GitHubConfig, tokens GitHubTokenService) *GitHubAuthorizer {
	oacfg := &oauth2.Config{
		ClientID:     cfg.ID,
		ClientSecret: cfg.Secret,
		RedirectURL:  cfg.RedirectURL,
		Endpoint:     githuboauth.Endpoint,
	}
	oa := NewOAuthAuthorizer(GitHubProvider, oacfg, githubOAuthTokenService{
		oacfg:  oacfg,
		tokens: tokens,
	})
	return &GitHubAuthorizer{*oa}
}

type githubOAuthTokenService struct {
	oacfg  *oauth2.Config
	tokens GitHubTokenService
}

// SaveToken saves the OAuth token for the given user.
func (s githubOAuthTokenService) SaveToken(ctx context.Context, token, provider string, oauthToken *oauth2.Token) error {
	if err := s.tokens.SaveToken(ctx, token, provider, oauthToken); err != nil {
		return err
	}

	client := gh.NewClient(ctx, s.oacfg.TokenSource(ctx, oauthToken))
	user, err := client.Me()
	if err != nil {
		return errors.Wrap(err, "failed to get user")
	}

	if err := s.tokens.UpdateUser(ctx, user); err != nil {
		return errors.Wrap(err, "failed to update user into database")
	}

	return nil
}

// RetrieveToken retrieves the OAuth token for the given user.
func (s githubOAuthTokenService) RetrieveToken(ctx context.Context, token, provider string) (*oauth2.Token, error) {
	return s.tokens.RetrieveToken(ctx, token, provider)
}
