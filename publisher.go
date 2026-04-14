package metron

import (
	"context"
	"fmt"
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
func (c *Client) PublisherByID(ctx context.Context, id int) (*Publisher, error) {
	return byID(ctx, c.cache, c.maxRetries, fmt.Sprintf("publisher/%d", id), c.client.ApiPublisherRetrieve, publisherMapper, id)
}

// Publishers returns an iterator over all publishers.
func (c *Client) Publishers(ctx context.Context, filters ...Filter) iter.Seq2[*PublisherList, error] {
	params := &internal.ApiPublisherListParams{}

	for _, f := range filters {
		f(params)
	}

	return paginate[internal.PaginatedPublisherListList](ctx, c.cache, c.maxRetries, "publisher", c.client.ApiPublisherList, publisherListMapper, params)
}

func publisherMapper(in internal.Publisher) (*Publisher, error) {
	if in.Id == nil {
		return nil, fmt.Errorf("publisher: nil Id")
	}

	if in.Modified == nil {
		return nil, fmt.Errorf("publisher: nil Modified")
	}

	if in.ResourceUrl == nil {
		return nil, fmt.Errorf("publisher: nil ResourceUrl")
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

	var countryCode *string

	if in.Country != nil {
		countryCode = func() *string {
			s := string(*in.Country)

			return &s
		}()
	}

	return &Publisher{
		ID:                    *in.Id,
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
		return nil, fmt.Errorf("publisher: nil Id")
	}

	if in.Modified == nil {
		return nil, fmt.Errorf("publisher: nil Modified")
	}

	return &PublisherList{
		ID:       *in.Id,
		Name:     in.Name,
		Modified: *in.Modified,
	}, nil
}
