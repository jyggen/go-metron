package metron

import (
	"context"
	"fmt"
	"iter"
	"net/url"
	"time"

	"github.com/jyggen/go-metron/internal"
)

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

type ImprintList struct {
	ID       int
	Name     string
	Modified time.Time
}

// ImprintByID returns the information of an individual imprint.
func (c *Client) ImprintByID(ctx context.Context, id int) (*Imprint, error) {
	return newByID(ctx, c.cache, fmt.Sprintf("imprint/%d", id), c.client.ApiImprintRetrieve, imprintMapper, id)
}

// Imprints returns a list of all the imprints.
func (c *Client) Imprints(ctx context.Context, filters ...Filter) iter.Seq2[*ImprintList, error] {
	params := &internal.ApiImprintListParams{}

	for _, f := range filters {
		f(params)
	}

	return newPaginate[internal.PaginatedImprintListList](ctx, c.client.ApiImprintList, imprintListMapper, params)
}

func imprintMapper(in internal.ImprintRead) (*Imprint, error) {
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

	return &Imprint{
		ID:                    *in.Id,
		Name:                  in.Name,
		Founded:               in.Founded,
		Description:           in.Desc,
		ImageURL:              imageURL,
		ComicVineID:           in.CvId,
		GrandComicsDatabaseID: in.GcdId,
		Publisher: Reference{
			ID:   *in.Publisher.Id,
			Name: in.Publisher.Name,
		},
		ResourceURL: *resourceURL,
		Modified:    *in.Modified,
	}, nil
}

func imprintListMapper(in internal.ImprintList) (*ImprintList, error) {
	return &ImprintList{
		ID:       *in.Id,
		Name:     in.Name,
		Modified: *in.Modified,
	}, nil
}
