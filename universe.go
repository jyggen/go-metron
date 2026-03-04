package metron

import (
	"context"
	"fmt"
	"iter"
	"net/url"
	"time"

	"github.com/jyggen/go-metron/internal"
)

type Universe struct {
	ID                    int
	Publisher             Reference
	Name                  string
	Designation           string
	Description           *string
	GrandComicsDatabaseID *int
	ImageURL              *url.URL
	ResourceURL           url.URL
	Modified              time.Time
}

type UniverseList struct {
	ID       int
	Name     string
	Modified time.Time
}

// UniverseByID returns a universe by its ID.
func (c *Client) UniverseByID(ctx context.Context, id int) (*Universe, error) {
	return newByID(ctx, c.cache, fmt.Sprintf("universe/%d", id), c.client.ApiUniverseRetrieve, universeMapper, id)
}

// Universes returns an iterator over all universes.
func (c *Client) Universes(ctx context.Context, filters ...Filter) iter.Seq2[*UniverseList, error] {
	params := &internal.ApiUniverseListParams{}

	for _, f := range filters {
		f(params)
	}

	return newPaginate[internal.PaginatedUniverseListList](ctx, c.client.ApiUniverseList, universeListMapper, params)
}

func universeMapper(in internal.UniverseRead) (*Universe, error) {
	if in.Id == nil {
		return nil, fmt.Errorf("universe: nil Id")
	}

	if in.Modified == nil {
		return nil, fmt.Errorf("universe: nil Modified")
	}

	if in.ResourceUrl == nil {
		return nil, fmt.Errorf("universe: nil ResourceUrl")
	}

	if in.Publisher == nil {
		return nil, fmt.Errorf("universe: nil Publisher")
	}

	if in.Publisher.Id == nil {
		return nil, fmt.Errorf("universe: nil Publisher.Id")
	}

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

	return &Universe{
		ID: *in.Id,
		Publisher: Reference{
			ID:   *in.Publisher.Id,
			Name: in.Publisher.Name,
		},
		Name:                  in.Name,
		Designation:           in.Designation,
		Description:           in.Desc,
		GrandComicsDatabaseID: nullableToPtr(in.GcdId),
		ImageURL:              imageURL,
		ResourceURL:           *resourceURL,
		Modified:              *in.Modified,
	}, nil
}

func universeListMapper(in internal.UniverseList) (*UniverseList, error) {
	if in.Id == nil {
		return nil, fmt.Errorf("universe: nil Id")
	}

	if in.Modified == nil {
		return nil, fmt.Errorf("universe: nil Modified")
	}

	return &UniverseList{
		ID:       *in.Id,
		Name:     in.Name,
		Modified: *in.Modified,
	}, nil
}
