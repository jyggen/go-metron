package metron_test

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"testing"

	"github.com/jyggen/go-metron"
	"github.com/stretchr/testify/require"
)

const issueURL = "https://metron.cloud/api/issue/1/"

var errUnusable = errors.New("unusable")

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
			expectedError:  `metron: GET ` + issueURL + `: unexpected status code: 404: {"detail":"Not found."}`,
		},
		{
			name:           "unauthorized",
			responseStatus: http.StatusUnauthorized,
			responseBody:   `{"detail":"Invalid token."}`,
			expectedError:  `metron: GET ` + issueURL + `: unexpected status code: 401: {"detail":"Invalid token."}`,
		},
		{
			name:           "server error with empty body",
			responseStatus: http.StatusInternalServerError,
			responseBody:   "",
			expectedError:  "metron: GET " + issueURL + ": unexpected status code: 500",
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

			// The request is captured from the outgoing *http.Request, so it is
			// present even though the mock transport never sets res.Request.
			require.Equal(t, http.MethodGet, apiErr.Method)
			require.Equal(t, issueURL, apiErr.URL.String())

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

func TestAPIErrorWithoutRequest(t *testing.T) {
	t.Parallel()

	// A caller-constructed APIError, as used for errors.Is matching, has no
	// request to name.
	err := &metron.APIError{StatusCode: http.StatusNotFound}

	require.Equal(t, "metron: unexpected status code: 404", err.Error())
}

func TestAPIErrorNamesTheRequestBehindAnIterator(t *testing.T) {
	t.Parallel()

	c := newTestClient(t, []requestMock{
		{
			expectedURL:         "https://metron.cloud/api/issue/?page=1",
			responseBodyFixture: "fixtures/issue_list_1.json",
		},
		{
			expectedURL:    "https://metron.cloud/api/issue/?page=2",
			responseStatus: http.StatusServiceUnavailable,
		},
	})

	var err error

	for _, iterErr := range c.Issues(context.Background()) {
		if iterErr != nil {
			err = iterErr

			break
		}
	}

	var apiErr *metron.APIError

	require.ErrorAs(t, err, &apiErr)

	// The page is the context the caller cannot supply: the iterator hides
	// pagination, so nothing on their side knows how far it had got.
	require.Equal(t, "https://metron.cloud/api/issue/?page=2", apiErr.URL.String())
	require.Contains(t, err.Error(), "metron: GET https://metron.cloud/api/issue/?page=2: unexpected status code: 503")
}

func TestMapErrorMessage(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name     string
		err      *metron.MapError
		expected string
	}{
		{
			name:     "missing field",
			err:      &metron.MapError{Kind: "issue", ID: 2558, Field: "RatingCount"},
			expected: "metron: issue 2558: nil RatingCount",
		},
		{
			name:     "nested field",
			err:      &metron.MapError{Kind: "issue", ID: 2558, Field: "Credits[].Creator"},
			expected: "metron: issue 2558: nil Credits[].Creator",
		},
		{
			name:     "unusable field",
			err:      &metron.MapError{Kind: "issue", ID: 2558, Field: "Price", Err: errUnusable},
			expected: "metron: issue 2558: Price: unusable",
		},
		{
			name:     "record without an ID",
			err:      &metron.MapError{Kind: "issue", Field: "Id"},
			expected: "metron: issue: nil Id",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tc.expected, tc.err.Error())
		})
	}
}

func TestMapErrorMatching(t *testing.T) {
	t.Parallel()

	err := &metron.MapError{Kind: "issue", ID: 2558, Field: "RatingCount", Err: errUnusable}

	// Zero fields on the target are wildcards, so a caller can match as broadly
	// or as narrowly as it needs.
	require.ErrorIs(t, err, &metron.MapError{})
	require.ErrorIs(t, err, &metron.MapError{Kind: "issue"})
	require.ErrorIs(t, err, &metron.MapError{ID: 2558})
	require.ErrorIs(t, err, &metron.MapError{Field: "RatingCount"})
	require.ErrorIs(t, err, &metron.MapError{Kind: "issue", ID: 2558, Field: "RatingCount"})

	require.NotErrorIs(t, err, &metron.MapError{Kind: "series"})
	require.NotErrorIs(t, err, &metron.MapError{ID: 1})
	require.NotErrorIs(t, err, &metron.MapError{Field: "Modified"})

	// The cause stays reachable.
	require.ErrorIs(t, err, errUnusable)

	// An APIError is not a MapError, and vice versa.
	require.NotErrorIs(t, err, &metron.APIError{StatusCode: http.StatusNotFound})
	require.NotErrorIs(t, &metron.APIError{StatusCode: http.StatusNotFound}, &metron.MapError{})
}

func TestListContinuesPastAMapError(t *testing.T) {
	t.Parallel()

	// The middle record has no id, so it cannot be mapped. The two either side
	// must still arrive.
	body := `{"count":3,"next":null,"previous":null,"results":[
		{"id":1,"series":{"id":9,"name":"S","volume":1,"year_began":2000},"number":"1","cover_date":"2022-09-01","modified":"2024-12-23T16:17:29.333263-05:00"},
		{"series":{"id":9,"name":"S","volume":1,"year_began":2000},"number":"2","cover_date":"2022-09-01","modified":"2024-12-23T16:17:29.333263-05:00"},
		{"id":3,"series":{"id":9,"name":"S","volume":1,"year_began":2000},"number":"3","cover_date":"2022-09-01","modified":"2024-12-23T16:17:29.333263-05:00"}
	]}`

	c := newTestClient(t, []requestMock{
		{expectedURL: "https://metron.cloud/api/issue/?page=1", responseBody: body},
	})

	var ids []int
	var errs []error

	for issue, err := range c.Issues(context.Background()) {
		if err != nil {
			errs = append(errs, err)

			continue
		}

		ids = append(ids, issue.ID)
	}

	require.Equal(t, []int{1, 3}, ids)
	require.Len(t, errs, 1)
	require.ErrorIs(t, errs[0], &metron.MapError{Kind: "issue", Field: "Id"})
}
