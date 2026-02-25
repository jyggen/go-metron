package metron

import (
	"context"
	"fmt"
	"iter"
	"net/url"
	"time"

	"github.com/jyggen/go-metron/internal"
)

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

type PublisherList struct {
	ID       int
	Name     string
	Modified time.Time
}

// PublisherByID returns the information of an individual publisher.
func (c *Client) PublisherByID(ctx context.Context, id int) (*Publisher, error) {
	return newByID(ctx, c.cache, fmt.Sprintf("publisher/%d", id), c.client.ApiPublisherRetrieve, publisherMapper, id)
}

// Publishers returns a list of all the publishers.
func (c *Client) Publishers(ctx context.Context, filters ...Filter) iter.Seq2[*PublisherList, error] {
	params := &internal.ApiPublisherListParams{}

	for _, f := range filters {
		f(params)
	}

	return newPaginate[internal.PaginatedPublisherListList](ctx, c.client.ApiPublisherList, publisherListMapper, params)
}

func publisherMapper(in internal.Publisher) (*Publisher, error) {
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
		Founded:               in.Founded,
		CountryCode:           countryCode,
		Description:           in.Desc,
		ImageURL:              imageURL,
		ComicVineID:           in.CvId,
		GrandComicsDatabaseID: in.GcdId,
		ResourceURL:           *resourceURL,
		Modified:              *in.Modified,
	}, nil
}

func publisherListMapper(in internal.PublisherList) (*PublisherList, error) {
	return &PublisherList{
		ID:       *in.Id,
		Name:     in.Name,
		Modified: *in.Modified,
	}, nil
}
