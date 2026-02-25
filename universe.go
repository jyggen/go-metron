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

// UniverseByID returns the information of an individual universe.
func (c *Client) UniverseByID(ctx context.Context, id int) (*Universe, error) {
	return newByID(ctx, c.cache, fmt.Sprintf("universe/%d", id), c.client.ApiUniverseRetrieve, universeMapper, id)
}

// Universes returns a list of all the universes.
func (c *Client) Universes(ctx context.Context, filters ...Filter) iter.Seq2[*UniverseList, error] {
	params := &internal.ApiUniverseListParams{}

	for _, f := range filters {
		f(params)
	}

	return newPaginate[internal.PaginatedUniverseListList](ctx, c.client.ApiUniverseList, universeListMapper, params)
}

func universeMapper(in internal.UniverseRead) (*Universe, error) {
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
		GrandComicsDatabaseID: in.GcdId,
		ImageURL:              imageURL,
		ResourceURL:           *resourceURL,
		Modified:              *in.Modified,
	}, nil
}

func universeListMapper(in internal.UniverseList) (*UniverseList, error) {
	return &UniverseList{
		ID:       *in.Id,
		Name:     in.Name,
		Modified: *in.Modified,
	}, nil
}
