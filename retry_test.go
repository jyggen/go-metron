package metron_test

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"codeberg.org/jyggen/go-httpkit"
	"github.com/jyggen/go-metron"
	"github.com/stretchr/testify/require"
)

// failingRoundTripper is a transport that can fail, which roundTripFunc cannot.
type failingRoundTripper func(req *http.Request) (*http.Response, error)

func (f failingRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

// retryClient returns a client whose transport replies from responses in order,
// along with the request counter. A nil entry means a network failure.
func retryClient(t *testing.T, retries int, responses []*http.Response) (*metron.Client, *atomic.Int64) {
	t.Helper()

	var count atomic.Int64

	c, err := metron.NewClient("foobar",
		metron.WithRetry(retries),
		metron.WithClient(&http.Client{
			Transport: failingRoundTripper(func(_ *http.Request) (*http.Response, error) {
				idx := int(count.Add(1)) - 1
				require.Less(t, idx, len(responses), "unexpected extra HTTP request")

				if responses[idx] == nil {
					return nil, errors.New("dial tcp: connection refused")
				}

				return responses[idx], nil
			}),
		}),
	)
	require.NoError(t, err)

	return c, &count
}

func status(code int) *http.Response {
	return &http.Response{
		StatusCode: code,
		Body:       io.NopCloser(strings.NewReader(`{"detail":"nope"}`)),
		Header:     make(http.Header),
	}
}

func fixture(t *testing.T, name string) *http.Response {
	t.Helper()

	f, err := fs.Open(name)
	require.NoError(t, err)

	return &http.Response{StatusCode: http.StatusOK, Body: f, Header: make(http.Header)}
}

func TestRetryTransient(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name             string
		responses        []*http.Response
		expectedRequests int64
		expectedStatus   int
	}{
		{
			name:             "503 is retried",
			responses:        []*http.Response{status(503), status(503)},
			expectedRequests: 2,
			expectedStatus:   503,
		},
		{
			name:             "502 is retried",
			responses:        []*http.Response{status(502), status(502)},
			expectedRequests: 2,
			expectedStatus:   502,
		},
		{
			name:             "504 is retried",
			responses:        []*http.Response{status(504), status(504)},
			expectedRequests: 2,
			expectedStatus:   504,
		},
		{
			name:             "500 is not retried",
			responses:        []*http.Response{status(500)},
			expectedRequests: 1,
			expectedStatus:   500,
		},
		{
			name:             "501 is not retried",
			responses:        []*http.Response{status(501)},
			expectedRequests: 1,
			expectedStatus:   501,
		},
		{
			name:             "404 is not retried",
			responses:        []*http.Response{status(404)},
			expectedRequests: 1,
			expectedStatus:   404,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			c, count := retryClient(t, 1, tc.responses)

			_, err := c.ArcByID(context.Background(), 659)

			require.ErrorIs(t, err, &metron.APIError{StatusCode: tc.expectedStatus})
			require.Equal(t, tc.expectedRequests, count.Load())
		})
	}
}

func TestRetryNetworkError(t *testing.T) {
	t.Parallel()

	c, count := retryClient(t, 1, []*http.Response{nil, nil})

	_, err := c.ArcByID(context.Background(), 659)

	require.Error(t, err)
	require.Equal(t, int64(2), count.Load())
}

func TestRetrySucceedsOnSecondAttempt(t *testing.T) {
	t.Parallel()

	c, count := retryClient(t, 1, []*http.Response{status(503), fixture(t, "fixtures/arc_659.json")})

	arc, err := c.ArcByID(context.Background(), 659)

	require.NoError(t, err)
	require.Equal(t, 659, arc.ID)
	require.Equal(t, int64(2), count.Load())
}

// TestTransientRetryNotAppliedToWrites covers the non-idempotent path. Writes
// still retry on rate limits, which TestRetryAfter covers.
func TestTransientRetryNotAppliedToWrites(t *testing.T) {
	t.Parallel()

	c, count := retryClient(t, 3, []*http.Response{status(503)})

	_, err := c.Scrobble(context.Background(), 12345)

	require.ErrorIs(t, err, &metron.APIError{StatusCode: 503})
	require.Equal(t, int64(1), count.Load())
}

func TestRetryAfter(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name        string
		retryAfter  func() string
		expectRetry bool
		minWait     time.Duration
		maxWait     time.Duration
	}{
		{
			name:        "delta seconds",
			retryAfter:  func() string { return "30" },
			expectRetry: true,
			minWait:     30 * time.Second,
			maxWait:     30 * time.Second,
		},
		{
			name: "http date in the future",
			retryAfter: func() string {
				return time.Now().Add(2 * time.Second).UTC().Format(http.TimeFormat)
			},
			expectRetry: true,
			minWait:     1,
			maxWait:     2 * time.Second,
		},
		{
			name: "http date in the past",
			retryAfter: func() string {
				return time.Now().Add(-time.Hour).UTC().Format(http.TimeFormat)
			},
			expectRetry: true,
			minWait:     0,
			maxWait:     0,
		},
		{
			name:        "malformed",
			retryAfter:  func() string { return "later please" },
			expectRetry: false,
		},
		{
			name:        "absent",
			retryAfter:  func() string { return "" },
			expectRetry: false,
		},
		{
			name:        "negative",
			retryAfter:  func() string { return "-5" },
			expectRetry: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			header := make(http.Header)

			if v := tc.retryAfter(); v != "" {
				header.Set("Retry-After", v)
			}

			c := newTestClient(t, []requestMock{
				{
					expectedURL:    issueURL,
					responseStatus: http.StatusTooManyRequests,
					responseHeader: header,
					responseBody:   `{"detail":"Request was throttled."}`,
				},
			})

			_, err := c.IssueByID(context.Background(), 1)

			require.Error(t, err)

			var retryErr *httpkit.RetryAfterError

			if !tc.expectRetry {
				// Nothing to wait on, so it surfaces as a plain APIError.
				require.NotErrorAs(t, err, &retryErr)
				require.ErrorIs(t, err, &metron.APIError{StatusCode: http.StatusTooManyRequests})

				return
			}

			require.ErrorAs(t, err, &retryErr)
			require.GreaterOrEqual(t, retryErr.RetryAfter(), tc.minWait)
			require.LessOrEqual(t, retryErr.RetryAfter(), tc.maxWait)
		})
	}
}
