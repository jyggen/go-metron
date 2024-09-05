package metron

import (
	"context"
	"fmt"
	"time"

	"cloud.google.com/go/civil"
)

type Issue struct {
	ID        int        `json:"id"`
	Publisher Reference  `json:"publisher"`
	Imprint   *Reference `json:"imprint"`
	Series    struct {
		ID       int         `json:"id"`
		Name     string      `json:"name"`
		SortName string      `json:"sort_name"`
		Volume   int         `json:"volume"`
		Type     Reference   `json:"series_type"`
		Genres   []Reference `json:"genres"`
	} `json:"series"`
	Number      string      `json:"number"`
	Title       *string     `json:"title"`
	Name        []string    `json:"name"`
	CoverDate   civil.Date  `json:"cover_date"`
	StoreDate   *civil.Date `json:"store_date"`
	Price       string      `json:"price"`
	Rating      Reference   `json:"rating"`
	SKU         *string     `json:"sku"`
	ISBN        *string     `json:"isbn"`
	UPC         *string     `json:"upc"`
	PageCount   *int        `json:"page"`
	Description *string     `json:"desc"`
	ImageURL    *URL        `json:"image"`
	CoverHash   *string     `json:"cover_hash"`
	Arcs        []ArcList   `json:"arcs"`
	Credits     []struct {
		ID    int         `json:"id"`
		Name  string      `json:"creator"`
		Roles []Reference `json:"role"`
	} `json:"credits"`
	Characters []CharacterList `json:"characters"`
	Teams      []TeamList      `json:"teams"`
	Universes  []UniverseList  `json:"universes"`
	Reprints   []struct {
		ID    int    `json:"id"`
		Issue string `json:"issue"`
	} `json:"reprints"`
	Variants []struct {
		Name     *string `json:"name"`
		SKU      *string `json:"sku"`
		UPC      *string `json:"upc"`
		ImageURL URL     `json:"image"`
	} `json:"variants"`
	ComicVineID *int      `json:"cv_id"`
	ResourceURL URL       `json:"resource_url"`
	Modified    time.Time `json:"modified"`
}

func (i Issue) modified() time.Time {
	return i.Modified
}

type IssueList struct {
	ID     int `json:"id"`
	Series struct {
		Name      string `json:"name"`
		Volume    int    `json:"volume"`
		YearBegan int    `json:"year_began"`
	} `json:"series"`
	Number    string      `json:"number"`
	Issue     string      `json:"issue"`
	CoverDate civil.Date  `json:"cover_date"`
	StoreDate *civil.Date `json:"store_date"`
	ImageURL  *URL        `json:"image"`
	CoverHash *string     `json:"cover_hash"`
	Modified  time.Time   `json:"modified"`
}

func (c *Client) IssuesByArcID(ctx context.Context, id int) func(func(IssueList, error) bool) {
	return paginate[IssueList](ctx, c, fmt.Sprintf("arc/%d/issue_list/", id))
}

func (c *Client) IssuesByCharacterID(ctx context.Context, id int) func(func(IssueList, error) bool) {
	return paginate[IssueList](ctx, c, fmt.Sprintf("character/%d/issue_list/", id))
}

func (c *Client) IssuesByPublisherID(ctx context.Context, id int) func(func(IssueList, error) bool) {
	return paginate[IssueList](ctx, c, fmt.Sprintf("publisher/%d/issue_list/", id))
}

func (c *Client) IssuesBySeriesID(ctx context.Context, id int) func(func(IssueList, error) bool) {
	return paginate[IssueList](ctx, c, fmt.Sprintf("series/%d/issue_list/", id))
}

func (c *Client) IssuesByTeamID(ctx context.Context, id int) func(func(IssueList, error) bool) {
	return paginate[IssueList](ctx, c, fmt.Sprintf("team/%d/issue_list/", id))
}

func (c *Client) IssueByID(ctx context.Context, id int) (Issue, error) {
	return request[Issue](ctx, c, fmt.Sprintf("issue/%d/", id))
}

func (c *Client) Issues(ctx context.Context) func(func(IssueList, error) bool) {
	return paginate[IssueList](ctx, c, "issue/")
}
