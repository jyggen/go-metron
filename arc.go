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

	if err := applyFilters("Arcs", params, filters); err != nil {
		return errIter[*ArcList](err)
	}

	return paginate[internal.PaginatedArcListList](ctx, c, "arc", c.client.ApiArcList, arcListMapper, params)
}

func arcMapper(in internal.Arc) (*Arc, error) {
	if in.Id == nil {
		return nil, &MapError{Kind: "arc", Field: "Id"}
	}

	id := *in.Id

	if in.Modified == nil {
		return nil, &MapError{Kind: "arc", ID: id, Field: "Modified"}
	}

	if in.ResourceUrl == nil {
		return nil, &MapError{Kind: "arc", ID: id, Field: "ResourceUrl"}
	}

	var imageURL *url.URL
	var err error

	if image := nullableToPtr(in.Image); image != nil {
		imageURL, err = url.Parse(*image)
		if err != nil {
			return nil, &MapError{Kind: "arc", ID: id, Field: "Image", Err: err}
		}
	}

	resourceURL, err := url.Parse(*in.ResourceUrl)
	if err != nil {
		return nil, &MapError{Kind: "arc", ID: id, Field: "ResourceUrl", Err: err}
	}

	return &Arc{
		ID:                    id,
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
		return nil, &MapError{Kind: "arc", Field: "Id"}
	}

	id := *in.Id

	if in.Modified == nil {
		return nil, &MapError{Kind: "arc", ID: id, Field: "Modified"}
	}

	return &ArcList{
		ID:       id,
		Name:     in.Name,
		Modified: *in.Modified,
	}, nil
}
