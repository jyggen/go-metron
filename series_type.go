package metron

import "context"

type SeriesTypeList struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func (c *Client) SeriesTypes(ctx context.Context) func(func(SeriesTypeList, error) bool) {
	return paginate[SeriesTypeList](ctx, c, "series_type/")
}
