package db

import (
	"io"

	"libdb.so/onlygithub"
	"libdb.so/onlygithub/internal/auth"
)

//go:generate sqlc generate

// Database is the interfaces implemented by the database.
type Database interface {
	io.Closer
	auth.GitHubTokenService
	auth.OAuthTokenService
	onlygithub.UserService
	onlygithub.ConfigService
	onlygithub.PrivilegedUserService
}
