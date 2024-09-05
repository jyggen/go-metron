package metron

import (
	"context"
	"fmt"
	"time"
)

type Universe struct {
	ID          int       `json:"id"`
	Publisher   Reference `json:"publisher"`
	Name        string    `json:"name"`
	Designation *string   `json:"designation"`
	Description *string   `json:"desc"`
	ImageURL    *URL      `json:"image"`
	ResourceURL URL       `json:"resource_url"`
	Modified    time.Time `json:"modified"`
}

func (u Universe) modified() time.Time {
	return u.Modified
}

type UniverseList struct {
	ID       int       `json:"id"`
	Name     string    `json:"name"`
	Modified time.Time `json:"modified"`
}

func (c *Client) UniverseByID(ctx context.Context, id int) (Universe, error) {
	return request[Universe](ctx, c, fmt.Sprintf("universe/%d/", id))
}

func (c *Client) Universes(ctx context.Context) func(func(UniverseList, error) bool) {
	return paginate[UniverseList](ctx, c, "universe/")
}
