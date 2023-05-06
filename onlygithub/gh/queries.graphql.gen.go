// Code generated by github.com/Khan/genqlient, DO NOT EDIT.

package gh

import (
	"context"
	"encoding/json"
	"fmt"
	"template"
	"time"

	"github.com/Khan/genqlient/graphql"
)

// The privacy of a sponsorship
type SponsorshipPrivacy string

const (
	// Private
	SponsorshipPrivacyPrivate SponsorshipPrivacy = "PRIVATE"
	// Public
	SponsorshipPrivacyPublic SponsorshipPrivacy = "PUBLIC"
)

// __sponsorsInput is used internally by genqlient
type __sponsorsInput struct {
	EndCursor string `json:"endCursor"`
	Limit     int32  `json:"limit"`
}

// GetEndCursor returns __sponsorsInput.EndCursor, and is useful for accessing the field via an interface.
func (v *__sponsorsInput) GetEndCursor() string { return v.EndCursor }

// GetLimit returns __sponsorsInput.Limit, and is useful for accessing the field via an interface.
func (v *__sponsorsInput) GetLimit() int32 { return v.Limit }

// __tiersInput is used internally by genqlient
type __tiersInput struct {
	EndCursor string `json:"endCursor"`
	Limit     int32  `json:"limit"`
}

// GetEndCursor returns __tiersInput.EndCursor, and is useful for accessing the field via an interface.
func (v *__tiersInput) GetEndCursor() string { return v.EndCursor }

// GetLimit returns __tiersInput.Limit, and is useful for accessing the field via an interface.
func (v *__tiersInput) GetLimit() int32 { return v.Limit }

// sponsorsResponse is returned by sponsors on success.
type sponsorsResponse struct {
	// The currently authenticated user.
	Viewer sponsorsViewerUser `json:"viewer"`
}

// GetViewer returns sponsorsResponse.Viewer, and is useful for accessing the field via an interface.
func (v *sponsorsResponse) GetViewer() sponsorsViewerUser { return v.Viewer }

// sponsorsViewerUser includes the requested fields of the GraphQL type User.
// The GraphQL type's documentation follows.
//
// A user is an individual's account on GitHub that owns repositories and can make new content.
type sponsorsViewerUser struct {
	// The sponsorships where this user or organization is the maintainer receiving the funds.
	SponsorshipsAsMaintainer sponsorsViewerUserSponsorshipsAsMaintainerSponsorshipConnection `json:"sponsorshipsAsMaintainer"`
}

// GetSponsorshipsAsMaintainer returns sponsorsViewerUser.SponsorshipsAsMaintainer, and is useful for accessing the field via an interface.
func (v *sponsorsViewerUser) GetSponsorshipsAsMaintainer() sponsorsViewerUserSponsorshipsAsMaintainerSponsorshipConnection {
	return v.SponsorshipsAsMaintainer
}

// sponsorsViewerUserSponsorshipsAsMaintainerSponsorshipConnection includes the requested fields of the GraphQL type SponsorshipConnection.
// The GraphQL type's documentation follows.
//
// The connection type for Sponsorship.
type sponsorsViewerUserSponsorshipsAsMaintainerSponsorshipConnection struct {
	// A list of edges.
	Edges []sponsorsViewerUserSponsorshipsAsMaintainerSponsorshipConnectionEdgesSponsorshipEdge `json:"edges"`
	// Information to aid in pagination.
	PageInfo sponsorsViewerUserSponsorshipsAsMaintainerSponsorshipConnectionPageInfo `json:"pageInfo"`
}

// GetEdges returns sponsorsViewerUserSponsorshipsAsMaintainerSponsorshipConnection.Edges, and is useful for accessing the field via an interface.
func (v *sponsorsViewerUserSponsorshipsAsMaintainerSponsorshipConnection) GetEdges() []sponsorsViewerUserSponsorshipsAsMaintainerSponsorshipConnectionEdgesSponsorshipEdge {
	return v.Edges
}

// GetPageInfo returns sponsorsViewerUserSponsorshipsAsMaintainerSponsorshipConnection.PageInfo, and is useful for accessing the field via an interface.
func (v *sponsorsViewerUserSponsorshipsAsMaintainerSponsorshipConnection) GetPageInfo() sponsorsViewerUserSponsorshipsAsMaintainerSponsorshipConnectionPageInfo {
	return v.PageInfo
}

// sponsorsViewerUserSponsorshipsAsMaintainerSponsorshipConnectionEdgesSponsorshipEdge includes the requested fields of the GraphQL type SponsorshipEdge.
// The GraphQL type's documentation follows.
//
// An edge in a connection.
type sponsorsViewerUserSponsorshipsAsMaintainerSponsorshipConnectionEdgesSponsorshipEdge struct {
	// The item at the end of the edge.
	Node sponsorsViewerUserSponsorshipsAsMaintainerSponsorshipConnectionEdgesSponsorshipEdgeNodeSponsorship `json:"node"`
}

// GetNode returns sponsorsViewerUserSponsorshipsAsMaintainerSponsorshipConnectionEdgesSponsorshipEdge.Node, and is useful for accessing the field via an interface.
func (v *sponsorsViewerUserSponsorshipsAsMaintainerSponsorshipConnectionEdgesSponsorshipEdge) GetNode() sponsorsViewerUserSponsorshipsAsMaintainerSponsorshipConnectionEdgesSponsorshipEdgeNodeSponsorship {
	return v.Node
}

// sponsorsViewerUserSponsorshipsAsMaintainerSponsorshipConnectionEdgesSponsorshipEdgeNodeSponsorship includes the requested fields of the GraphQL type Sponsorship.
// The GraphQL type's documentation follows.
//
// A sponsorship relationship between a sponsor and a maintainer
type sponsorsViewerUserSponsorshipsAsMaintainerSponsorshipConnectionEdgesSponsorshipEdgeNodeSponsorship struct {
	Id string `json:"id"`
	// The user or organization that is sponsoring, if you have permission to view them.
	SponsorEntity sponsorsViewerUserSponsorshipsAsMaintainerSponsorshipConnectionEdgesSponsorshipEdgeNodeSponsorshipSponsorEntitySponsor `json:"-"`
	// The privacy level for this sponsorship.
	PrivacyLevel SponsorshipPrivacy `json:"privacyLevel"`
	// Whether this sponsorship represents a one-time payment versus a recurring sponsorship.
	IsOneTimePayment bool `json:"isOneTimePayment"`
	// Whether the sponsorship is active. False implies the sponsor is a past sponsor
	// of the maintainer, while true implies they are a current sponsor.
	IsActive bool `json:"isActive"`
	// Identifies the date and time when the object was created.
	CreatedAt time.Time `json:"createdAt"`
	// Identifies the date and time when the current tier was chosen for this sponsorship.
	TierSelectedAt time.Time `json:"tierSelectedAt"`
	// The associated sponsorship tier
	Tier sponsorsViewerUserSponsorshipsAsMaintainerSponsorshipConnectionEdgesSponsorshipEdgeNodeSponsorshipTierSponsorsTier `json:"tier"`
}

// GetId returns sponsorsViewerUserSponsorshipsAsMaintainerSponsorshipConnectionEdgesSponsorshipEdgeNodeSponsorship.Id, and is useful for accessing the field via an interface.
func (v *sponsorsViewerUserSponsorshipsAsMaintainerSponsorshipConnectionEdgesSponsorshipEdgeNodeSponsorship) GetId() string {
	return v.Id
}

// GetSponsorEntity returns sponsorsViewerUserSponsorshipsAsMaintainerSponsorshipConnectionEdgesSponsorshipEdgeNodeSponsorship.SponsorEntity, and is useful for accessing the field via an interface.
func (v *sponsorsViewerUserSponsorshipsAsMaintainerSponsorshipConnectionEdgesSponsorshipEdgeNodeSponsorship) GetSponsorEntity() sponsorsViewerUserSponsorshipsAsMaintainerSponsorshipConnectionEdgesSponsorshipEdgeNodeSponsorshipSponsorEntitySponsor {
	return v.SponsorEntity
}

// GetPrivacyLevel returns sponsorsViewerUserSponsorshipsAsMaintainerSponsorshipConnectionEdgesSponsorshipEdgeNodeSponsorship.PrivacyLevel, and is useful for accessing the field via an interface.
func (v *sponsorsViewerUserSponsorshipsAsMaintainerSponsorshipConnectionEdgesSponsorshipEdgeNodeSponsorship) GetPrivacyLevel() SponsorshipPrivacy {
	return v.PrivacyLevel
}

// GetIsOneTimePayment returns sponsorsViewerUserSponsorshipsAsMaintainerSponsorshipConnectionEdgesSponsorshipEdgeNodeSponsorship.IsOneTimePayment, and is useful for accessing the field via an interface.
func (v *sponsorsViewerUserSponsorshipsAsMaintainerSponsorshipConnectionEdgesSponsorshipEdgeNodeSponsorship) GetIsOneTimePayment() bool {
	return v.IsOneTimePayment
}

// GetIsActive returns sponsorsViewerUserSponsorshipsAsMaintainerSponsorshipConnectionEdgesSponsorshipEdgeNodeSponsorship.IsActive, and is useful for accessing the field via an interface.
func (v *sponsorsViewerUserSponsorshipsAsMaintainerSponsorshipConnectionEdgesSponsorshipEdgeNodeSponsorship) GetIsActive() bool {
	return v.IsActive
}

// GetCreatedAt returns sponsorsViewerUserSponsorshipsAsMaintainerSponsorshipConnectionEdgesSponsorshipEdgeNodeSponsorship.CreatedAt, and is useful for accessing the field via an interface.
func (v *sponsorsViewerUserSponsorshipsAsMaintainerSponsorshipConnectionEdgesSponsorshipEdgeNodeSponsorship) GetCreatedAt() time.Time {
	return v.CreatedAt
}

// GetTierSelectedAt returns sponsorsViewerUserSponsorshipsAsMaintainerSponsorshipConnectionEdgesSponsorshipEdgeNodeSponsorship.TierSelectedAt, and is useful for accessing the field via an interface.
func (v *sponsorsViewerUserSponsorshipsAsMaintainerSponsorshipConnectionEdgesSponsorshipEdgeNodeSponsorship) GetTierSelectedAt() time.Time {
	return v.TierSelectedAt
}

// GetTier returns sponsorsViewerUserSponsorshipsAsMaintainerSponsorshipConnectionEdgesSponsorshipEdgeNodeSponsorship.Tier, and is useful for accessing the field via an interface.
func (v *sponsorsViewerUserSponsorshipsAsMaintainerSponsorshipConnectionEdgesSponsorshipEdgeNodeSponsorship) GetTier() sponsorsViewerUserSponsorshipsAsMaintainerSponsorshipConnectionEdgesSponsorshipEdgeNodeSponsorshipTierSponsorsTier {
	return v.Tier
}

func (v *sponsorsViewerUserSponsorshipsAsMaintainerSponsorshipConnectionEdgesSponsorshipEdgeNodeSponsorship) UnmarshalJSON(b []byte) error {

	if string(b) == "null" {
		return nil
	}

	var firstPass struct {
		*sponsorsViewerUserSponsorshipsAsMaintainerSponsorshipConnectionEdgesSponsorshipEdgeNodeSponsorship
		SponsorEntity json.RawMessage `json:"sponsorEntity"`
		graphql.NoUnmarshalJSON
	}
	firstPass.sponsorsViewerUserSponsorshipsAsMaintainerSponsorshipConnectionEdgesSponsorshipEdgeNodeSponsorship = v

	err := json.Unmarshal(b, &firstPass)
	if err != nil {
		return err
	}

	{
		dst := &v.SponsorEntity
		src := firstPass.SponsorEntity
		if len(src) != 0 && string(src) != "null" {
			err = __unmarshalsponsorsViewerUserSponsorshipsAsMaintainerSponsorshipConnectionEdgesSponsorshipEdgeNodeSponsorshipSponsorEntitySponsor(
				src, dst)
			if err != nil {
				return fmt.Errorf(
					"unable to unmarshal sponsorsViewerUserSponsorshipsAsMaintainerSponsorshipConnectionEdgesSponsorshipEdgeNodeSponsorship.SponsorEntity: %w", err)
			}
		}
	}
	return nil
}

type __premarshalsponsorsViewerUserSponsorshipsAsMaintainerSponsorshipConnectionEdgesSponsorshipEdgeNodeSponsorship struct {
	Id string `json:"id"`

	SponsorEntity json.RawMessage `json:"sponsorEntity"`

	PrivacyLevel SponsorshipPrivacy `json:"privacyLevel"`

	IsOneTimePayment bool `json:"isOneTimePayment"`

	IsActive bool `json:"isActive"`

	CreatedAt time.Time `json:"createdAt"`

	TierSelectedAt time.Time `json:"tierSelectedAt"`

	Tier sponsorsViewerUserSponsorshipsAsMaintainerSponsorshipConnectionEdgesSponsorshipEdgeNodeSponsorshipTierSponsorsTier `json:"tier"`
}

func (v *sponsorsViewerUserSponsorshipsAsMaintainerSponsorshipConnectionEdgesSponsorshipEdgeNodeSponsorship) MarshalJSON() ([]byte, error) {
	premarshaled, err := v.__premarshalJSON()
	if err != nil {
		return nil, err
	}
	return json.Marshal(premarshaled)
}

func (v *sponsorsViewerUserSponsorshipsAsMaintainerSponsorshipConnectionEdgesSponsorshipEdgeNodeSponsorship) __premarshalJSON() (*__premarshalsponsorsViewerUserSponsorshipsAsMaintainerSponsorshipConnectionEdgesSponsorshipEdgeNodeSponsorship, error) {
	var retval __premarshalsponsorsViewerUserSponsorshipsAsMaintainerSponsorshipConnectionEdgesSponsorshipEdgeNodeSponsorship

	retval.Id = v.Id
	{

		dst := &retval.SponsorEntity
		src := v.SponsorEntity
		var err error
		*dst, err = __marshalsponsorsViewerUserSponsorshipsAsMaintainerSponsorshipConnectionEdgesSponsorshipEdgeNodeSponsorshipSponsorEntitySponsor(
			&src)
		if err != nil {
			return nil, fmt.Errorf(
				"unable to marshal sponsorsViewerUserSponsorshipsAsMaintainerSponsorshipConnectionEdgesSponsorshipEdgeNodeSponsorship.SponsorEntity: %w", err)
		}
	}
	retval.PrivacyLevel = v.PrivacyLevel
	retval.IsOneTimePayment = v.IsOneTimePayment
	retval.IsActive = v.IsActive
	retval.CreatedAt = v.CreatedAt
	retval.TierSelectedAt = v.TierSelectedAt
	retval.Tier = v.Tier
	return &retval, nil
}

// sponsorsViewerUserSponsorshipsAsMaintainerSponsorshipConnectionEdgesSponsorshipEdgeNodeSponsorshipSponsorEntityOrganization includes the requested fields of the GraphQL type Organization.
// The GraphQL type's documentation follows.
//
// An account on GitHub, with one or more owners, that has repositories, members and teams.
type sponsorsViewerUserSponsorshipsAsMaintainerSponsorshipConnectionEdgesSponsorshipEdgeNodeSponsorshipSponsorEntityOrganization struct {
	Typename string `json:"__typename"`
	Id       string `json:"id"`
	// The organization's login name.
	Login string `json:"login"`
	// The organization's public profile name.
	Name string `json:"name"`
	// A URL pointing to the organization's public avatar.
	AvatarUrl string `json:"avatarUrl"`
}

// GetTypename returns sponsorsViewerUserSponsorshipsAsMaintainerSponsorshipConnectionEdgesSponsorshipEdgeNodeSponsorshipSponsorEntityOrganization.Typename, and is useful for accessing the field via an interface.
func (v *sponsorsViewerUserSponsorshipsAsMaintainerSponsorshipConnectionEdgesSponsorshipEdgeNodeSponsorshipSponsorEntityOrganization) GetTypename() string {
	return v.Typename
}

// GetId returns sponsorsViewerUserSponsorshipsAsMaintainerSponsorshipConnectionEdgesSponsorshipEdgeNodeSponsorshipSponsorEntityOrganization.Id, and is useful for accessing the field via an interface.
func (v *sponsorsViewerUserSponsorshipsAsMaintainerSponsorshipConnectionEdgesSponsorshipEdgeNodeSponsorshipSponsorEntityOrganization) GetId() string {
	return v.Id
}

// GetLogin returns sponsorsViewerUserSponsorshipsAsMaintainerSponsorshipConnectionEdgesSponsorshipEdgeNodeSponsorshipSponsorEntityOrganization.Login, and is useful for accessing the field via an interface.
func (v *sponsorsViewerUserSponsorshipsAsMaintainerSponsorshipConnectionEdgesSponsorshipEdgeNodeSponsorshipSponsorEntityOrganization) GetLogin() string {
	return v.Login
}

// GetName returns sponsorsViewerUserSponsorshipsAsMaintainerSponsorshipConnectionEdgesSponsorshipEdgeNodeSponsorshipSponsorEntityOrganization.Name, and is useful for accessing the field via an interface.
func (v *sponsorsViewerUserSponsorshipsAsMaintainerSponsorshipConnectionEdgesSponsorshipEdgeNodeSponsorshipSponsorEntityOrganization) GetName() string {
	return v.Name
}

// GetAvatarUrl returns sponsorsViewerUserSponsorshipsAsMaintainerSponsorshipConnectionEdgesSponsorshipEdgeNodeSponsorshipSponsorEntityOrganization.AvatarUrl, and is useful for accessing the field via an interface.
func (v *sponsorsViewerUserSponsorshipsAsMaintainerSponsorshipConnectionEdgesSponsorshipEdgeNodeSponsorshipSponsorEntityOrganization) GetAvatarUrl() string {
	return v.AvatarUrl
}

// sponsorsViewerUserSponsorshipsAsMaintainerSponsorshipConnectionEdgesSponsorshipEdgeNodeSponsorshipSponsorEntitySponsor includes the requested fields of the GraphQL interface Sponsor.
//
// sponsorsViewerUserSponsorshipsAsMaintainerSponsorshipConnectionEdgesSponsorshipEdgeNodeSponsorshipSponsorEntitySponsor is implemented by the following types:
// sponsorsViewerUserSponsorshipsAsMaintainerSponsorshipConnectionEdgesSponsorshipEdgeNodeSponsorshipSponsorEntityOrganization
// sponsorsViewerUserSponsorshipsAsMaintainerSponsorshipConnectionEdgesSponsorshipEdgeNodeSponsorshipSponsorEntityUser
// The GraphQL type's documentation follows.
//
// Entities that can sponsor others via GitHub Sponsors
type sponsorsViewerUserSponsorshipsAsMaintainerSponsorshipConnectionEdgesSponsorshipEdgeNodeSponsorshipSponsorEntitySponsor interface {
	implementsGraphQLInterfacesponsorsViewerUserSponsorshipsAsMaintainerSponsorshipConnectionEdgesSponsorshipEdgeNodeSponsorshipSponsorEntitySponsor()
	// GetTypename returns the receiver's concrete GraphQL type-name (see interface doc for possible values).
	GetTypename() string
}

func (v *sponsorsViewerUserSponsorshipsAsMaintainerSponsorshipConnectionEdgesSponsorshipEdgeNodeSponsorshipSponsorEntityOrganization) implementsGraphQLInterfacesponsorsViewerUserSponsorshipsAsMaintainerSponsorshipConnectionEdgesSponsorshipEdgeNodeSponsorshipSponsorEntitySponsor() {
}
func (v *sponsorsViewerUserSponsorshipsAsMaintainerSponsorshipConnectionEdgesSponsorshipEdgeNodeSponsorshipSponsorEntityUser) implementsGraphQLInterfacesponsorsViewerUserSponsorshipsAsMaintainerSponsorshipConnectionEdgesSponsorshipEdgeNodeSponsorshipSponsorEntitySponsor() {
}

func __unmarshalsponsorsViewerUserSponsorshipsAsMaintainerSponsorshipConnectionEdgesSponsorshipEdgeNodeSponsorshipSponsorEntitySponsor(b []byte, v *sponsorsViewerUserSponsorshipsAsMaintainerSponsorshipConnectionEdgesSponsorshipEdgeNodeSponsorshipSponsorEntitySponsor) error {
	if string(b) == "null" {
		return nil
	}

	var tn struct {
		TypeName string `json:"__typename"`
	}
	err := json.Unmarshal(b, &tn)
	if err != nil {
		return err
	}

	switch tn.TypeName {
	case "Organization":
		*v = new(sponsorsViewerUserSponsorshipsAsMaintainerSponsorshipConnectionEdgesSponsorshipEdgeNodeSponsorshipSponsorEntityOrganization)
		return json.Unmarshal(b, *v)
	case "User":
		*v = new(sponsorsViewerUserSponsorshipsAsMaintainerSponsorshipConnectionEdgesSponsorshipEdgeNodeSponsorshipSponsorEntityUser)
		return json.Unmarshal(b, *v)
	case "":
		return fmt.Errorf(
			"response was missing Sponsor.__typename")
	default:
		return fmt.Errorf(
			`unexpected concrete type for sponsorsViewerUserSponsorshipsAsMaintainerSponsorshipConnectionEdgesSponsorshipEdgeNodeSponsorshipSponsorEntitySponsor: "%v"`, tn.TypeName)
	}
}

func __marshalsponsorsViewerUserSponsorshipsAsMaintainerSponsorshipConnectionEdgesSponsorshipEdgeNodeSponsorshipSponsorEntitySponsor(v *sponsorsViewerUserSponsorshipsAsMaintainerSponsorshipConnectionEdgesSponsorshipEdgeNodeSponsorshipSponsorEntitySponsor) ([]byte, error) {

	var typename string
	switch v := (*v).(type) {
	case *sponsorsViewerUserSponsorshipsAsMaintainerSponsorshipConnectionEdgesSponsorshipEdgeNodeSponsorshipSponsorEntityOrganization:
		typename = "Organization"

		result := struct {
			TypeName string `json:"__typename"`
			*sponsorsViewerUserSponsorshipsAsMaintainerSponsorshipConnectionEdgesSponsorshipEdgeNodeSponsorshipSponsorEntityOrganization
		}{typename, v}
		return json.Marshal(result)
	case *sponsorsViewerUserSponsorshipsAsMaintainerSponsorshipConnectionEdgesSponsorshipEdgeNodeSponsorshipSponsorEntityUser:
		typename = "User"

		result := struct {
			TypeName string `json:"__typename"`
			*sponsorsViewerUserSponsorshipsAsMaintainerSponsorshipConnectionEdgesSponsorshipEdgeNodeSponsorshipSponsorEntityUser
		}{typename, v}
		return json.Marshal(result)
	case nil:
		return []byte("null"), nil
	default:
		return nil, fmt.Errorf(
			`unexpected concrete type for sponsorsViewerUserSponsorshipsAsMaintainerSponsorshipConnectionEdgesSponsorshipEdgeNodeSponsorshipSponsorEntitySponsor: "%T"`, v)
	}
}

// sponsorsViewerUserSponsorshipsAsMaintainerSponsorshipConnectionEdgesSponsorshipEdgeNodeSponsorshipSponsorEntityUser includes the requested fields of the GraphQL type User.
// The GraphQL type's documentation follows.
//
// A user is an individual's account on GitHub that owns repositories and can make new content.
type sponsorsViewerUserSponsorshipsAsMaintainerSponsorshipConnectionEdgesSponsorshipEdgeNodeSponsorshipSponsorEntityUser struct {
	Typename string `json:"__typename"`
	Id       string `json:"id"`
	// The username used to login.
	Login string `json:"login"`
	// The user's public profile name.
	Name string `json:"name"`
	// A URL pointing to the user's public avatar.
	AvatarUrl string `json:"avatarUrl"`
}

// GetTypename returns sponsorsViewerUserSponsorshipsAsMaintainerSponsorshipConnectionEdgesSponsorshipEdgeNodeSponsorshipSponsorEntityUser.Typename, and is useful for accessing the field via an interface.
func (v *sponsorsViewerUserSponsorshipsAsMaintainerSponsorshipConnectionEdgesSponsorshipEdgeNodeSponsorshipSponsorEntityUser) GetTypename() string {
	return v.Typename
}

// GetId returns sponsorsViewerUserSponsorshipsAsMaintainerSponsorshipConnectionEdgesSponsorshipEdgeNodeSponsorshipSponsorEntityUser.Id, and is useful for accessing the field via an interface.
func (v *sponsorsViewerUserSponsorshipsAsMaintainerSponsorshipConnectionEdgesSponsorshipEdgeNodeSponsorshipSponsorEntityUser) GetId() string {
	return v.Id
}

// GetLogin returns sponsorsViewerUserSponsorshipsAsMaintainerSponsorshipConnectionEdgesSponsorshipEdgeNodeSponsorshipSponsorEntityUser.Login, and is useful for accessing the field via an interface.
func (v *sponsorsViewerUserSponsorshipsAsMaintainerSponsorshipConnectionEdgesSponsorshipEdgeNodeSponsorshipSponsorEntityUser) GetLogin() string {
	return v.Login
}

// GetName returns sponsorsViewerUserSponsorshipsAsMaintainerSponsorshipConnectionEdgesSponsorshipEdgeNodeSponsorshipSponsorEntityUser.Name, and is useful for accessing the field via an interface.
func (v *sponsorsViewerUserSponsorshipsAsMaintainerSponsorshipConnectionEdgesSponsorshipEdgeNodeSponsorshipSponsorEntityUser) GetName() string {
	return v.Name
}

// GetAvatarUrl returns sponsorsViewerUserSponsorshipsAsMaintainerSponsorshipConnectionEdgesSponsorshipEdgeNodeSponsorshipSponsorEntityUser.AvatarUrl, and is useful for accessing the field via an interface.
func (v *sponsorsViewerUserSponsorshipsAsMaintainerSponsorshipConnectionEdgesSponsorshipEdgeNodeSponsorshipSponsorEntityUser) GetAvatarUrl() string {
	return v.AvatarUrl
}

// sponsorsViewerUserSponsorshipsAsMaintainerSponsorshipConnectionEdgesSponsorshipEdgeNodeSponsorshipTierSponsorsTier includes the requested fields of the GraphQL type SponsorsTier.
// The GraphQL type's documentation follows.
//
// A GitHub Sponsors tier associated with a GitHub Sponsors listing.
type sponsorsViewerUserSponsorshipsAsMaintainerSponsorshipConnectionEdgesSponsorshipEdgeNodeSponsorshipTierSponsorsTier struct {
	Id string `json:"id"`
	// The name of the tier.
	Name string `json:"name"`
	// How much this tier costs per month in cents.
	MonthlyPriceInCents int32 `json:"monthlyPriceInCents"`
	// Whether this tier is only for use with one-time sponsorships.
	IsOneTime bool `json:"isOneTime"`
	// Whether this tier was chosen at checkout time by the sponsor rather than
	// defined ahead of time by the maintainer who manages the Sponsors listing.
	IsCustomAmount bool `json:"isCustomAmount"`
	// The description of the tier.
	Description string `json:"description"`
	// The tier description rendered to HTML
	DescriptionHTML template.HTML `json:"descriptionHTML"`
}

// GetId returns sponsorsViewerUserSponsorshipsAsMaintainerSponsorshipConnectionEdgesSponsorshipEdgeNodeSponsorshipTierSponsorsTier.Id, and is useful for accessing the field via an interface.
func (v *sponsorsViewerUserSponsorshipsAsMaintainerSponsorshipConnectionEdgesSponsorshipEdgeNodeSponsorshipTierSponsorsTier) GetId() string {
	return v.Id
}

// GetName returns sponsorsViewerUserSponsorshipsAsMaintainerSponsorshipConnectionEdgesSponsorshipEdgeNodeSponsorshipTierSponsorsTier.Name, and is useful for accessing the field via an interface.
func (v *sponsorsViewerUserSponsorshipsAsMaintainerSponsorshipConnectionEdgesSponsorshipEdgeNodeSponsorshipTierSponsorsTier) GetName() string {
	return v.Name
}

// GetMonthlyPriceInCents returns sponsorsViewerUserSponsorshipsAsMaintainerSponsorshipConnectionEdgesSponsorshipEdgeNodeSponsorshipTierSponsorsTier.MonthlyPriceInCents, and is useful for accessing the field via an interface.
func (v *sponsorsViewerUserSponsorshipsAsMaintainerSponsorshipConnectionEdgesSponsorshipEdgeNodeSponsorshipTierSponsorsTier) GetMonthlyPriceInCents() int32 {
	return v.MonthlyPriceInCents
}

// GetIsOneTime returns sponsorsViewerUserSponsorshipsAsMaintainerSponsorshipConnectionEdgesSponsorshipEdgeNodeSponsorshipTierSponsorsTier.IsOneTime, and is useful for accessing the field via an interface.
func (v *sponsorsViewerUserSponsorshipsAsMaintainerSponsorshipConnectionEdgesSponsorshipEdgeNodeSponsorshipTierSponsorsTier) GetIsOneTime() bool {
	return v.IsOneTime
}

// GetIsCustomAmount returns sponsorsViewerUserSponsorshipsAsMaintainerSponsorshipConnectionEdgesSponsorshipEdgeNodeSponsorshipTierSponsorsTier.IsCustomAmount, and is useful for accessing the field via an interface.
func (v *sponsorsViewerUserSponsorshipsAsMaintainerSponsorshipConnectionEdgesSponsorshipEdgeNodeSponsorshipTierSponsorsTier) GetIsCustomAmount() bool {
	return v.IsCustomAmount
}

// GetDescription returns sponsorsViewerUserSponsorshipsAsMaintainerSponsorshipConnectionEdgesSponsorshipEdgeNodeSponsorshipTierSponsorsTier.Description, and is useful for accessing the field via an interface.
func (v *sponsorsViewerUserSponsorshipsAsMaintainerSponsorshipConnectionEdgesSponsorshipEdgeNodeSponsorshipTierSponsorsTier) GetDescription() string {
	return v.Description
}

// GetDescriptionHTML returns sponsorsViewerUserSponsorshipsAsMaintainerSponsorshipConnectionEdgesSponsorshipEdgeNodeSponsorshipTierSponsorsTier.DescriptionHTML, and is useful for accessing the field via an interface.
func (v *sponsorsViewerUserSponsorshipsAsMaintainerSponsorshipConnectionEdgesSponsorshipEdgeNodeSponsorshipTierSponsorsTier) GetDescriptionHTML() template.HTML {
	return v.DescriptionHTML
}

// sponsorsViewerUserSponsorshipsAsMaintainerSponsorshipConnectionPageInfo includes the requested fields of the GraphQL type PageInfo.
// The GraphQL type's documentation follows.
//
// Information about pagination in a connection.
type sponsorsViewerUserSponsorshipsAsMaintainerSponsorshipConnectionPageInfo struct {
	// When paginating forwards, the cursor to continue.
	EndCursor string `json:"endCursor"`
	// When paginating forwards, are there more items?
	HasNextPage bool `json:"hasNextPage"`
}

// GetEndCursor returns sponsorsViewerUserSponsorshipsAsMaintainerSponsorshipConnectionPageInfo.EndCursor, and is useful for accessing the field via an interface.
func (v *sponsorsViewerUserSponsorshipsAsMaintainerSponsorshipConnectionPageInfo) GetEndCursor() string {
	return v.EndCursor
}

// GetHasNextPage returns sponsorsViewerUserSponsorshipsAsMaintainerSponsorshipConnectionPageInfo.HasNextPage, and is useful for accessing the field via an interface.
func (v *sponsorsViewerUserSponsorshipsAsMaintainerSponsorshipConnectionPageInfo) GetHasNextPage() bool {
	return v.HasNextPage
}

// tiersResponse is returned by tiers on success.
type tiersResponse struct {
	// The currently authenticated user.
	Viewer tiersViewerUser `json:"viewer"`
}

// GetViewer returns tiersResponse.Viewer, and is useful for accessing the field via an interface.
func (v *tiersResponse) GetViewer() tiersViewerUser { return v.Viewer }

// tiersViewerUser includes the requested fields of the GraphQL type User.
// The GraphQL type's documentation follows.
//
// A user is an individual's account on GitHub that owns repositories and can make new content.
type tiersViewerUser struct {
	// The GitHub Sponsors listing for this user or organization.
	SponsorsListing tiersViewerUserSponsorsListing `json:"sponsorsListing"`
}

// GetSponsorsListing returns tiersViewerUser.SponsorsListing, and is useful for accessing the field via an interface.
func (v *tiersViewerUser) GetSponsorsListing() tiersViewerUserSponsorsListing {
	return v.SponsorsListing
}

// tiersViewerUserSponsorsListing includes the requested fields of the GraphQL type SponsorsListing.
// The GraphQL type's documentation follows.
//
// A GitHub Sponsors listing.
type tiersViewerUserSponsorsListing struct {
	// The tiers for this GitHub Sponsors profile.
	Tiers tiersViewerUserSponsorsListingTiersSponsorsTierConnection `json:"tiers"`
}

// GetTiers returns tiersViewerUserSponsorsListing.Tiers, and is useful for accessing the field via an interface.
func (v *tiersViewerUserSponsorsListing) GetTiers() tiersViewerUserSponsorsListingTiersSponsorsTierConnection {
	return v.Tiers
}

// tiersViewerUserSponsorsListingTiersSponsorsTierConnection includes the requested fields of the GraphQL type SponsorsTierConnection.
// The GraphQL type's documentation follows.
//
// The connection type for SponsorsTier.
type tiersViewerUserSponsorsListingTiersSponsorsTierConnection struct {
	// A list of edges.
	Edges []tiersViewerUserSponsorsListingTiersSponsorsTierConnectionEdgesSponsorsTierEdge `json:"edges"`
	// Information to aid in pagination.
	PageInfo tiersViewerUserSponsorsListingTiersSponsorsTierConnectionPageInfo `json:"pageInfo"`
}

// GetEdges returns tiersViewerUserSponsorsListingTiersSponsorsTierConnection.Edges, and is useful for accessing the field via an interface.
func (v *tiersViewerUserSponsorsListingTiersSponsorsTierConnection) GetEdges() []tiersViewerUserSponsorsListingTiersSponsorsTierConnectionEdgesSponsorsTierEdge {
	return v.Edges
}

// GetPageInfo returns tiersViewerUserSponsorsListingTiersSponsorsTierConnection.PageInfo, and is useful for accessing the field via an interface.
func (v *tiersViewerUserSponsorsListingTiersSponsorsTierConnection) GetPageInfo() tiersViewerUserSponsorsListingTiersSponsorsTierConnectionPageInfo {
	return v.PageInfo
}

// tiersViewerUserSponsorsListingTiersSponsorsTierConnectionEdgesSponsorsTierEdge includes the requested fields of the GraphQL type SponsorsTierEdge.
// The GraphQL type's documentation follows.
//
// An edge in a connection.
type tiersViewerUserSponsorsListingTiersSponsorsTierConnectionEdgesSponsorsTierEdge struct {
	// The item at the end of the edge.
	Node tiersViewerUserSponsorsListingTiersSponsorsTierConnectionEdgesSponsorsTierEdgeNodeSponsorsTier `json:"node"`
}

// GetNode returns tiersViewerUserSponsorsListingTiersSponsorsTierConnectionEdgesSponsorsTierEdge.Node, and is useful for accessing the field via an interface.
func (v *tiersViewerUserSponsorsListingTiersSponsorsTierConnectionEdgesSponsorsTierEdge) GetNode() tiersViewerUserSponsorsListingTiersSponsorsTierConnectionEdgesSponsorsTierEdgeNodeSponsorsTier {
	return v.Node
}

// tiersViewerUserSponsorsListingTiersSponsorsTierConnectionEdgesSponsorsTierEdgeNodeSponsorsTier includes the requested fields of the GraphQL type SponsorsTier.
// The GraphQL type's documentation follows.
//
// A GitHub Sponsors tier associated with a GitHub Sponsors listing.
type tiersViewerUserSponsorsListingTiersSponsorsTierConnectionEdgesSponsorsTierEdgeNodeSponsorsTier struct {
	Id string `json:"id"`
	// The name of the tier.
	Name string `json:"name"`
	// How much this tier costs per month in cents.
	MonthlyPriceInCents int32 `json:"monthlyPriceInCents"`
	// Whether this tier is only for use with one-time sponsorships.
	IsOneTime bool `json:"isOneTime"`
	// Whether this tier was chosen at checkout time by the sponsor rather than
	// defined ahead of time by the maintainer who manages the Sponsors listing.
	IsCustomAmount bool `json:"isCustomAmount"`
	// The description of the tier.
	Description string `json:"description"`
	// The tier description rendered to HTML
	DescriptionHTML template.HTML `json:"descriptionHTML"`
}

// GetId returns tiersViewerUserSponsorsListingTiersSponsorsTierConnectionEdgesSponsorsTierEdgeNodeSponsorsTier.Id, and is useful for accessing the field via an interface.
func (v *tiersViewerUserSponsorsListingTiersSponsorsTierConnectionEdgesSponsorsTierEdgeNodeSponsorsTier) GetId() string {
	return v.Id
}

// GetName returns tiersViewerUserSponsorsListingTiersSponsorsTierConnectionEdgesSponsorsTierEdgeNodeSponsorsTier.Name, and is useful for accessing the field via an interface.
func (v *tiersViewerUserSponsorsListingTiersSponsorsTierConnectionEdgesSponsorsTierEdgeNodeSponsorsTier) GetName() string {
	return v.Name
}

// GetMonthlyPriceInCents returns tiersViewerUserSponsorsListingTiersSponsorsTierConnectionEdgesSponsorsTierEdgeNodeSponsorsTier.MonthlyPriceInCents, and is useful for accessing the field via an interface.
func (v *tiersViewerUserSponsorsListingTiersSponsorsTierConnectionEdgesSponsorsTierEdgeNodeSponsorsTier) GetMonthlyPriceInCents() int32 {
	return v.MonthlyPriceInCents
}

// GetIsOneTime returns tiersViewerUserSponsorsListingTiersSponsorsTierConnectionEdgesSponsorsTierEdgeNodeSponsorsTier.IsOneTime, and is useful for accessing the field via an interface.
func (v *tiersViewerUserSponsorsListingTiersSponsorsTierConnectionEdgesSponsorsTierEdgeNodeSponsorsTier) GetIsOneTime() bool {
	return v.IsOneTime
}

// GetIsCustomAmount returns tiersViewerUserSponsorsListingTiersSponsorsTierConnectionEdgesSponsorsTierEdgeNodeSponsorsTier.IsCustomAmount, and is useful for accessing the field via an interface.
func (v *tiersViewerUserSponsorsListingTiersSponsorsTierConnectionEdgesSponsorsTierEdgeNodeSponsorsTier) GetIsCustomAmount() bool {
	return v.IsCustomAmount
}

// GetDescription returns tiersViewerUserSponsorsListingTiersSponsorsTierConnectionEdgesSponsorsTierEdgeNodeSponsorsTier.Description, and is useful for accessing the field via an interface.
func (v *tiersViewerUserSponsorsListingTiersSponsorsTierConnectionEdgesSponsorsTierEdgeNodeSponsorsTier) GetDescription() string {
	return v.Description
}

// GetDescriptionHTML returns tiersViewerUserSponsorsListingTiersSponsorsTierConnectionEdgesSponsorsTierEdgeNodeSponsorsTier.DescriptionHTML, and is useful for accessing the field via an interface.
func (v *tiersViewerUserSponsorsListingTiersSponsorsTierConnectionEdgesSponsorsTierEdgeNodeSponsorsTier) GetDescriptionHTML() template.HTML {
	return v.DescriptionHTML
}

// tiersViewerUserSponsorsListingTiersSponsorsTierConnectionPageInfo includes the requested fields of the GraphQL type PageInfo.
// The GraphQL type's documentation follows.
//
// Information about pagination in a connection.
type tiersViewerUserSponsorsListingTiersSponsorsTierConnectionPageInfo struct {
	// When paginating forwards, the cursor to continue.
	EndCursor string `json:"endCursor"`
	// When paginating forwards, are there more items?
	HasNextPage bool `json:"hasNextPage"`
}

// GetEndCursor returns tiersViewerUserSponsorsListingTiersSponsorsTierConnectionPageInfo.EndCursor, and is useful for accessing the field via an interface.
func (v *tiersViewerUserSponsorsListingTiersSponsorsTierConnectionPageInfo) GetEndCursor() string {
	return v.EndCursor
}

// GetHasNextPage returns tiersViewerUserSponsorsListingTiersSponsorsTierConnectionPageInfo.HasNextPage, and is useful for accessing the field via an interface.
func (v *tiersViewerUserSponsorsListingTiersSponsorsTierConnectionPageInfo) GetHasNextPage() bool {
	return v.HasNextPage
}

// The query or mutation executed by sponsors.
const sponsors_Operation = `
query sponsors ($endCursor: String, $limit: Int = 100) {
	viewer {
		sponsorshipsAsMaintainer(first: $limit, after: $endCursor, includePrivate: true, activeOnly: true) {
			edges {
				node {
					id
					sponsorEntity {
						__typename
						... on User {
							id
							login
							name
							avatarUrl
						}
						... on Organization {
							id
							login
							name
							avatarUrl
						}
					}
					privacyLevel
					isOneTimePayment
					isActive
					createdAt
					tierSelectedAt
					tier {
						id
						name
						monthlyPriceInCents
						isOneTime
						isCustomAmount
						description
						descriptionHTML
					}
				}
			}
			pageInfo {
				endCursor
				hasNextPage
			}
		}
	}
}
`

func sponsors(
	ctx context.Context,
	client graphql.Client,
	endCursor string,
	limit int32,
) (*sponsorsResponse, error) {
	req := &graphql.Request{
		OpName: "sponsors",
		Query:  sponsors_Operation,
		Variables: &__sponsorsInput{
			EndCursor: endCursor,
			Limit:     limit,
		},
	}
	var err error

	var data sponsorsResponse
	resp := &graphql.Response{Data: &data}

	err = client.MakeRequest(
		ctx,
		req,
		resp,
	)

	return &data, err
}

// The query or mutation executed by tiers.
const tiers_Operation = `
query tiers ($endCursor: String, $limit: Int = 100) {
	viewer {
		sponsorsListing {
			tiers(first: $limit, after: $endCursor) {
				edges {
					node {
						id
						name
						monthlyPriceInCents
						isOneTime
						isCustomAmount
						description
						descriptionHTML
					}
				}
				pageInfo {
					endCursor
					hasNextPage
				}
			}
		}
	}
}
`

func tiers(
	ctx context.Context,
	client graphql.Client,
	endCursor string,
	limit int32,
) (*tiersResponse, error) {
	req := &graphql.Request{
		OpName: "tiers",
		Query:  tiers_Operation,
		Variables: &__tiersInput{
			EndCursor: endCursor,
			Limit:     limit,
		},
	}
	var err error

	var data tiersResponse
	resp := &graphql.Response{Data: &data}

	err = client.MakeRequest(
		ctx,
		req,
		resp,
	)

	return &data, err
}
