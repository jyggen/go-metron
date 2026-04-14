package metron

import (
	"context"
	"fmt"
	"iter"
	"net/url"
	"time"

	"github.com/jyggen/go-metron/internal"
	"github.com/oapi-codegen/nullable"
)

func nullableToPtr[T any](n nullable.Nullable[T]) *T {
	if v, err := n.Get(); err == nil {
		return &v
	}

	return nil
}

// Arc is a story arc.
type Arc struct {
	ID                    int
	Name                  string
	Description           *string
	ImageURL              *url.URL
	ComicVineID           *int
	GrandComicsDatabaseID *int
	ResourceURL           url.URL
	Modified              time.Time
}

// ArcList is a story arc as it appears in list responses.
type ArcList struct {
	ID       int
	Name     string
	Modified time.Time
}

// ArcByID returns a story arc by its ID.
func (c *Client) ArcByID(ctx context.Context, id int) (*Arc, error) {
	return byID(ctx, c, fmt.Sprintf("arc/%d", id), c.client.ApiArcRetrieve, arcMapper, id)
}

// Arcs returns an iterator over all story arcs.
func (c *Client) Arcs(ctx context.Context, filters ...Filter) iter.Seq2[*ArcList, error] {
	params := &internal.ApiArcListParams{}

	for _, f := range filters {
		f(params)
	}

	return paginate[internal.PaginatedArcListList](ctx, c, "arc", c.client.ApiArcList, arcListMapper, params)
}

func arcMapper(in internal.Arc) (*Arc, error) {
	if in.Id == nil {
		return nil, fmt.Errorf("arc: nil Id")
	}

	if in.Modified == nil {
		return nil, fmt.Errorf("arc: nil Modified")
	}

	if in.ResourceUrl == nil {
		return nil, fmt.Errorf("arc: nil ResourceUrl")
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

	return &Arc{
		ID:                    *in.Id,
		Name:                  in.Name,
		Description:           in.Desc,
		ImageURL:              imageURL,
		ComicVineID:           nullableToPtr(in.CvId),
		GrandComicsDatabaseID: nullableToPtr(in.GcdId),
		ResourceURL:           *resourceURL,
		Modified:              *in.Modified,
	}, nil
}

func arcListMapper(in internal.ArcList) (*ArcList, error) {
	if in.Id == nil {
		return nil, fmt.Errorf("arc: nil Id")
	}

	if in.Modified == nil {
		return nil, fmt.Errorf("arc: nil Modified")
	}

	return &ArcList{
		ID:       *in.Id,
		Name:     in.Name,
		Modified: *in.Modified,
	}, nil
}
