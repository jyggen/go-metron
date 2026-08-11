package metron

import (
	"context"
	"iter"

	"github.com/jyggen/go-metron/internal"
)

// SeriesTypeList is a series type as it appears in list responses.
type SeriesTypeList struct {
	ID   int
	Name string
}

// SeriesTypes returns an iterator over all series types.
func (c *Client) SeriesTypes(ctx context.Context, filters *SeriesTypeFilters, opts ...RequestOption) iter.Seq2[*SeriesTypeList, error] {
	return paginate[internal.PaginatedSeriesTypeList](ctx, c, c.client.ApiSeriesTypeList, seriesTypeMapper, filters.params, everyPage(requestEditors(opts)))
}

func seriesTypeMapper(in internal.SeriesType) (*SeriesTypeList, error) {
	if in.Id == nil {
		return nil, &MapError{Kind: "series_type", Field: "Id"}
	}

	return &SeriesTypeList{
		ID:   *in.Id,
		Name: in.Name,
	}, nil
}
