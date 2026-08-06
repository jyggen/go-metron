package metron_test

import (
	"context"
	"iter"
	"testing"

	"github.com/jyggen/go-metron"
	"github.com/stretchr/testify/require"
)

// TestFilterUnsupported checks an inapplicable filter is reported, not dropped.
// ApiRoleListParams takes only name and modified_gt, so ByPublisherID cannot.
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
		require.Equal(t, "Roles", filterErr.Endpoint)
		require.Equal(t, "metron: ByPublisherID does not apply to Roles", filterErr.Error())

		// Callers can match a specific filter, a specific endpoint, or both.
		require.ErrorIs(t, err, &metron.FilterError{Filter: "ByPublisherID"})
		require.ErrorIs(t, err, &metron.FilterError{Endpoint: "Roles"})
		require.ErrorIs(t, err, &metron.FilterError{Filter: "ByPublisherID", Endpoint: "Roles"})
		require.NotErrorIs(t, err, &metron.FilterError{Filter: "ByName"})
		require.NotErrorIs(t, err, &metron.FilterError{Endpoint: "Issues"})
	}

	require.Equal(t, 1, iterations)
}

// TestFilterSupported checks a filter reaches the query string.
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

// TestFilterUnsupportedAmongSupported checks a bad filter wins over good ones.
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

// TestFilterUnsupportedNamesTheCallingMethod checks the endpoint is the method
// that rejected the filter, not a constant. The two params structs here are
// identical, so only the caller distinguishes them.
func TestFilterUnsupportedNamesTheCallingMethod(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		endpoint string
		list     func(*metron.Client) iter.Seq2[any, error]
	}{
		{
			endpoint: "Roles",
			list: func(c *metron.Client) iter.Seq2[any, error] {
				return anyIter(c.Roles(context.Background(), metron.ByPublisherID(2)))
			},
		},
		{
			endpoint: "SeriesTypes",
			list: func(c *metron.Client) iter.Seq2[any, error] {
				return anyIter(c.SeriesTypes(context.Background(), metron.ByPublisherID(2)))
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.endpoint, func(t *testing.T) {
			t.Parallel()

			c := newTestClient(t, nil)

			for _, err := range tc.list(c) {
				require.EqualError(t, err, "metron: ByPublisherID does not apply to "+tc.endpoint)
				require.ErrorIs(t, err, &metron.FilterError{Endpoint: tc.endpoint})
			}
		})
	}
}

// anyIter erases a list method's element type so iterators over different
// resources can share a table.
func anyIter[T any](seq iter.Seq2[T, error]) iter.Seq2[any, error] {
	return func(yield func(any, error) bool) {
		for v, err := range seq {
			if !yield(v, err) {
				return
			}
		}
	}
}
