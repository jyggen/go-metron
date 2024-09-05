package metron

import (
	"context"
	"fmt"
	"time"
)

type Series struct {
	ID          int         `json:"id"`
	Name        string      `json:"name"`
	SortName    string      `json:"sort_name"`
	Volume      int         `json:"volume"`
	Type        Reference   `json:"series_type"`
	Status      string      `json:"status"`
	Publisher   Reference   `json:"publisher"`
	Imprint     *Reference  `json:"imprint"`
	YearBegan   int         `json:"year_began"`
	YearEnded   *int        `json:"year_end"`
	Description *string     `json:"desc"`
	IssueCount  int         `json:"issue_count"`
	Genres      []Reference `json:"genres"`
	Associated  []struct {
		ID   int    `json:"id"`
		Name string `json:"series"`
	} `json:"associated"`
	ComicVineID *int      `json:"cv_id"`
	ResourceURL URL       `json:"resource_url"`
	Modified    time.Time `json:"modified"`
}

func (s Series) modified() time.Time {
	return s.Modified
}

type SeriesList struct {
	ID         int       `json:"id"`
	Name       string    `json:"series"`
	YearBegan  int       `json:"year_began"`
	Volume     int       `json:"volume"`
	IssueCount int       `json:"issue_count"`
	Modified   time.Time `json:"modified"`
}

func (c *Client) SeriesByID(ctx context.Context, id int) (Series, error) {
	return request[Series](ctx, c, fmt.Sprintf("series/%d/", id))
}

func (c *Client) Series(ctx context.Context) func(func(SeriesList, error) bool) {
	return paginate[SeriesList](ctx, c, "series/")
}
