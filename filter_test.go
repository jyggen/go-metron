package metron_test

import (
	"context"
	"os/exec"
	"strings"
	"testing"
	"time"

	"cloud.google.com/go/civil"
	"github.com/jyggen/go-metron"
	"github.com/stretchr/testify/require"
)

// issueQuery ranges a single-page issue listing and returns the URL requested.
// The mock serves an empty results page, so the iterator stops after one call.
func issueQuery(t *testing.T, filters *metron.IssueFilters) string {
	t.Helper()

	var got string

	c := newTestClient(t, []requestMock{
		{captureURL: &got, responseBody: `{"count":0,"next":null,"previous":null,"results":[]}`},
	})

	for _, err := range c.Issues(context.Background(), filters) {
		require.NoError(t, err)
	}

	return got
}

func TestIssueFiltersQueryString(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name     string
		filters  *metron.IssueFilters
		expected string
	}{
		{
			name:     "nil sends only the page",
			filters:  nil,
			expected: "https://metron.cloud/api/issue/?page=1",
		},
		{
			// The zero civil.Date normalises to a valid-looking -0001-11-30, so the
			// params method has to test IsValid.
			name:     "zero value is indistinguishable from nil",
			filters:  &metron.IssueFilters{},
			expected: "https://metron.cloud/api/issue/?page=1",
		},
		{
			name:     "string",
			filters:  &metron.IssueFilters{PublisherName: "marvel"},
			expected: "https://metron.cloud/api/issue/?page=1&publisher_name=marvel",
		},
		{
			name:     "int",
			filters:  &metron.IssueFilters{SeriesID: 793},
			expected: "https://metron.cloud/api/issue/?page=1&series_id=793",
		},
		{
			name:     "time.Month renders as a bare number",
			filters:  &metron.IssueFilters{CoverMonth: time.March},
			expected: "https://metron.cloud/api/issue/?cover_month=3&page=1",
		},
		{
			name:     "int widened to float32",
			filters:  &metron.IssueFilters{CoverYear: 2024},
			expected: "https://metron.cloud/api/issue/?cover_year=2024&page=1",
		},
		{
			name:     "civil.Date",
			filters:  &metron.IssueFilters{StoreDate: civil.Date{Year: 2021, Month: time.June, Day: 7}},
			expected: "https://metron.cloud/api/issue/?page=1&store_date=2021-06-07",
		},
		{
			name: "inclusive date bounds",
			filters: &metron.IssueFilters{
				StoreDateFrom: civil.Date{Year: 2021, Month: time.June, Day: 7},
				StoreDateTo:   civil.Date{Year: 2021, Month: time.June, Day: 13},
			},
			expected: "https://metron.cloud/api/issue/?page=1&store_date_range_after=2021-06-07&store_date_range_before=2021-06-13",
		},
		{
			// Deliberately not normalised to UTC, unlike the date filters.
			name:     "time.Time keeps its offset",
			filters:  &metron.IssueFilters{ModifiedAfter: parseTime(t, "2024-03-07T12:36:10-05:00")},
			expected: "https://metron.cloud/api/issue/?modified_gt=2024-03-07T12%3A36%3A10-05%3A00&page=1",
		},
		{
			name:     "slice joins with commas",
			filters:  &metron.IssueFilters{RoleIDs: []int{1, 2, 3}},
			expected: "https://metron.cloud/api/issue/?page=1&role_id=1%2C2%2C3",
		},
		{
			name:     "missing true",
			filters:  &metron.IssueFilters{MissingComicVineID: new(true)},
			expected: "https://metron.cloud/api/issue/?missing_cv_id=true&page=1",
		},
		{
			// false is a distinct query: "has a Comic Vine ID".
			name:     "missing false",
			filters:  &metron.IssueFilters{MissingComicVineID: new(false)},
			expected: "https://metron.cloud/api/issue/?missing_cv_id=false&page=1",
		},
		{
			name:     "cross-resource filter",
			filters:  &metron.IssueFilters{CharacterID: 42, CoverYear: 2024},
			expected: "https://metron.cloud/api/issue/?character_id=42&cover_year=2024&page=1",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tc.expected, issueQuery(t, tc.filters))
		})
	}
}

func TestSeriesFiltersQueryString(t *testing.T) {
	t.Parallel()

	var got string

	c := newTestClient(t, []requestMock{
		{captureURL: &got, responseBody: `{"count":0,"next":null,"previous":null,"results":[]}`},
	})

	for _, err := range c.Series(context.Background(), &metron.SeriesFilters{
		PublisherID: 2,
		Status:      metron.SeriesOngoing,
		Search:      "spider",
	}) {
		require.NoError(t, err)
	}

	require.Equal(t, "https://metron.cloud/api/series/?page=1&publisher_id=2&q=spider&status=4", got)
}

func TestFiltersReachTheQueryStringAcrossPages(t *testing.T) {
	t.Parallel()

	c := newTestClient(t, []requestMock{
		{expectedURL: "https://metron.cloud/api/role/?name=Writer&page=1", responseBodyFixture: "fixtures/role_list_1.json"},
		{expectedURL: "https://metron.cloud/api/role/?name=Writer&page=2", responseBodyFixture: "fixtures/role_list_2.json"},
	})

	var count int

	for _, err := range c.Roles(context.Background(), &metron.RoleFilters{Name: "Writer"}) {
		require.NoError(t, err)

		count++
	}

	require.Positive(t, count)
}

// TestInvalidUsageDoesNotCompile covers the calls the compiler must reject,
// which no compiling test can. Each testdata/negative package pairs a Bad
// function with a Good one; asserting a single error keeps the Good one honest.
func TestInvalidUsageDoesNotCompile(t *testing.T) {
	t.Parallel()

	if _, err := exec.LookPath("go"); err != nil {
		t.Skip("needs the go toolchain")
	}

	testCases := []struct {
		pkg      string
		expected string
	}{
		{
			pkg:      "conditional_on_list",
			expected: "metron.ConditionalOption does not implement metron.RequestOption",
		},
		{
			pkg:      "wrong_filters_type",
			expected: "cannot use &metron.IssueFilters{} (value of type *metron.IssueFilters) as *metron.RoleFilters value",
		},
		{
			pkg:      "unknown_filter_field",
			expected: "unknown field PublisherID in struct literal of type metron.ArcFilters",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.pkg, func(t *testing.T) {
			t.Parallel()

			out, err := exec.CommandContext(t.Context(), "go", "build", "./testdata/negative/"+tc.pkg).CombinedOutput()

			require.Error(t, err, "expected a compile error, got none:\n%s", out)
			require.Contains(t, string(out), tc.expected)
			require.Equal(t, 1, strings.Count(string(out), "testdata/negative/"+tc.pkg+"/main.go:"),
				"expected exactly one error, so the Good function is known to compile:\n%s", out)
		})
	}
}

func TestFiltersAreReusable(t *testing.T) {
	t.Parallel()

	// The same struct drives any number of listings; nothing consumes it.
	filters := &metron.IssueFilters{RoleIDs: []int{1, 2}}
	expected := "https://metron.cloud/api/issue/?page=1&role_id=1%2C2"

	require.Equal(t, expected, issueQuery(t, filters))
	require.Equal(t, expected, issueQuery(t, filters))
}
