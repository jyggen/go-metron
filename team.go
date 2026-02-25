package metron

import (
	"context"
	"fmt"
	"iter"
	"net/url"
	"time"

	"github.com/jyggen/go-metron/internal"
)

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

type TeamList struct {
	ID       int
	Name     string
	Modified time.Time
}

// TeamByID returns the information of an individual team.
func (c *Client) TeamByID(ctx context.Context, id int) (*Team, error) {
	return newByID(ctx, c.cache, fmt.Sprintf("team/%d", id), c.client.ApiTeamRetrieve, teamMapper, id)
}

// Teams returns a list of all the teams.
func (c *Client) Teams(ctx context.Context, filters ...Filter) iter.Seq2[*TeamList, error] {
	params := &internal.ApiTeamListParams{}

	for _, f := range filters {
		f(params)
	}

	return newPaginate[internal.PaginatedTeamListList](ctx, c.client.ApiTeamList, teamListMapper, params)
}

func teamMapper(in internal.TeamRead) (*Team, error) {
	var imageURL *url.URL
	var err error

	if in.Image != nil {
		imageURL, err = url.Parse(*in.Image)
		if err != nil {
			return nil, err
		}
	}

	resourceURL, err := url.Parse(*in.ResourceUrl)
	if err != nil {
		return nil, err
	}

	creators := make([]CreatorList, 0, len(*in.Creators))

	for _, creator := range *in.Creators {
		c, innerErr := creatorListMapper(creator)
		if innerErr != nil {
			return nil, innerErr
		}

		creators = append(creators, *c)
	}

	universes := make([]UniverseList, 0, len(*in.Universes))

	for _, universe := range *in.Universes {
		u, innerErr := universeListMapper(universe)
		if innerErr != nil {
			return nil, innerErr
		}

		universes = append(universes, *u)
	}

	return &Team{
		ID:                    *in.Id,
		Name:                  in.Name,
		Description:           in.Desc,
		ImageURL:              imageURL,
		Creators:              creators,
		Universes:             universes,
		ComicVineID:           in.CvId,
		GrandComicsDatabaseID: in.GcdId,
		ResourceURL:           *resourceURL,
		Modified:              *in.Modified,
	}, nil
}

func teamListMapper(in internal.TeamList) (*TeamList, error) {
	return &TeamList{
		ID:       *in.Id,
		Name:     in.Name,
		Modified: *in.Modified,
	}, nil
}
