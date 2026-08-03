package metron_test

import (
	"context"
	"testing"

	"github.com/jyggen/go-metron"
	"github.com/stretchr/testify/require"
)

// TestFilterUnsupported checks that a filter the endpoint does not support is
// reported rather than silently dropped. ApiRoleListParams only accepts name
// and modified_gt, so ByPublisherID cannot apply to Roles.
func TestFilterUnsupported(t *testing.T) {
	t.Parallel()

	// No mocks: the filter must fail before any request is made. If that ever
	// regresses, the transport indexes an empty slice and the test panics.
	c := newTestClient(t, nil)

	var iterations int

	for role, err := range c.Roles(context.Background(), metron.ByPublisherID(2)) {
		iterations++

		require.Nil(t, role)
		require.Error(t, err)

		var filterErr *metron.FilterError

		require.ErrorAs(t, err, &filterErr)
		require.Equal(t, "ByPublisherID", filterErr.Filter)
		require.Equal(t, "metron: ByPublisherID does not apply to this endpoint", filterErr.Error())

		// Callers can match a specific filter.
		require.ErrorIs(t, err, &metron.FilterError{Filter: "ByPublisherID"})
		require.NotErrorIs(t, err, &metron.FilterError{Filter: "ByName"})
	}

	require.Equal(t, 1, iterations)
}

// TestFilterSupported checks the happy path still reaches the API with the
// filter applied to the query string.
func TestFilterSupported(t *testing.T) {
	t.Parallel()

	c := newTestClient(t, []requestMock{
		{
			expectedURL:         "https://metron.cloud/api/role/?name=Writer&page=1",
			responseBodyFixture: "fixtures/role_list_1.json",
		},
		{
			expectedURL:         "https://metron.cloud/api/role/?name=Writer&page=2",
			responseBodyFixture: "fixtures/role_list_2.json",
		},
	})

	for _, err := range c.Roles(context.Background(), metron.ByName("Writer")) {
		require.NoError(t, err)
	}
}

// TestFilterUnsupportedAmongSupported checks the first unsupported filter wins
// even when valid ones are present.
func TestFilterUnsupportedAmongSupported(t *testing.T) {
	t.Parallel()

	c := newTestClient(t, nil)

	for _, err := range c.Roles(
		context.Background(),
		metron.ByName("Writer"),
		metron.ByPublisherID(2),
	) {
		require.ErrorIs(t, err, &metron.FilterError{Filter: "ByPublisherID"})
	}
}
