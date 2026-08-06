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
	return byID(ctx, c, fmt.Sprintf("creator/%d", id), c.client.ApiCreatorRetrieve, creatorMapper, id)
}

// Creators returns an iterator over all creators.
func (c *Client) Creators(ctx context.Context, filters ...Filter) iter.Seq2[*CreatorList, error] {
	params := &internal.ApiCreatorListParams{}

	if err := applyFilters("Creators", params, filters); err != nil {
		return errIter[*CreatorList](err)
	}

	return paginate[internal.PaginatedCreatorListList](ctx, c, "creator", c.client.ApiCreatorList, creatorListMapper, params)
}

func creatorMapper(in internal.Creator) (*Creator, error) {
	if in.Id == nil {
		return nil, &MapError{Kind: "creator", Field: "Id"}
	}

	id := *in.Id

	if in.Modified == nil {
		return nil, &MapError{Kind: "creator", ID: id, Field: "Modified"}
	}

	if in.ResourceUrl == nil {
		return nil, &MapError{Kind: "creator", ID: id, Field: "ResourceUrl"}
	}

	var imageURL *url.URL
	var err error

	if image := nullableToPtr(in.Image); image != nil {
		imageURL, err = url.Parse(*image)
		if err != nil {
			return nil, &MapError{Kind: "creator", ID: id, Field: "Image", Err: err}
		}
	}

	resourceURL, err := url.Parse(*in.ResourceUrl)
	if err != nil {
		return nil, &MapError{Kind: "creator", ID: id, Field: "ResourceUrl", Err: err}
	}

	var birth *civil.Date
	var death *civil.Date

	if d, err := in.Birth.Get(); err == nil {
		birth = new(civil.DateOf(d.Time))
	}

	if d, err := in.Death.Get(); err == nil {
		death = new(civil.DateOf(d.Time))
	}

	return &Creator{
		ID:                    id,
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
		return nil, &MapError{Kind: "creator", Field: "Id"}
	}

	id := *in.Id

	if in.Modified == nil {
		return nil, &MapError{Kind: "creator", ID: id, Field: "Modified"}
	}

	return &CreatorList{
		ID:       id,
		Name:     in.Name,
		Modified: *in.Modified,
	}, nil
}
