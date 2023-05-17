package db

import (
	"io"

	"libdb.so/onlygithub"
	"libdb.so/onlygithub/internal/auth"
)

// Database is the interfaces implemented by the database.
type Database interface {
	io.Closer
	auth.GitHubTokenService
	onlygithub.UserService
	onlygithub.TierService
	onlygithub.ConfigService
	onlygithub.OAuthTokenService
	onlygithub.PrivilegedUserService
}
