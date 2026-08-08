package metron

import (
	"context"
	"iter"
	"net/url"
	"time"

	"github.com/jyggen/go-metron/internal"
)

// Universe is a fictional universe.
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

// UniverseList is a universe as it appears in list responses.
type UniverseList struct {
	ID       int
	Name     string
	Modified time.Time
}

// UniverseByID returns a universe by its ID.
func (c *Client) UniverseByID(ctx context.Context, id int) (*Universe, error) {
	return byID(ctx, c, c.client.ApiUniverseRetrieve, universeMapper, id)
}

// Universes returns an iterator over all universes.
func (c *Client) Universes(ctx context.Context, filters ...Filter) iter.Seq2[*UniverseList, error] {
	params := &internal.ApiUniverseListParams{}

	if err := applyFilters("Universes", params, filters); err != nil {
		return errIter[*UniverseList](err)
	}

	return paginate[internal.PaginatedUniverseListList](ctx, c, c.client.ApiUniverseList, universeListMapper, params)
}

func universeMapper(in internal.UniverseRead) (*Universe, error) {
	if in.Id == nil {
		return nil, &MapError{Kind: "universe", Field: "Id"}
	}

	id := *in.Id

	if in.Modified == nil {
		return nil, &MapError{Kind: "universe", ID: id, Field: "Modified"}
	}

	if in.ResourceUrl == nil {
		return nil, &MapError{Kind: "universe", ID: id, Field: "ResourceUrl"}
	}

	if in.Publisher == nil {
		return nil, &MapError{Kind: "universe", ID: id, Field: "Publisher"}
	}

	if in.Publisher.Id == nil {
		return nil, &MapError{Kind: "universe", ID: id, Field: "Publisher.Id"}
	}

	var imageURL *url.URL
	var err error

	if image := nullableToPtr(in.Image); image != nil {
		imageURL, err = url.Parse(*image)
		if err != nil {
			return nil, &MapError{Kind: "universe", ID: id, Field: "Image", Err: err}
		}
	}

	resourceURL, err := url.Parse(*in.ResourceUrl)
	if err != nil {
		return nil, &MapError{Kind: "universe", ID: id, Field: "ResourceUrl", Err: err}
	}

	return &Universe{
		ID: id,
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
		return nil, &MapError{Kind: "universe", Field: "Id"}
	}

	id := *in.Id

	if in.Modified == nil {
		return nil, &MapError{Kind: "universe", ID: id, Field: "Modified"}
	}

	return &UniverseList{
		ID:       id,
		Name:     in.Name,
		Modified: *in.Modified,
	}, nil
}
