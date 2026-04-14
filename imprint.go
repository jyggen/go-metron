package metron

import (
	"context"
	"fmt"
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
func (c *Client) ImprintByID(ctx context.Context, id int) (*Imprint, error) {
	return byID(ctx, c.cache, c.maxRetries, fmt.Sprintf("imprint/%d", id), c.client.ApiImprintRetrieve, imprintMapper, id)
}

// Imprints returns an iterator over all imprints.
func (c *Client) Imprints(ctx context.Context, filters ...Filter) iter.Seq2[*ImprintList, error] {
	params := &internal.ApiImprintListParams{}

	for _, f := range filters {
		f(params)
	}

	return paginate[internal.PaginatedImprintListList](ctx, c.cache, c.maxRetries, "imprint", c.client.ApiImprintList, imprintListMapper, params)
}

func imprintMapper(in internal.ImprintRead) (*Imprint, error) {
	if in.Id == nil {
		return nil, fmt.Errorf("imprint: nil Id")
	}

	if in.Modified == nil {
		return nil, fmt.Errorf("imprint: nil Modified")
	}

	if in.ResourceUrl == nil {
		return nil, fmt.Errorf("imprint: nil ResourceUrl")
	}

	if in.Publisher == nil {
		return nil, fmt.Errorf("imprint: nil Publisher")
	}

	if in.Publisher.Id == nil {
		return nil, fmt.Errorf("imprint: nil Publisher.Id")
	}

	var imageURL *url.URL
	var err error

	if imgStr, imgErr := in.Image.Get(); imgErr == nil {
		imageURL, err = url.Parse(imgStr)
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
		return nil, fmt.Errorf("imprint: nil Id")
	}

	if in.Modified == nil {
		return nil, fmt.Errorf("imprint: nil Modified")
	}

	return &ImprintList{
		ID:       *in.Id,
		Name:     in.Name,
		Modified: *in.Modified,
	}, nil
}
