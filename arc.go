package metron

import (
	"context"
	"fmt"
	"time"
)

type Arc struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description *string   `json:"desc"`
	ImageURL    *URL      `json:"image"`
	ComicVineID *int      `json:"cv_id"`
	ResourceURL URL       `json:"resource_url"`
	Modified    time.Time `json:"modified"`
}

func (a Arc) modified() time.Time {
	return a.Modified
}

type ArcList struct {
	ID       int       `json:"id"`
	Name     string    `json:"name"`
	Modified time.Time `json:"modified"`
}

func (c *Client) ArcByID(ctx context.Context, id int) (Arc, error) {
	return request[Arc](ctx, c, fmt.Sprintf("arc/%d/", id))
}

func (c *Client) Arcs(ctx context.Context) func(func(ArcList, error) bool) {
	return paginate[ArcList](ctx, c, "arc/")
}
