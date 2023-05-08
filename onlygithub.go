package onlygithub

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/rs/xid"
)

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
	// Username is the name of the user. It is known as `login` on GitHub's API.
	Username string `json:"username"`
	// JoinedAt is the time the user joined this service.
	JoinedAt time.Time `json:"joinedAt"`
	// Email is the email of the user. It is not guaranteed to be set.
	// GitHub also has an API to send an email to a sponsor group, so that
	// should be used instead.
	Email string `json:"email,omitempty"`
	// RealName is the nickname of the user.
	RealName string `json:"realName,omitempty"`
	// Pronouns is the pronouns of the user.
	Pronouns string `json:"pronouns,omitempty"`
	// AvatarURL is the URL of the avatar of the user.
	AvatarURL string `json:"avatarURL,omitempty"`
	// Sponsorship is the sponsorship of the user. If null, the user is not
	// currently sponsoring.
	Sponsorship *Sponsorship `json:"sponsorship,omitempty"`
}

// Sponsorship describes a sponsorship of a user.
type Sponsorship struct {
	// Price is the price of the sponsorship in cents. It may not correspond to
	// any tier.
	Price Cents `json:"price"`
	// StartedAt is the time the user was first sponsored.
	StartedAt time.Time `json:"firstSponsoredAt"`
	// RenewedAt is the time the user last renewed their sponsorship.
	RenewedAt time.Time `json:"renewedAt"`
	// IsOneTime is true if the tier is a one-time payment. Otherwise, it is a
	// monthly subscription.
	IsOneTime bool `json:"isOneTime,omitempty"`
	// IsCustomAmount is true if the tier is a custom amount. Otherwise, it is
	// predefined. Note that for custom amounts, the actual tier is derived from
	// the lowest predefined tier that is greater than the custom amount.
	IsCustomAmount bool `json:"isCustomAmount,omitempty"`
	// Tier is the tier the user is currently on. If null, the user is not
	// currently on any tier or is on a custom tier.
	Tier *Tier `json:"tier"`
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

// SiteConfig describes the configuration of the entire site. There's only one
// site-wide configuration.
type SiteConfig struct {
	// OwnerID is the ID of the owner of the site. This field should not be
	// modified; the user should manually set the owner ID in the database.
	OwnerID GitHubID `json:"ownerID,omitempty"`
	// BannerURL is the URL of the banner image that will be displayed on
	// the home page.
	BannerURL string `json:"bannerURL,omitempty"`
	// Description is the description of the site.
	Description template.HTML `json:"bio,omitempty"`
	// About is the about section of the site.
	About template.HTML `json:"about,omitempty"`
	// Socials is a list of social media accounts to display.
	Socials Socials `json:"socials,omitempty"`
	// AllowDMs controls whether or not users can send DMs to the
	// owner.
	AllowDMs bool `json:"allowDMs,omitempty"`
	// AllowComments controls whether or not users can comment on the
	// homepage.
	AllowComments bool `json:"allowComments,omitempty"`
	// TierWhitelist controls whether or not tiers must be explicitly
	// whitelisted in order to access sponsored content.
	TierWhitelist *TierWhitelist `json:"tierWhitelist,omitempty"`
	// CustomCSS is custom CSS that will be injected into the page.
	// Each page will be identified by a unique ID that can be used to
	// scope the CSS to that page.
	CustomCSS string `json:"customCSS,omitempty"`
}

// DefaultSiteConfig returns the default site-wide configuration.
func DefaultSiteConfig() *SiteConfig {
	return &SiteConfig{
		Socials: Socials{
			ShowGitHub: true,
		},
		AllowDMs:      true,
		AllowComments: true,
	}
}

// UserConfig describes the configuration of a user.
type UserConfig struct {
	// ShowComments controls whether or not comments are shown on the
	// homepage.
	ShowComments bool `json:"showComments,omitempty"`
	// Anonymous controls whether or not the user is anonymous. If true, then
	// user information will not be shown to other users.
	Anonymous bool `json:"anonymous,omitempty"`
}

// DefaultUserConfig returns the default configuration for a user.
func DefaultUserConfig() *UserConfig {
	return &UserConfig{
		ShowComments: true,
	}
}

// Socials is a list of social media accounts.
type Socials struct {
	// Twitter is the Twitter username.
	Twitter string `json:"twitter,omitempty"`
	// YouTube is the YouTube username.
	YouTube string `json:"youtube,omitempty"`
	// GitHub is the GitHub username.
	GitHub string `json:"github,omitempty"`
	// Twitch is the Twitch username.
	Twitch string `json:"twitch,omitempty"`
	// Discord is the Discord username.
	Discord string `json:"discord,omitempty"`
	// Instagram is the Instagram username.
	Instagram string `json:"instagram,omitempty"`
	// Matrix is the Matrix username.
	Matrix string `json:"matrix,omitempty"`
	// Reddit is the Reddit username.
	Reddit string `json:"reddit,omitempty"`
	// Facebook is the Facebook username.
	Facebook string `json:"facebook,omitempty"`
	// Mastodon is the Mastodon username.
	Mastodon string `json:"mastodon,omitempty"`
	// ShowGitHub controls whether or not the GitHub icon is shown.
	// The owner GitHub will be shown.
	ShowGitHub bool `json:"showGitHub,omitempty"`
}

// TierWhitelist is a list of tiers that are allowed to access sponsored
// content.
type TierWhitelist struct {
	// Requires is a list of requirements that a tier must satisfy in order to
	// be allowed to access sponsored content.
	Requires []TierRequirement
	// AllowCustom, if true, means that custom amounts are allowed to access
	// sponsored content as long as they are greater than any of the tiers
	// specified in Requires.
	AllowCustom bool
}

// TierRequirement is a description of a tier (or what a tier should have) that
// the author desires. At least one of these fields must be set.
//
// Some fields can have values surrounded by slashes to indicate regular
// expressions. For example, if the name is "/^Tier [0-9]+$/", then it will
// match any tier name that starts with "Tier " and ends with a number.
//
// The following fields can have regular expressions:
// - Name
// - Description
type TierRequirement struct {
	Name        string `json:"name,omitempty"`
	Price       Cents  `json:"price,omitempty"`
	Description string `json:"description,omitempty"`
	// AllowHigher means that tiers with a higher price than the one specified
	// are allowed to access sponsored content.
	AllowHigher bool `json:"allowHigher,omitempty"`
}

// Validate validates the TierRequirement.
func (r TierRequirement) Validate() error {
	if r.Name == "" && r.Price == 0 && r.Description == "" {
		return fmt.Errorf("at least one of name, price, or description must be set")
	}
	if isRegex(r.Name) {
		_, err := regexp.Compile(trimRegex(r.Name))
		if err != nil {
			return errors.Wrap(err, "invalid name regex")
		}
	}
	if isRegex(r.Description) {
		_, err := regexp.Compile(trimRegex(r.Name))
		if err != nil {
			return errors.Wrap(err, "invalid description regex")
		}
	}
	return nil
}

// Matches returns true if the given tier matches the requirement.
func (r TierRequirement) Matches(t Tier) bool {
	return false ||
		matches(r.Name, t.Name) ||
		matches(r.Description, string(t.Description)) ||
		(r.Price != 0 && t.Price == r.Price) ||
		(r.Price != 0 && r.AllowHigher && r.Price >= r.Price)
}

func isRegex(str string) bool {
	return strings.HasPrefix(str, "/") && strings.HasSuffix(str, "/")
}

func trimRegex(str string) string {
	return strings.TrimPrefix(strings.TrimSuffix(str, "/"), "/")
}

func matches(maybeRegex, str string) bool {
	if maybeRegex == "" {
		return false
	}
	if isRegex(maybeRegex) {
		re := regexp.MustCompile(trimRegex(maybeRegex))
		return re.MatchString(str)
	}
	return maybeRegex == str
}
