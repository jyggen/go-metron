package metron_test

import (
	"testing"

	"github.com/jyggen/go-metron"
)

func TestRoles(t *testing.T) {
	t.Parallel()
	testList(t, "role", (*metron.Client).Roles, []testCase[metron.RoleList]{
		{
			id: 1,
			expected: metron.RoleList{
				ID:   1,
				Name: "Writer",
			},
		},
		{
			id: 22,
			expected: metron.RoleList{
				ID:   22,
				Name: "Script",
			},
		},
		{
			id: 30,
			expected: metron.RoleList{
				ID:   30,
				Name: "Story",
			},
		},
		{
			id: 23,
			expected: metron.RoleList{
				ID:   23,
				Name: "Plot",
			},
		},
	})
}
