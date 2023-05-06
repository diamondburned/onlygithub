package gh

//go:generate wget -N https://docs.github.com/public/schema.docs.graphql -O /tmp/github.graphqls
//go:generate genqlient

import (
	"context"
	"encoding/json"
	"html/template"

	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
	"libdb.so/onlygithub/onlygithub"

	genqlient "github.com/Khan/genqlient/graphql"
)

// Storage is a storage interface.
type Storage interface {
	// User returns the user with the given ID.
	User(id onlygithub.GitHubID) (*onlygithub.User, error)
}

// Paginator describes a paginator for any resource.
type Paginator[T any] interface {
	json.Marshaler
	json.Unmarshaler
	// Next returns the next page of results.
	Next() ([]T, error)
	// All returns all results.
	All() ([]T, error)
}

// Client wraps around the GitHub API client.
type Client struct {
	*githubv4.Client
	genqlient genqlient.Client
}

// NewClient creates a new GitHub API client.
// To make a new tokenSource, simply call config.TokenSource(ctx, oauthToken).
func NewClient(ctx context.Context, tokenSource oauth2.TokenSource) *Client {
	client := oauth2.NewClient(ctx, tokenSource)
	return &Client{
		Client:    githubv4.NewClient(client),
		genqlient: genqlient.NewClient("https://api.github.com/graphql", client),
	}
}

// Sponsors returns a paginator for fetching sponsors.
func (c *Client) Sponsors(ctx context.Context, limit int) Paginator[onlygithub.User] {
	return &paginator[sponsorsResponse, onlygithub.User]{
		client: c,
		ctx:    ctx,
		limit:  limit,
		queryFunc: func(ctx context.Context, client *Client, cursor string, limit int32) (*sponsorsResponse, error) {
			return sponsors(ctx, client.genqlient, cursor, limit)
		},
		mapFunc: func(resp *sponsorsResponse) (paginatedResource[onlygithub.User], error) {
			users := make([]onlygithub.User, 0, len(resp.Viewer.SponsorshipsAsMaintainer.Edges))
			for _, edge := range resp.Viewer.SponsorshipsAsMaintainer.Edges {
				node := edge.Node
				user := onlygithub.User{
					Sponsorship: &onlygithub.Sponsorship{
						Tier: onlygithub.Tier{
							ID:             onlygithub.GitHubID(node.Tier.Id),
							Name:           node.Tier.Name,
							Price:          onlygithub.Cents(node.Tier.MonthlyPriceInCents),
							Description:    template.HTML(node.Tier.DescriptionHTML),
							IsOneTime:      node.Tier.IsOneTime,
							IsCustomAmount: node.Tier.IsCustomAmount,
						},
						SponsoredAt: node.TierSelectedAt,
					},
				}

				switch sponsor := node.SponsorEntity.(type) {
				case *sponsorsViewerUserSponsorshipsAsMaintainerSponsorshipConnectionEdgesSponsorshipEdgeNodeSponsorshipSponsorEntityOrganization:
					user.ID = onlygithub.GitHubID(sponsor.Id)
					user.Name = sponsor.Login
					user.Nickname = sponsor.Name
					user.AvatarURL = sponsor.AvatarUrl
				case *sponsorsViewerUserSponsorshipsAsMaintainerSponsorshipConnectionEdgesSponsorshipEdgeNodeSponsorshipSponsorEntityUser:
					user.ID = onlygithub.GitHubID(sponsor.Id)
					user.Name = sponsor.Login
					user.Nickname = sponsor.Name
					user.AvatarURL = sponsor.AvatarUrl
				}

				users = append(users, user)
			}

			return paginatedResource[onlygithub.User]{
				Cursor:      resp.Viewer.SponsorshipsAsMaintainer.PageInfo.EndCursor,
				HasNextPage: resp.Viewer.SponsorshipsAsMaintainer.PageInfo.HasNextPage,
				Resources:   users,
			}, nil
		},
	}
}

func (c *Client) Tiers(ctx context.Context, limit int) Paginator[onlygithub.Tier] {
	return &paginator[tiersResponse, onlygithub.Tier]{
		client: c,
		ctx:    ctx,
		limit:  limit,
		queryFunc: func(ctx context.Context, client *Client, cursor string, limit int32) (*tiersResponse, error) {
			return tiers(ctx, client.genqlient, cursor, limit)
		},
		mapFunc: func(resp *tiersResponse) (paginatedResource[onlygithub.Tier], error) {
			tiers := make([]onlygithub.Tier, 0, len(resp.Viewer.SponsorsListing.Tiers.Edges))
			for _, edge := range resp.Viewer.SponsorsListing.Tiers.Edges {
				node := edge.Node
				tiers = append(tiers, onlygithub.Tier{
					ID:             onlygithub.GitHubID(node.Id),
					Name:           node.Name,
					Price:          onlygithub.Cents(node.MonthlyPriceInCents),
					Description:    template.HTML(node.DescriptionHTML),
					IsOneTime:      node.IsOneTime,
					IsCustomAmount: node.IsCustomAmount,
				})
			}
			return paginatedResource[onlygithub.Tier]{
				Cursor:      resp.Viewer.SponsorsListing.Tiers.PageInfo.EndCursor,
				HasNextPage: resp.Viewer.SponsorsListing.Tiers.PageInfo.HasNextPage,
				Resources:   tiers,
			}, nil
		},
	}
}

type paginator[RespT, ResourceT any] struct {
	client    *Client
	ctx       context.Context
	limit     int
	queryFunc func(ctx context.Context, client *Client, cursor string, limit int32) (*RespT, error)
	mapFunc   func(*RespT) (paginatedResource[ResourceT], error)

	Cursor  string `json:"cursor"`
	HasNext bool   `json:"hasNext"`
}

type paginatedResource[T any] struct {
	Cursor      string
	HasNextPage bool
	Resources   []T
}

func (p *paginator[RespT, ResourceT]) MarshalJSON() ([]byte, error) {
	type raw paginator[RespT, ResourceT]
	return json.Marshal((*raw)(p))
}

func (p *paginator[RespT, ResourceT]) UnmarshalJSON(data []byte) error {
	type raw paginator[RespT, ResourceT]
	return json.Unmarshal(data, (*raw)(p))
}

func (p *paginator[RespT, ResourceT]) Next() ([]ResourceT, error) {
	resp, err := p.queryFunc(p.ctx, p.client, p.Cursor, int32(p.limit))
	if err != nil {
		return nil, err
	}

	resources, err := p.mapFunc(resp)
	if err != nil {
		return nil, err
	}

	p.HasNext = resources.HasNextPage
	p.Cursor = resources.Cursor
	return resources.Resources, nil
}

func (p *paginator[RespT, ResourceT]) All() ([]ResourceT, error) {
	var all []ResourceT
	for p.HasNext {
		vs, err := p.Next()
		if err != nil {
			return nil, err
		}
		all = append(all, vs...)
	}
	return all, nil
}
