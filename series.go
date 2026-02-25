package metron

import (
	"context"
	"fmt"
	"iter"
	"net/url"
	"time"

	"github.com/jyggen/go-metron/internal"
)

type Series struct {
	ID                    int
	Name                  string
	SortName              string
	Volume                int
	Type                  Reference
	Status                string
	Publisher             Reference
	Imprint               *Reference
	YearBegan             int
	YearEnded             *int
	Description           *string
	IssueCount            int
	Genres                []Reference
	Associated            []Reference
	ComicVineID           *int
	GrandComicsDatabaseID *int
	ResourceURL           url.URL
	Modified              time.Time
}

type SeriesList struct {
	ID         int
	Name       string
	YearBegan  int
	Volume     int
	IssueCount int
	Modified   time.Time
}

// SeriesByID returns the information of an individual series.
func (c *Client) SeriesByID(ctx context.Context, id int) (*Series, error) {
	return newByID(ctx, c.cache, fmt.Sprintf("series/%d", id), c.client.ApiSeriesRetrieve, seriesMapper, id)
}

// Series returns a list of all the series.
func (c *Client) Series(ctx context.Context, filters ...Filter) iter.Seq2[*SeriesList, error] {
	params := &internal.ApiSeriesListParams{}

	for _, f := range filters {
		f(params)
	}

	return newPaginate[internal.PaginatedSeriesListList](ctx, c.client.ApiSeriesList, seriesListMapper, params)
}

// SeriesByPublisherID returns a list of all the series for a publisher.
func (c *Client) SeriesByPublisherID(ctx context.Context, id int, filters ...Filter) iter.Seq2[*SeriesList, error] {
	params := &internal.ApiPublisherSeriesListListParams{}

	for _, f := range filters {
		f(params)
	}

	return newIDPaginate[internal.PaginatedSeriesListList](ctx, c.client.ApiPublisherSeriesListList, seriesListMapper, id, params)
}

func seriesMapper(in internal.SeriesRead) (*Series, error) {
	resourceURL, err := url.Parse(*in.ResourceUrl)
	if err != nil {
		return nil, err
	}

	var imprint *Reference

	if in.Imprint != nil {
		imprint = &Reference{
			ID:   *in.Imprint.Id,
			Name: in.Imprint.Name,
		}
	}

	genres := make([]Reference, 0, len(*in.Genres))

	for _, g := range *in.Genres {
		genres = append(genres, Reference{
			ID:   *g.Id,
			Name: g.Name,
		})
	}

	associated := make([]Reference, 0, len(*in.Associated))

	for _, a := range *in.Associated {
		associated = append(associated, Reference{
			ID:   *a.Id,
			Name: a.Series,
		})
	}

	return &Series{
		ID:       *in.Id,
		Name:     in.Name,
		SortName: in.SortName,
		Volume:   in.Volume,
		Type: Reference{
			ID:   *in.SeriesType.Id,
			Name: in.SeriesType.Name,
		},
		Status: *in.Status,
		Publisher: Reference{
			ID:   *in.Publisher.Id,
			Name: in.Publisher.Name,
		},
		Imprint:               imprint,
		YearBegan:             in.YearBegan,
		YearEnded:             in.YearEnd,
		Description:           in.Desc,
		IssueCount:            *in.IssueCount,
		Genres:                genres,
		Associated:            associated,
		ComicVineID:           in.CvId,
		GrandComicsDatabaseID: in.GcdId,
		ResourceURL:           *resourceURL,
		Modified:              *in.Modified,
	}, nil
}

func seriesListMapper(in internal.SeriesList) (*SeriesList, error) {
	return &SeriesList{
		ID:         *in.Id,
		Name:       in.Series,
		YearBegan:  in.YearBegan,
		Volume:     in.Volume,
		IssueCount: *in.IssueCount,
		Modified:   *in.Modified,
	}, nil
}
