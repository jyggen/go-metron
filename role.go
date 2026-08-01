package metron

import (
	"context"
	"fmt"
	"iter"

	"github.com/jyggen/go-metron/internal"
)

// RoleList is a credit role as it appears in list responses.
type RoleList struct {
	ID   int
	Name string
}

// Roles returns an iterator over all roles.
func (c *Client) Roles(ctx context.Context, filters ...Filter) iter.Seq2[*RoleList, error] {
	params := &internal.ApiRoleListParams{}

	for _, f := range filters {
		f(params)
	}

	return paginate[internal.PaginatedRoleList](ctx, c, "role", c.client.ApiRoleList, roleListMapper, params)
}

func roleListMapper(in internal.Role) (*RoleList, error) {
	if in.Id == nil {
		return nil, fmt.Errorf("role: nil Id")
	}

	return &RoleList{
		ID:   *in.Id,
		Name: in.Name,
	}, nil
}
