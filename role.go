package metron

import (
	"context"
	"iter"

	"github.com/jyggen/go-metron/internal"
)

// RoleList is a credit role as it appears in list responses.
type RoleList struct {
	ID   int
	Name string
}

// Roles returns an iterator over all roles.
func (c *Client) Roles(ctx context.Context, filters *RoleFilters, opts ...RequestOption) iter.Seq2[*RoleList, error] {
	return paginate[internal.PaginatedRoleList](ctx, c, c.client.ApiRoleList, roleListMapper, filters.params, everyPage(requestEditors(opts)))
}

func roleListMapper(in internal.Role) (*RoleList, error) {
	if in.Id == nil {
		return nil, &MapError{Kind: "role", Field: "Id"}
	}

	return &RoleList{
		ID:   *in.Id,
		Name: in.Name,
	}, nil
}
