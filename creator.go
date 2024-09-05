package metron

import (
	"context"
	"fmt"
	"time"

	"cloud.google.com/go/civil"
)

type Creator struct {
	ID          int         `json:"id"`
	Name        string      `json:"name"`
	Birth       *civil.Date `json:"birth"`
	Death       *civil.Date `json:"death"`
	Description *string     `json:"desc"`
	ImageURL    *URL        `json:"image"`
	Alias       *[]string   `json:"alias"`
	ComicVineID *int        `json:"cv_id"`
	ResourceURL URL         `json:"resource_url"`
	Modified    time.Time   `json:"modified"`
}

func (c Creator) modified() time.Time {
	return c.Modified
}

type CreatorList struct {
	ID       int       `json:"id"`
	Name     string    `json:"name"`
	Modified time.Time `json:"modified"`
}

func (c *Client) CreatorByID(ctx context.Context, id int) (Creator, error) {
	return request[Creator](ctx, c, fmt.Sprintf("creator/%d/", id))
}

func (c *Client) Creators(ctx context.Context, filters ...Filter) func(func(CreatorList, error) bool) {
	return paginate[CreatorList](ctx, c, "creator/", filters...)
}
