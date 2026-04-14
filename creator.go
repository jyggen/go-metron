package metron

import (
	"context"
	"fmt"
	"iter"
	"net/url"
	"time"

	"cloud.google.com/go/civil"
	"github.com/jyggen/go-metron/internal"
)

// Creator is a comic book creator.
type Creator struct {
	ID                    int
	Name                  string
	Birth                 *civil.Date
	Death                 *civil.Date
	Description           *string
	ImageURL              *url.URL
	Alias                 *[]string
	ComicVineID           *int
	GrandComicsDatabaseID *int
	ResourceURL           url.URL
	Modified              time.Time
}

// CreatorList is a creator as it appears in list responses.
type CreatorList struct {
	ID       int
	Name     string
	Modified time.Time
}

// CreatorByID returns a creator by its ID.
func (c *Client) CreatorByID(ctx context.Context, id int) (*Creator, error) {
	return newByID(ctx, c.cache, c.maxRetries, fmt.Sprintf("creator/%d", id), c.client.ApiCreatorRetrieve, creatorMapper, id)
}

// Creators returns an iterator over all creators.
func (c *Client) Creators(ctx context.Context, filters ...Filter) iter.Seq2[*CreatorList, error] {
	params := &internal.ApiCreatorListParams{}

	for _, f := range filters {
		f(params)
	}

	return newPaginate[internal.PaginatedCreatorListList](ctx, c.cache, c.maxRetries, "creator", c.client.ApiCreatorList, creatorListMapper, params)
}

func creatorMapper(in internal.Creator) (*Creator, error) {
	if in.Id == nil {
		return nil, fmt.Errorf("creator: nil Id")
	}

	if in.Modified == nil {
		return nil, fmt.Errorf("creator: nil Modified")
	}

	if in.ResourceUrl == nil {
		return nil, fmt.Errorf("creator: nil ResourceUrl")
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

	var birth *civil.Date
	var death *civil.Date

	if d, err := in.Birth.Get(); err == nil {
		noRefBirth := civil.DateOf(d.Time)
		birth = &noRefBirth
	}

	if d, err := in.Death.Get(); err == nil {
		noRefDeath := civil.DateOf(d.Time)
		death = &noRefDeath
	}

	return &Creator{
		ID:                    *in.Id,
		Name:                  in.Name,
		Birth:                 birth,
		Death:                 death,
		Description:           in.Desc,
		ImageURL:              imageURL,
		Alias:                 in.Alias,
		ComicVineID:           nullableToPtr(in.CvId),
		GrandComicsDatabaseID: nullableToPtr(in.GcdId),
		ResourceURL:           *resourceURL,
		Modified:              *in.Modified,
	}, nil
}

func creatorListMapper(in internal.CreatorList) (*CreatorList, error) {
	if in.Id == nil {
		return nil, fmt.Errorf("creator: nil Id")
	}

	if in.Modified == nil {
		return nil, fmt.Errorf("creator: nil Modified")
	}

	return &CreatorList{
		ID:       *in.Id,
		Name:     in.Name,
		Modified: *in.Modified,
	}, nil
}
