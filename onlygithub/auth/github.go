package auth

import (
	"golang.org/x/oauth2"

	githuboauth "golang.org/x/oauth2/github"
)

// GitHubAuthorizer is an OAuthAuthorizer for GitHub.
type GitHubAuthorizer struct {
	OAuthAuthorizer
}

type GitHubConfig struct {
	ID     string
	Secret string
}

// NewGitHubAuthorizer returns a new GitHubAuthorizer.
func NewGitHubAuthorizer(cfg GitHubConfig, tokens OAuthTokenService) *GitHubAuthorizer {
	oacfg := &oauth2.Config{
		ClientID:     cfg.ID,
		ClientSecret: cfg.Secret,
		Endpoint:     githuboauth.Endpoint,
	}
	oa := NewOAuthAuthorizer("github", oacfg, tokens)
	return &GitHubAuthorizer{*oa}
}
