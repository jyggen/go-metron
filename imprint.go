package metron

import (
	"context"
	"fmt"
	"time"
)

type Imprint struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Founded     *int      `json:"founded"`
	Description *string   `json:"desc"`
	ImageURL    *URL      `json:"image"`
	ComicVineID *int      `json:"cv_id"`
	Publisher   Reference `json:"publisher"`
	ResourceURL URL       `json:"resource_url"`
	Modified    time.Time `json:"modified"`
}

func (i Imprint) modified() time.Time {
	return i.Modified
}

type ImprintList struct {
	ID       int       `json:"id"`
	Name     string    `json:"name"`
	Modified time.Time `json:"modified"`
}

func (c *Client) ImprintByID(ctx context.Context, id int) (Imprint, error) {
	return request[Imprint](ctx, c, fmt.Sprintf("imprint/%d/", id))
}

func (c *Client) Imprints(ctx context.Context, filters ...Filter) func(func(ImprintList, error) bool) {
	return paginate[ImprintList](ctx, c, "imprint/", filters...)
}
