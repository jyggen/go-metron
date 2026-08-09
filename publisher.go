package metron

import (
	"context"
	"iter"
	"net/url"
	"time"

	"github.com/jyggen/go-metron/internal"
)

// Publisher is a comic book publisher.
type Publisher struct {
	ID                    int
	Name                  string
	Founded               *int
	CountryCode           *string
	Description           *string
	ImageURL              *url.URL
	ComicVineID           *int
	GrandComicsDatabaseID *int
	ResourceURL           url.URL
	Modified              time.Time
}

// PublisherList is a publisher as it appears in list responses.
type PublisherList struct {
	ID       int
	Name     string
	Modified time.Time
}

// PublisherByID returns a publisher by its ID.
func (c *Client) PublisherByID(ctx context.Context, id int, opts ...RequestOption) (*Publisher, error) {
	return byID(ctx, c, c.client.ApiPublisherRetrieve, publisherMapper, id, opts)
}

// Publishers returns an iterator over all publishers.
func (c *Client) Publishers(ctx context.Context, filters ...Filter) iter.Seq2[*PublisherList, error] {
	params := &internal.ApiPublisherListParams{}

	if err := applyFilters("Publishers", params, filters); err != nil {
		return errIter[*PublisherList](err)
	}

	return paginate[internal.PaginatedPublisherListList](ctx, c, c.client.ApiPublisherList, publisherListMapper, params)
}

func publisherMapper(in internal.Publisher) (*Publisher, error) {
	if in.Id == nil {
		return nil, &MapError{Kind: "publisher", Field: "Id"}
	}

	id := *in.Id

	if in.Modified == nil {
		return nil, &MapError{Kind: "publisher", ID: id, Field: "Modified"}
	}

	if in.ResourceUrl == nil {
		return nil, &MapError{Kind: "publisher", ID: id, Field: "ResourceUrl"}
	}

	var imageURL *url.URL
	var err error

	if image := nullableToPtr(in.Image); image != nil {
		imageURL, err = url.Parse(*image)
		if err != nil {
			return nil, &MapError{Kind: "publisher", ID: id, Field: "Image", Err: err}
		}
	}

	resourceURL, err := url.Parse(*in.ResourceUrl)
	if err != nil {
		return nil, &MapError{Kind: "publisher", ID: id, Field: "ResourceUrl", Err: err}
	}

	var countryCode *string

	if in.Country != nil {
		countryCode = func() *string {
			s := string(*in.Country)

			return &s
		}()
	}

	return &Publisher{
		ID:                    id,
		Name:                  in.Name,
		Founded:               nullableToPtr(in.Founded),
		CountryCode:           countryCode,
		Description:           in.Desc,
		ImageURL:              imageURL,
		ComicVineID:           nullableToPtr(in.CvId),
		GrandComicsDatabaseID: nullableToPtr(in.GcdId),
		ResourceURL:           *resourceURL,
		Modified:              *in.Modified,
	}, nil
}

func publisherListMapper(in internal.PublisherList) (*PublisherList, error) {
	if in.Id == nil {
		return nil, &MapError{Kind: "publisher", Field: "Id"}
	}

	id := *in.Id

	if in.Modified == nil {
		return nil, &MapError{Kind: "publisher", ID: id, Field: "Modified"}
	}

	return &PublisherList{
		ID:       id,
		Name:     in.Name,
		Modified: *in.Modified,
	}, nil
}
