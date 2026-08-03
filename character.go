package metron

import (
	"context"
	"fmt"
	"iter"
	"net/url"
	"time"

	"github.com/jyggen/go-metron/internal"
)

// Character is a comic book character.
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

// CharacterList is a character as it appears in list responses.
type CharacterList struct {
	ID       int
	Name     string
	Modified time.Time
}

// CharacterByID returns a character by its ID.
func (c *Client) CharacterByID(ctx context.Context, id int) (*Character, error) {
	return byID(ctx, c, fmt.Sprintf("character/%d", id), c.client.ApiCharacterRetrieve, characterMapper, id)
}

// Characters returns an iterator over all characters.
func (c *Client) Characters(ctx context.Context, filters ...Filter) iter.Seq2[*CharacterList, error] {
	params := &internal.ApiCharacterListParams{}

	if err := applyFilters(params, filters); err != nil {
		return errIter[*CharacterList](err)
	}

	return paginate[internal.PaginatedCharacterListList](ctx, c, "character", c.client.ApiCharacterList, characterListMapper, params)
}

func characterMapper(in internal.CharacterRead) (*Character, error) {
	if in.Id == nil {
		return nil, fmt.Errorf("character: nil Id")
	}

	if in.Modified == nil {
		return nil, fmt.Errorf("character: nil Modified")
	}

	if in.ResourceUrl == nil {
		return nil, fmt.Errorf("character: nil ResourceUrl")
	}

	var imageURL *url.URL
	var err error

	if image := nullableToPtr(in.Image); image != nil {
		imageURL, err = url.Parse(*image)
		if err != nil {
			return nil, err
		}
	}

	resourceURL, err := url.Parse(*in.ResourceUrl)
	if err != nil {
		return nil, err
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

	var teams []TeamList

	if in.Teams != nil {
		teams = make([]TeamList, 0, len(*in.Teams))

		for _, team := range *in.Teams {
			t, innerErr := teamListMapper(team)
			if innerErr != nil {
				return nil, innerErr
			}
			teams = append(teams, *t)
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

	return &Character{
		ID:                    *in.Id,
		Name:                  in.Name,
		Alias:                 in.Alias,
		Description:           in.Desc,
		ImageURL:              imageURL,
		Creators:              creators,
		Teams:                 teams,
		Universes:             universes,
		ComicVineID:           nullableToPtr(in.CvId),
		GrandComicsDatabaseID: nullableToPtr(in.GcdId),
		ResourceURL:           *resourceURL,
		Modified:              *in.Modified,
	}, nil
}

func characterListMapper(in internal.CharacterList) (*CharacterList, error) {
	if in.Id == nil {
		return nil, fmt.Errorf("character: nil Id")
	}

	if in.Modified == nil {
		return nil, fmt.Errorf("character: nil Modified")
	}

	return &CharacterList{
		ID:       *in.Id,
		Name:     in.Name,
		Modified: *in.Modified,
	}, nil
}
