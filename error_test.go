package metron_test

import (
	"context"
	"net/http"
	"strings"
	"testing"

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
