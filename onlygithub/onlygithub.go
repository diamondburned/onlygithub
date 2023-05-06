package onlygithub

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/rs/xid"
)

// ErrorResponse is a response that contains an error message.
type ErrorResponse struct {
	Message string `json:"message"`
}

// GitHubID is a GitHub ID for a resource originating from GitHub.
type GitHubID string

// ID is an ID for our resource.
type ID = xid.ID

// GenerateID generates a new ID.
func GenerateID() ID {
	return xid.New()
}

// Cents is a monetary value in cents.
type Cents int64

// String returns the string representation of the monetary value with a dollar
// sign and two decimal places.
func (c Cents) String() string {
	dollars := c / 100
	cents := c % 100
	return fmt.Sprintf("$%d.%02d", dollars, cents)
}

// UnmarshalText implements encoding.TextUnmarshaler.
func (c *Cents) UnmarshalText(text []byte) error {
	str := string(text)

	if strings.HasPrefix(str, "$") {
		_, err := fmt.Sscanf(str, "$%d.%02d", &c, &c)
		if err != nil {
			return fmt.Errorf("invalid dollars: %s", str)
		}
		return nil
	}

	v, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return fmt.Errorf("invalid cents: %s", str)
	}

	*c = Cents(v)
	return nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (c *Cents) UnmarshalJSON(data []byte) error {
	if bytes.HasPrefix(data, []byte{'"'}) {
		var s string
		if err := json.Unmarshal(data, &s); err != nil {
			return err
		}
		return c.UnmarshalText([]byte(s))
	}

	var v int64
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	*c = Cents(v)
	return nil
}

// Tier describes a tier of service. Tiers are ordered from lowest to highest
// using the index of each tier in the Tiers slice.
type Tier struct {
	ID GitHubID `json:"id"`
	// Name is the name of the tier.
	Name string `json:"name"`
	// Price is the price of the tier in cents.
	Price Cents `json:"price"`
	// Description is the description of the tier.
	Description template.HTML `json:"description,omitempty"`
	// IsOneTime is true if the tier is a one-time payment. Otherwise, it is a
	// monthly subscription.
	IsOneTime bool `json:"isOneTime,omitempty"`
	// IsCustomAmount is true if the tier is a custom amount. Otherwise, it is
	// predefined. Note that for custom amounts, the actual tier is derived from
	// the lowest predefined tier that is greater than the custom amount.
	IsCustomAmount bool `json:"isCustomAmount,omitempty"`
}

// Tiers is a list of Tier. Tiers should be ordered from lowest to highest.
// Tiers implements sort.Interface, which sorts from lowest to highest by
// default.
type Tiers []Tier

// SelectFromPrice returns the highest tier that is less than or equal to the
// given price.
func (t Tiers) SelectFromPrice(price Cents) (*Tier, bool) {
	var tier *Tier
	sort.Search(len(t), func(i int) bool {
		tier = &t[i]
		return tier.Price > price
	})
	if tier == nil {
		return nil, false
	}
	return tier, true
}

func (t Tiers) Len() int           { return len(t) }
func (t Tiers) Less(i, j int) bool { return t[i].Price < t[j].Price }
func (t Tiers) Swap(i, j int)      { t[i], t[j] = t[j], t[i] }

// User is a GitHub user. It may technically be an organization on GitHub's end.
type User struct {
	// ID is the ID of the user on GitHub.
	ID GitHubID `json:"id"`
	// Name is the name of the user.
	Name string `json:"name"`
	// JoinedAt is the time the user joined this service.
	JoinedAt time.Time `json:"joinedAt"`
	// Email is the email of the user. It is not guaranteed to be set.
	// GitHub also has an API to send an email to a sponsor group, so that
	// should be used instead.
	Email string `json:"email,omitempty"`
	// Nickname is the nickname of the user.
	Nickname string `json:"nickname,omitempty"`
	// AvatarURL is the URL of the avatar of the user.
	AvatarURL string `json:"avatarURL,omitempty"`
	// Sponsorship is the sponsorship of the user. If null, the user is not
	// currently sponsoring.
	Sponsorship *Sponsorship `json:"sponsorship,omitempty"`
}

// Sponsorship describes a sponsorship of a user.
type Sponsorship struct {
	// Tier is the tier the user is currently on. If null, the user is not
	// currently on any tier.
	Tier Tier `json:"tier"`
	// SponsoredAt is the time the user was sponsored.
	SponsoredAt time.Time `json:"sponsoredAt"`
}

// Visibility is the visibility of a piece of content.
type Visibility string

const (
	// NotVisible means the content is not visible to anyone except the author.
	// This is a good default for drafts, so it's also the zero value.
	NotVisible Visibility = ""
	// VisibleToSponsors means the content is only visible to users who have
	// logged in and paid some amount.
	VisibleToSponsors Visibility = "sponsor"
	// VisibleToPrivate means the content is only visible to users who have
	// logged in but not necessarily paid.
	VisibleToPrivate Visibility = "private"
	// VisibleToPublic means the content is visible to everyone, including
	// users who have not logged in. This automatically implies that links
	// leading to the content are public.
	VisibleToPublic Visibility = "public"
)

// Asset is the metadata of an asset. It does not contain the actual asset data.
//
// An asset can be a piece of content, a comment, or a reaction. Each asset
// comes with a visibility and a minimum cost. The visibility determines who can
// access the asset. The minimum cost determines how much the user has to pay to
// access the asset. The actual cost may be higher than the minimum cost if the
// user chooses to pay more.
//
// An asset is also mutable and can be updated. The updated time is stored in
// LastUpdated.
type Asset struct {
	// ID is the ID of the asset.
	ID ID `json:"id"`
	// Visibility is the visibility of the asset.
	Visibility Visibility `json:"visibility"`
	// MinimumCost is the minimum cost to access the content.
	// The cost is derived from the tier; the user never pays this amount unless
	// the tier with this exact amount. This allows us to change the tier
	// structure without breaking existing content.
	//
	// If the content is free, MinimumCost is 0.
	MinimumCost Cents `json:"minimumCost"`
	// LastUpdated is the time the content was last updated.
	LastUpdated *time.Time `json:"updatedAt,omitempty"`
}

// Post is a single post. It is also an asset.
type Post struct {
	Asset
	// Markdown is the Markdown content.
	Markdown string `json:"markdown"`
	// Assets is a list of assets within the content. It acts as a reference to
	// the actual assets, preventing them to be garbage collected.
	Assets []ImageAsset `json:"assets,omitempty"`
	// AllowComments is true if comments are allowed on the content.
	AllowComments bool `json:"allowComments,omitempty"`
	// AllowReactions is true if reactions are allowed on the content.
	// Reactions will be allowed for comments as well.
	AllowReactions bool `json:"allowReactions,omitempty"`
}

// ImageAsset is an image asset. Image assets can be referenced by a special URL
// scheme: onlygithub://image/ID/Filename. The ID is the ID of the content, and
// the Filename is the filename of the asset.
type ImageAsset struct {
	Asset
	// Filename is the filename of the asset. It is used to reference the asset
	// in the content.
	Filename string `json:"filename"`
	// Width is the width of the image.
	Width int `json:"width"`
	// Height is the height of the image.
	Height int `json:"height"`
}

// Comment is a comment on a piece of content.
type Comment struct {
	ID ID `json:"id"`
	// Markdown is the Markdown content.
	Markdown string `json:"markdown"`
	// Author is the author of the comment.
	Author User `json:"author"`
}

// Reaction is a reaction to a piece of content. Each reaction belongs to a
// user.
type Reaction struct {
	Name string `json:"name"`
	// AssetID is the ID of the asset the reaction belongs to. If the asset is
	// zero, then name should be displayed instead. If there is an asset, then
	// it must be of type image.
	AssetID ID `json:"assetID"`
}
