package metron

import "context"

type RoleList struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func (c *Client) Roles(ctx context.Context, filters ...Filter) func(func(RoleList, error) bool) {
	return paginate[RoleList](ctx, c, "role/", filters...)
}
