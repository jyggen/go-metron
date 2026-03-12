package metron

import (
	"context"
	"fmt"
	"iter"

	"github.com/jyggen/go-metron/internal"
)

type SeriesTypeList struct {
	ID   int
	Name string
}

// SeriesTypes returns an iterator over all series types.
func (c *Client) SeriesTypes(ctx context.Context, filters ...Filter) iter.Seq2[*SeriesTypeList, error] {
	params := &internal.ApiSeriesTypeListParams{}

	for _, f := range filters {
		f(params)
	}

	return newPaginate[internal.PaginatedSeriesTypeList](ctx, c.cache, c.maxRetries, "series_type", c.client.ApiSeriesTypeList, seriesTypeMapper, params)
}

func seriesTypeMapper(in internal.SeriesType) (*SeriesTypeList, error) {
	if in.Id == nil {
		return nil, fmt.Errorf("series_type: nil Id")
	}

	return &SeriesTypeList{
		ID:   *in.Id,
		Name: in.Name,
	}, nil
}
