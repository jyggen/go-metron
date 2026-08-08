package metron

import (
	"context"
	"iter"
	"net/url"
	"time"

	"github.com/jyggen/go-metron/internal"
)

// Series is a comic book series.
type Series struct {
	ID                    int
	Name                  string
	AlternativeNames      []string
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

// SeriesList is a series as it appears in list responses.
type SeriesList struct {
	ID         int
	Name       string
	YearBegan  int
	YearEnded  *int
	Volume     int
	IssueCount int
	Modified   time.Time
}

// SeriesByID returns a series by its ID.
func (c *Client) SeriesByID(ctx context.Context, id int) (*Series, error) {
	return byID(ctx, c, c.client.ApiSeriesRetrieve, seriesMapper, id)
}

// Series returns an iterator over all series.
func (c *Client) Series(ctx context.Context, filters ...Filter) iter.Seq2[*SeriesList, error] {
	params := &internal.ApiSeriesListParams{}

	if err := applyFilters("Series", params, filters); err != nil {
		return errIter[*SeriesList](err)
	}

	return paginate[internal.PaginatedSeriesListList](ctx, c, c.client.ApiSeriesList, seriesListMapper, params)
}

// SeriesByPublisherID returns an iterator over all series for a publisher.
func (c *Client) SeriesByPublisherID(ctx context.Context, id int, filters ...Filter) iter.Seq2[*SeriesList, error] {
	params := &internal.ApiPublisherSeriesListListParams{}

	if err := applyFilters("SeriesByPublisherID", params, filters); err != nil {
		return errIter[*SeriesList](err)
	}

	return idPaginate[internal.PaginatedSeriesListList](ctx, c, c.client.ApiPublisherSeriesListList, seriesListMapper, id, params)
}

func seriesMapper(in internal.SeriesRead) (*Series, error) {
	if in.Id == nil {
		return nil, &MapError{Kind: "series", Field: "Id"}
	}

	id := *in.Id

	if in.Modified == nil {
		return nil, &MapError{Kind: "series", ID: id, Field: "Modified"}
	}

	if in.ResourceUrl == nil {
		return nil, &MapError{Kind: "series", ID: id, Field: "ResourceUrl"}
	}

	if in.Status == nil {
		return nil, &MapError{Kind: "series", ID: id, Field: "Status"}
	}

	if in.Publisher == nil {
		return nil, &MapError{Kind: "series", ID: id, Field: "Publisher"}
	}

	if in.Publisher.Id == nil {
		return nil, &MapError{Kind: "series", ID: id, Field: "Publisher.Id"}
	}

	if in.SeriesType == nil {
		return nil, &MapError{Kind: "series", ID: id, Field: "SeriesType"}
	}

	if in.SeriesType.Id == nil {
		return nil, &MapError{Kind: "series", ID: id, Field: "SeriesType.Id"}
	}

	if in.IssueCount == nil {
		return nil, &MapError{Kind: "series", ID: id, Field: "IssueCount"}
	}

	resourceURL, err := url.Parse(*in.ResourceUrl)
	if err != nil {
		return nil, &MapError{Kind: "series", ID: id, Field: "ResourceUrl", Err: err}
	}

	var imprint *Reference

	if in.Imprint != nil {
		if in.Imprint.Id == nil {
			return nil, &MapError{Kind: "series", ID: id, Field: "Imprint.Id"}
		}

		imprint = &Reference{
			ID:   *in.Imprint.Id,
			Name: in.Imprint.Name,
		}
	}

	var altNames []string

	if in.AltNames != nil {
		altNames = *in.AltNames
	}

	var genres []Reference

	if in.Genres != nil {
		genres = make([]Reference, 0, len(*in.Genres))

		for _, g := range *in.Genres {
			if g.Id == nil {
				return nil, &MapError{Kind: "series", ID: id, Field: "Genres[].Id"}
			}

			genres = append(genres, Reference{
				ID:   *g.Id,
				Name: g.Name,
			})
		}
	}

	var associated []Reference

	if in.Associated != nil {
		associated = make([]Reference, 0, len(*in.Associated))

		for _, a := range *in.Associated {
			if a.Id == nil {
				return nil, &MapError{Kind: "series", ID: id, Field: "Associated[].Id"}
			}

			associated = append(associated, Reference{
				ID:   *a.Id,
				Name: a.Series,
			})
		}
	}

	return &Series{
		ID:               id,
		Name:             in.Name,
		AlternativeNames: altNames,
		SortName:         in.SortName,
		Volume:           in.Volume,
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
		YearEnded:             nullableToPtr(in.YearEnd),
		Description:           in.Desc,
		IssueCount:            *in.IssueCount,
		Genres:                genres,
		Associated:            associated,
		ComicVineID:           nullableToPtr(in.CvId),
		GrandComicsDatabaseID: nullableToPtr(in.GcdId),
		ResourceURL:           *resourceURL,
		Modified:              *in.Modified,
	}, nil
}

func seriesListMapper(in internal.SeriesList) (*SeriesList, error) {
	if in.Id == nil {
		return nil, &MapError{Kind: "series", Field: "Id"}
	}

	id := *in.Id

	if in.Modified == nil {
		return nil, &MapError{Kind: "series", ID: id, Field: "Modified"}
	}

	if in.IssueCount == nil {
		return nil, &MapError{Kind: "series", ID: id, Field: "IssueCount"}
	}

	return &SeriesList{
		ID:         id,
		Name:       in.Series,
		YearBegan:  in.YearBegan,
		YearEnded:  nullableToPtr(in.YearEnd),
		Volume:     in.Volume,
		IssueCount: *in.IssueCount,
		Modified:   *in.Modified,
	}, nil
}
