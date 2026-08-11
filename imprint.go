package metron

import (
	"context"
	"iter"
	"net/url"
	"time"

	"github.com/jyggen/go-metron/internal"
)

// Imprint is a publisher imprint.
type Imprint struct {
	ID                    int
	Name                  string
	Founded               *int
	Description           *string
	ImageURL              *url.URL
	ComicVineID           *int
	GrandComicsDatabaseID *int
	Publisher             Reference
	ResourceURL           url.URL
	Modified              time.Time
}

// ImprintList is an imprint as it appears in list responses.
type ImprintList struct {
	ID       int
	Name     string
	Modified time.Time
}

// ImprintByID returns an imprint by its ID.
func (c *Client) ImprintByID(ctx context.Context, id int, opts ...ConditionalOption) (*Imprint, error) {
	return byID(ctx, c, c.client.ApiImprintRetrieve, imprintMapper, id, conditionalEditors(opts))
}

// Imprints returns an iterator over all imprints.
func (c *Client) Imprints(ctx context.Context, filters *ImprintFilters, opts ...RequestOption) iter.Seq2[*ImprintList, error] {
	return paginate[internal.PaginatedImprintListList](ctx, c, c.client.ApiImprintList, imprintListMapper, filters.params, everyPage(requestEditors(opts)))
}

func imprintMapper(in internal.ImprintRead) (*Imprint, error) {
	if in.Id == nil {
		return nil, &MapError{Kind: "imprint", Field: "Id"}
	}

	id := *in.Id

	if in.Modified == nil {
		return nil, &MapError{Kind: "imprint", ID: id, Field: "Modified"}
	}

	if in.ResourceUrl == nil {
		return nil, &MapError{Kind: "imprint", ID: id, Field: "ResourceUrl"}
	}

	if in.Publisher == nil {
		return nil, &MapError{Kind: "imprint", ID: id, Field: "Publisher"}
	}

	if in.Publisher.Id == nil {
		return nil, &MapError{Kind: "imprint", ID: id, Field: "Publisher.Id"}
	}

	var imageURL *url.URL
	var err error

	if imgStr, imgErr := in.Image.Get(); imgErr == nil {
		imageURL, err = url.Parse(imgStr)
		if err != nil {
			return nil, &MapError{Kind: "imprint", ID: id, Field: "Image", Err: err}
		}
	}

	resourceURL, err := url.Parse(*in.ResourceUrl)
	if err != nil {
		return nil, &MapError{Kind: "imprint", ID: id, Field: "ResourceUrl", Err: err}
	}

	return &Imprint{
		ID:                    id,
		Name:                  in.Name,
		Founded:               nullableToPtr(in.Founded),
		Description:           in.Desc,
		ImageURL:              imageURL,
		ComicVineID:           nullableToPtr(in.CvId),
		GrandComicsDatabaseID: nullableToPtr(in.GcdId),
		Publisher: Reference{
			ID:   *in.Publisher.Id,
			Name: in.Publisher.Name,
		},
		ResourceURL: *resourceURL,
		Modified:    *in.Modified,
	}, nil
}

func imprintListMapper(in internal.ImprintList) (*ImprintList, error) {
	if in.Id == nil {
		return nil, &MapError{Kind: "imprint", Field: "Id"}
	}

	id := *in.Id

	if in.Modified == nil {
		return nil, &MapError{Kind: "imprint", ID: id, Field: "Modified"}
	}

	return &ImprintList{
		ID:       id,
		Name:     in.Name,
		Modified: *in.Modified,
	}, nil
}
