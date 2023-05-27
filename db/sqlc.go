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
	onlygithub.PostService
	onlygithub.TierService
	onlygithub.ImageService
	onlygithub.ConfigService
	onlygithub.OAuthTokenService
	onlygithub.PrivilegedUserService
}

/*
type ctxKey int

const (
	_ ctxKey = iota
	currentUserKey
)

// WithCurrentUser returns a new context with the current user set.
// This is used by database implementations to authorize actions.
func WithCurrentUser(ctx context.Context, user *onlygithub.User) context.Context {
	ctx = context.WithValue(ctx, currentUserKey, user)
	return ctx
}

func currentUser(ctx context.Context) *onlygithub.User {
	user, _ := ctx.Value(currentUserKey).(*onlygithub.User)
	return user
}
*/
