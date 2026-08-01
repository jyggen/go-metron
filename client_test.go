package metron_test

import (
	"context"
	"embed"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"testing"
	"time"

	"cloud.google.com/go/civil"
	"github.com/jyggen/go-metron"
	"github.com/stretchr/testify/require"
)

//go:embed fixtures/*
var fs embed.FS

type roundTripFunc func(req *http.Request) *http.Response

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req), nil
}

type requestMock struct {
	expectedURL         string
	responseBodyFixture string
}

type testCase[T any] struct {
	id       int
	expected T
}

func testList[T any](
	t *testing.T,
	kind string,
	method func(*metron.Client, context.Context, ...metron.Filter) func(func(T, error) bool),
	testCases []testCase[T],
) {
	c := newTestClient(t, []requestMock{
		{fmt.Sprintf("https://metron.cloud/api/%s/", kind), fmt.Sprintf("fixtures/%s_list_1.json", kind)},
		{fmt.Sprintf("https://metron.cloud/api/%s/?page=2", kind), fmt.Sprintf("fixtures/%s_list_2.json", kind)},
	})

	resources := make([]T, 0, 4)

	for res, err := range method(c, context.Background()) {
		require.NoError(t, err)
		resources = append(resources, res)
	}

	require.Len(t, resources, len(testCases))

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("%T #%d", tc.expected, tc.id), func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tc.expected, resources[i])
		})
	}
}

func testListByID[T any](
	t *testing.T,
	kind string,
	id int,
	listKind string,
	method func(*metron.Client, context.Context, int) func(func(T, error) bool),
	testCases []testCase[T],
) {
	c := newTestClient(t, []requestMock{
		{
			fmt.Sprintf("https://metron.cloud/api/%s/%d/%s_list/", kind, id, listKind),
			fmt.Sprintf("fixtures/%s_%d_%s_list_1.json", kind, id, listKind),
		},
		{
			fmt.Sprintf("https://metron.cloud/api/%s/%d/%s_list/?page=2", kind, id, listKind),
			fmt.Sprintf("fixtures/%s_%d_%s_list_2.json", kind, id, listKind),
		},
	})

	resources := make([]T, 0, 4)

	for res, err := range method(c, context.Background(), id) {
		require.NoError(t, err)
		resources = append(resources, res)
	}

	require.Len(t, resources, len(testCases))

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("%T #%d", tc.expected, tc.id), func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tc.expected, resources[i])
		})
	}
}

func testByID[T any](
	t *testing.T,
	kind string,
	method func(*metron.Client, context.Context, int) (T, error),
	testCases []testCase[T],
) {
	for _, tc := range testCases {
		t.Run(fmt.Sprintf("%T #%d", tc.expected, tc.id), func(t *testing.T) {
			t.Parallel()

			c := newTestClient(t, []requestMock{
				{
					fmt.Sprintf("https://metron.cloud/api/%s/%d/", kind, tc.id),
					fmt.Sprintf("fixtures/%s_%d.json", kind, tc.id),
				},
			})

			v, err := method(c, context.Background(), tc.id)

			require.NoError(t, err)
			require.Equal(t, tc.expected, v)
		})
	}
}

func newTestClient(t *testing.T, mocks []requestMock) *metron.Client {
	return metron.NewClient(metron.WithClient(&http.Client{
		Transport: roundTripFunc(func(req *http.Request) *http.Response {
			m := mocks[0]
			mocks = append(mocks[:0], mocks[1:]...)

			require.Equal(t, m.expectedURL, req.URL.String())

			username, password, _ := req.BasicAuth()

			require.Equal(t, "username", username)
			require.Equal(t, "password", password)

			f, err := fs.Open(m.responseBodyFixture)
			require.NoError(t, err)

			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       f,
				Header:     make(http.Header),
			}
		}),
	}), metron.WithAuthentication("username", "password"))
}

func TestRateLimitWaitsForHeaderReset(t *testing.T) {
	t.Parallel()

	// The reset header only has second-level precision, so the window must be
	// large enough that truncating it down to the nearest second still leaves
	// a reliably measurable wait.
	const window = 3 * time.Second

	reset := time.Now().Add(window)
	var callTimes []time.Time

	c := metron.NewClient(metron.WithClient(&http.Client{
		Transport: roundTripFunc(func(req *http.Request) *http.Response {
			callTimes = append(callTimes, time.Now())

			f, err := fs.Open(fmt.Sprintf("fixtures/series_%s.json", req.URL.Path[len("/api/series/"):len(req.URL.Path)-1]))
			require.NoError(t, err)

			h := make(http.Header)
			h.Set("X-RateLimit-Burst-Remaining", "0")
			h.Set("X-RateLimit-Burst-Reset", strconv.FormatInt(reset.Unix(), 10))
			h.Set("X-RateLimit-Sustained-Remaining", "100")
			h.Set("X-RateLimit-Sustained-Reset", strconv.FormatInt(time.Now().Add(24*time.Hour).Unix(), 10))

			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       f,
				Header:     h,
			}
		}),
	}), metron.WithAuthentication("username", "password"))

	_, err := c.SeriesByID(context.Background(), 793)
	require.NoError(t, err)

	_, err = c.SeriesByID(context.Background(), 3371)
	require.NoError(t, err)

	require.Len(t, callTimes, 2)

	elapsed := callTimes[1].Sub(callTimes[0])
	require.GreaterOrEqual(t, elapsed, window-1200*time.Millisecond)
	require.LessOrEqual(t, elapsed, window+500*time.Millisecond)
}

func parseDate(t *testing.T, dateString string) civil.Date {
	v, err := civil.ParseDate(dateString)

	require.NoError(t, err)

	return v
}

func parseTime(t *testing.T, timeString string) time.Time {
	v, err := time.Parse(time.RFC3339, timeString)

	require.NoError(t, err)

	return v
}

func parseURL(t *testing.T, urlString string) metron.URL {
	v, err := url.Parse(urlString)

	require.NoError(t, err)

	return metron.URL{URL: v}
}

func asReference[T any](v T) *T {
	return &v
}
