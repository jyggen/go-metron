package metron

import (
	"context"
	"fmt"
	"time"
)

type Team struct {
	ID          int            `json:"id"`
	Name        string         `json:"name"`
	Description *string        `json:"desc"`
	ImageURL    *URL           `json:"image"`
	Creators    []CreatorList  `json:"creators"`
	Universes   []UniverseList `json:"universes"`
	ComicVineID *int           `json:"cv_id"`
	ResourceURL URL            `json:"resource_url"`
	Modified    time.Time      `json:"modified"`
}

func (t Team) modified() time.Time {
	return t.Modified
}

type TeamList struct {
	ID       int       `json:"id"`
	Name     string    `json:"name"`
	Modified time.Time `json:"modified"`
}

func (c *Client) TeamByID(ctx context.Context, id int) (Team, error) {
	return request[Team](ctx, c, fmt.Sprintf("team/%d/", id))
}

func (c *Client) Teams(ctx context.Context, filters ...Filter) func(func(TeamList, error) bool) {
	return paginate[TeamList](ctx, c, "team/", filters...)
}
