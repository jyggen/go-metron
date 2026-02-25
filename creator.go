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

type CreatorList struct {
	ID       int
	Name     string
	Modified time.Time
}

// CreatorByID returns the information of an individual creator.
func (c *Client) CreatorByID(ctx context.Context, id int) (*Creator, error) {
	return newByID(ctx, c.cache, fmt.Sprintf("creator/%d", id), c.client.ApiCreatorRetrieve, creatorMapper, id)
}

// Creators returns a list of all the creators.
func (c *Client) Creators(ctx context.Context, filters ...Filter) iter.Seq2[*CreatorList, error] {
	params := &internal.ApiCreatorListParams{}

	for _, f := range filters {
		f(params)
	}

	return newPaginate[internal.PaginatedCreatorListList](ctx, c.client.ApiCreatorList, creatorListMapper, params)
}

func creatorMapper(in internal.Creator) (*Creator, error) {
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

	if in.Birth != nil {
		noRefBirth := civil.DateOf(in.Birth.Time)
		birth = &noRefBirth
	}

	if in.Death != nil {
		noRefDeath := civil.DateOf(in.Death.Time)
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
		ComicVineID:           in.CvId,
		GrandComicsDatabaseID: in.GcdId,
		ResourceURL:           *resourceURL,
		Modified:              *in.Modified,
	}, nil
}

func creatorListMapper(in internal.CreatorList) (*CreatorList, error) {
	return &CreatorList{
		ID:       *in.Id,
		Name:     in.Name,
		Modified: *in.Modified,
	}, nil
}
