package metron

import (
	"context"
	"fmt"
	"time"
)

type Character struct {
	ID          int            `json:"id"`
	Name        string         `json:"name"`
	Alias       *[]string      `json:"alias"`
	Description *string        `json:"desc"`
	ImageURL    *URL           `json:"image"`
	Creators    []CreatorList  `json:"creators"`
	Teams       []TeamList     `json:"teams"`
	Universes   []UniverseList `json:"universes"`
	ComicVineID *int           `json:"cv_id"`
	ResourceURL URL            `json:"resource_url"`
	Modified    time.Time      `json:"modified"`
}

func (c Character) modified() time.Time {
	return c.Modified
}

type CharacterList struct {
	ID       int       `json:"id"`
	Name     string    `json:"name"`
	Modified time.Time `json:"modified"`
}

func (c *Client) CharacterByID(ctx context.Context, id int) (Character, error) {
	return request[Character](ctx, c, fmt.Sprintf("character/%d/", id))
}

func (c *Client) Characters(ctx context.Context) func(func(CharacterList, error) bool) {
	return paginate[CharacterList](ctx, c, "character/")
}
