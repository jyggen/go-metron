package metron_test

import (
	"context"
	"net/http"
	"strings"
	"testing"
	"time"

	"codeberg.org/jyggen/go-httpkit"
	"github.com/jyggen/go-metron"
	"github.com/stretchr/testify/require"
)

const issueURL = "https://metron.cloud/api/issue/1/"

func TestAPIError(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name           string
		responseStatus int
		responseBody   string
		expectedError  string
	}{
		{
			name:           "not found",
			responseStatus: http.StatusNotFound,
			responseBody:   `{"detail":"Not found."}`,
			expectedError:  `metron: unexpected status code: 404: {"detail":"Not found."}`,
		},
		{
			name:           "unauthorized",
			responseStatus: http.StatusUnauthorized,
			responseBody:   `{"detail":"Invalid token."}`,
			expectedError:  `metron: unexpected status code: 401: {"detail":"Invalid token."}`,
		},
		{
			name:           "server error with empty body",
			responseStatus: http.StatusInternalServerError,
			responseBody:   "",
			expectedError:  "metron: unexpected status code: 500",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			c := newTestClient(t, []requestMock{
				{expectedURL: issueURL, responseStatus: tc.responseStatus, responseBody: tc.responseBody},
			})

			_, err := c.IssueByID(context.Background(), 1)

			require.Error(t, err)

			var apiErr *metron.APIError

			require.ErrorAs(t, err, &apiErr)
			require.Equal(t, tc.responseStatus, apiErr.StatusCode)
			require.Equal(t, tc.responseBody, string(apiErr.Body))
			require.Equal(t, tc.expectedError, apiErr.Error())

			// Callers can match on status alone.
			require.ErrorIs(t, err, &metron.APIError{StatusCode: tc.responseStatus})
			require.NotErrorIs(t, err, &metron.APIError{StatusCode: http.StatusTeapot})
		})
	}
}

func TestAPIErrorBodyIsTruncated(t *testing.T) {
	t.Parallel()

	c := newTestClient(t, []requestMock{
		{
			expectedURL:    issueURL,
			responseStatus: http.StatusInternalServerError,
			responseBody:   strings.Repeat("a", 8<<10),
		},
	})

	_, err := c.IssueByID(context.Background(), 1)

	var apiErr *metron.APIError

	require.ErrorAs(t, err, &apiErr)
	require.Len(t, apiErr.Body, 4<<10)
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
				// Without a usable Retry-After there is nothing to wait on, so
				// the response surfaces as a plain APIError instead.
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
