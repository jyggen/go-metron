package metron

import (
	"context"
	"iter"

	"github.com/jyggen/go-metron/internal"
)

type SeriesTypeList struct {
	ID   int
	Name string
}

// SeriesTypes returns a list of all the series types.
func (c *Client) SeriesTypes(ctx context.Context, filters ...Filter) iter.Seq2[*SeriesTypeList, error] {
	params := &internal.ApiSeriesTypeListParams{}

	for _, f := range filters {
		f(params)
	}

	return newPaginate[internal.PaginatedSeriesTypeList](ctx, c.client.ApiSeriesTypeList, seriesTypeMapper, params)
}

func seriesTypeMapper(in internal.SeriesType) (*SeriesTypeList, error) {
	return &SeriesTypeList{
		ID:   *in.Id,
		Name: in.Name,
	}, nil
}
