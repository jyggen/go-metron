package metron

import (
	"context"
	"iter"

	"github.com/jyggen/go-metron/internal"
)

type RoleList struct {
	ID   int
	Name string
}

// Roles returns a list of all the roles.
func (c *Client) Roles(ctx context.Context, filters ...Filter) iter.Seq2[*RoleList, error] {
	params := &internal.ApiRoleListParams{}

	for _, f := range filters {
		f(params)
	}

	return newPaginate[internal.PaginatedRoleList](ctx, c.client.ApiRoleList, roleListMapper, params)
}

func roleListMapper(in internal.Role) (*RoleList, error) {
	return &RoleList{
		ID:   *in.Id,
		Name: in.Name,
	}, nil
}
