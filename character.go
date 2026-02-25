package metron

import (
	"context"
	"fmt"
	"iter"
	"net/url"
	"time"

	"github.com/jyggen/go-metron/internal"
)

type Character struct {
	ID                    int
	Name                  string
	Alias                 *[]string
	Description           *string
	ImageURL              *url.URL
	Creators              []CreatorList
	Teams                 []TeamList
	Universes             []UniverseList
	ComicVineID           *int
	GrandComicsDatabaseID *int
	ResourceURL           url.URL
	Modified              time.Time
}

type CharacterList struct {
	ID       int
	Name     string
	Modified time.Time
}

// CharacterByID returns the information of an individual character.
func (c *Client) CharacterByID(ctx context.Context, id int) (*Character, error) {
	return newByID(ctx, c.cache, fmt.Sprintf("character/%d", id), c.client.ApiCharacterRetrieve, characterMapper, id)
}

// Characters returns a list of all the characters.
func (c *Client) Characters(ctx context.Context, filters ...Filter) iter.Seq2[*CharacterList, error] {
	params := &internal.ApiCharacterListParams{}

	for _, f := range filters {
		f(params)
	}

	return newPaginate[internal.PaginatedCharacterListList](ctx, c.client.ApiCharacterList, characterListMapper, params)
}

func characterMapper(in internal.CharacterRead) (*Character, error) {
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

	teams := make([]TeamList, 0, len(*in.Teams))

	for _, team := range *in.Teams {
		t, innerErr := teamListMapper(team)
		if innerErr != nil {
			return nil, innerErr
		}
		teams = append(teams, *t)
	}

	universes := make([]UniverseList, 0, len(*in.Universes))

	for _, universe := range *in.Universes {
		u, innerErr := universeListMapper(universe)
		if innerErr != nil {
			return nil, innerErr
		}
		universes = append(universes, *u)
	}

	return &Character{
		ID:                    *in.Id,
		Name:                  in.Name,
		Alias:                 in.Alias,
		Description:           in.Desc,
		ImageURL:              imageURL,
		Creators:              creators,
		Teams:                 teams,
		Universes:             universes,
		ComicVineID:           in.CvId,
		GrandComicsDatabaseID: in.GcdId,
		ResourceURL:           *resourceURL,
		Modified:              *in.Modified,
	}, nil
}

func characterListMapper(in internal.CharacterList) (*CharacterList, error) {
	return &CharacterList{
		ID:       *in.Id,
		Name:     in.Name,
		Modified: *in.Modified,
	}, nil
}
