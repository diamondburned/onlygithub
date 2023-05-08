package onlygithub

import (
	"context"
	"io"
)

// PrivilegedUserService should never be used by the frontend. It is only
// intended to be used by the backend.
type PrivilegedUserService interface {
	// MakeOwner makes the user with the given username the owner of the site.
	MakeOwner(ctx context.Context, username string) error
}

// UserService is a service that manages users.
type UserService interface {
	// User returns the user with the given ID. If no user exists with the given
	// ID, an error should be returned.
	User(ctx context.Context, id GitHubID) (*User, error)
	// UpdateUser updates the user with the given GitHub user.
	UpdateUser(ctx context.Context, user *User) error
	// Owner returns the owner of the site. This is the user that owns the site
	// and has full control over it.
	Owner(ctx context.Context) (*User, error)
}

// PostService is a service that manages posts.
type PostService interface {
	// Posts returns a list of posts. The page size can be anything, but it
	// should generally be over or around 10 at least.
	//
	// Posts should be sorted latest-first, meaning higher IDs should be
	// returned first.
	//
	// Visibility should act as a filter: if Public is specified, only public
	// posts should be returned, but if Sponsor is specified, both public and
	// sponsor-only posts should be returned.
	Posts(ctx context.Context, visibility Visibility, before ID) ([]Post, error)
	// UploadImage uploads an image. The image should be uploaded with the given
	// filename. The filename should be validated.
	UploadImage(ctx context.Context, r io.Reader, req UploadImageRequest) (*ImageAsset, error)
	// CreatePost creates a post. The post should be created with the given
	// visibility. The visibility should be validated.
	CreatePost(ctx context.Context, req CreatePostRequest) (*Post, error)
}

// UploadImageRequest is a request to upload an image.
type UploadImageRequest struct {
	// Filename is the filename of the image.
	Filename string `json:"filename"`
	// Visibility is the visibility of the image.
	Visibility Visibility `json:"visibility"`
	// MinimumCost is the minimum cost of the image.
	MinimumCost Cents `json:"minimumCost"`
}

// CreatePostRequest is a request to create a post.
type CreatePostRequest struct {
	// Visibility is the visibility of the image.
	Visibility Visibility `json:"visibility"`
	// MinimumCost is the minimum cost of the image.
	MinimumCost Cents `json:"minimumCost"`
	// Markdown is the markdown of the post.
	Markdown string `json:"markdown"`
	// AssetIDs is the list of asset IDs to attach to the post.
	AssetIDs []ID `json:"assetIDs"`
	// AllowComments is whether comments are allowed on the post.
	// If null, the default value should be used.
	AllowComments *bool `json:"allowComments,omitempty"`
	// AllowReactions is whether reactions are allowed on the post.
	// If null, the default value should be used.
	AllowReactions *bool `json:"allowReactions,omitempty"`
}

// ConfigService is a service that manages the configuration of the site.
type ConfigService interface {
	// SiteConfig returns the site-wide configuration. If no site-wide
	// configuration exists, a default configuration will be returned.
	SiteConfig(ctx context.Context) (*SiteConfig, error)
	// SetSiteConfig updates the site-wide configuration.
	SetSiteConfig(ctx context.Context, cfg *SiteConfig) error
	// UserConfig returns the configuration for the given user. If no
	// configuration exists for the user, a default configuration will be
	// returned.
	UserConfig(ctx context.Context, userID GitHubID) (*UserConfig, error)
	// SetUserConfig updates the configuration for the given user.
	SetUserConfig(ctx context.Context, userID GitHubID, cfg *UserConfig) error
}
