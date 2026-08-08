package metron

import (
	"context"
	"iter"
	"net/url"
	"time"

	"github.com/jyggen/go-metron/internal"
)

// Team is a team of characters.
type Team struct {
	ID                    int
	Name                  string
	Description           *string
	ImageURL              *url.URL
	Creators              []CreatorList
	Universes             []UniverseList
	ComicVineID           *int
	GrandComicsDatabaseID *int
	ResourceURL           url.URL
	Modified              time.Time
}

// TeamList is a team as it appears in list responses.
type TeamList struct {
	ID       int
	Name     string
	Modified time.Time
}

// TeamByID returns a team by its ID.
func (c *Client) TeamByID(ctx context.Context, id int) (*Team, error) {
	return byID(ctx, c, c.client.ApiTeamRetrieve, teamMapper, id)
}

// Teams returns an iterator over all teams.
func (c *Client) Teams(ctx context.Context, filters ...Filter) iter.Seq2[*TeamList, error] {
	params := &internal.ApiTeamListParams{}

	if err := applyFilters("Teams", params, filters); err != nil {
		return errIter[*TeamList](err)
	}

	return paginate[internal.PaginatedTeamListList](ctx, c, c.client.ApiTeamList, teamListMapper, params)
}

func teamMapper(in internal.TeamRead) (*Team, error) {
	if in.Id == nil {
		return nil, &MapError{Kind: "team", Field: "Id"}
	}

	id := *in.Id

	if in.Modified == nil {
		return nil, &MapError{Kind: "team", ID: id, Field: "Modified"}
	}

	if in.ResourceUrl == nil {
		return nil, &MapError{Kind: "team", ID: id, Field: "ResourceUrl"}
	}

	var imageURL *url.URL
	var err error

	if image := nullableToPtr(in.Image); image != nil {
		imageURL, err = url.Parse(*image)
		if err != nil {
			return nil, &MapError{Kind: "team", ID: id, Field: "Image", Err: err}
		}
	}

	resourceURL, err := url.Parse(*in.ResourceUrl)
	if err != nil {
		return nil, &MapError{Kind: "team", ID: id, Field: "ResourceUrl", Err: err}
	}

	var creators []CreatorList

	if in.Creators != nil {
		creators = make([]CreatorList, 0, len(*in.Creators))

		for _, creator := range *in.Creators {
			c, innerErr := creatorListMapper(creator)
			if innerErr != nil {
				return nil, innerErr
			}

			creators = append(creators, *c)
		}
	}

	var universes []UniverseList

	if in.Universes != nil {
		universes = make([]UniverseList, 0, len(*in.Universes))

		for _, universe := range *in.Universes {
			u, innerErr := universeListMapper(universe)
			if innerErr != nil {
				return nil, innerErr
			}

			universes = append(universes, *u)
		}
	}

	return &Team{
		ID:                    id,
		Name:                  in.Name,
		Description:           in.Desc,
		ImageURL:              imageURL,
		Creators:              creators,
		Universes:             universes,
		ComicVineID:           nullableToPtr(in.CvId),
		GrandComicsDatabaseID: nullableToPtr(in.GcdId),
		ResourceURL:           *resourceURL,
		Modified:              *in.Modified,
	}, nil
}

func teamListMapper(in internal.TeamList) (*TeamList, error) {
	if in.Id == nil {
		return nil, &MapError{Kind: "team", Field: "Id"}
	}

	id := *in.Id

	if in.Modified == nil {
		return nil, &MapError{Kind: "team", ID: id, Field: "Modified"}
	}

	return &TeamList{
		ID:       id,
		Name:     in.Name,
		Modified: *in.Modified,
	}, nil
}
