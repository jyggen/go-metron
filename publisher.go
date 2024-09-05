package metron

import (
	"context"
	"fmt"
	"time"
)

type Publisher struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Founded     *int      `json:"founded"`
	Description *string   `json:"desc"`
	ImageURL    *URL      `json:"image"`
	ComicVineID *int      `json:"cv_id"`
	ResourceURL URL       `json:"resource_url"`
	Modified    time.Time `json:"modified"`
}

func (p Publisher) modified() time.Time {
	return p.Modified
}

type PublisherList struct {
	ID       int       `json:"id"`
	Name     string    `json:"name"`
	Modified time.Time `json:"modified"`
}

func (c *Client) PublisherByID(ctx context.Context, id int) (Publisher, error) {
	return request[Publisher](ctx, c, fmt.Sprintf("publisher/%d/", id))
}

func (c *Client) Publishers(ctx context.Context, filters ...Filter) func(func(PublisherList, error) bool) {
	return paginate[PublisherList](ctx, c, "publisher/", filters...)
}
