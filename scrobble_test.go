package metron_test

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/jyggen/go-metron"
	"github.com/stretchr/testify/require"
)

type scrobbleTestCase struct {
	name         string
	fixture      string
	issueID      int
	opts         []metron.ScrobbleOption
	validateBody func(t *testing.T, body []byte)
	expected     *metron.ScrobbleResult
}

func TestScrobble(t *testing.T) {
	t.Parallel()

	dateRead := time.Date(2024, 3, 15, 10, 0, 0, 0, time.UTC)

	testCases := []scrobbleTestCase{
		{
			name:    "basic scrobble",
			fixture: "fixtures/scrobble_1.json",
			issueID: 12345,
			validateBody: func(t *testing.T, body []byte) {
				t.Helper()

				var req struct {
					IssueID  int  `json:"issue_id"`
					DateRead *any `json:"date_read"`
					Rating   *any `json:"rating"`
				}

				require.NoError(t, json.Unmarshal(body, &req))
				require.Equal(t, 12345, req.IssueID)
				require.Nil(t, req.DateRead)
				require.Nil(t, req.Rating)
			},
			expected: &metron.ScrobbleResult{
				ID: 42,
				Issue: metron.ScrobbleIssue{
					ID:        12345,
					Number:    "1",
					CoverDate: parseDate(t, "2023-06-15"),
					StoreDate: new(parseDate(t, "2023-06-13")),
					Series: metron.ScrobbleIssueSeries{
						Name:      "Amazing Spider-Man",
						Volume:    6,
						YearBegan: 2022,
					},
					Modified: parseTime(t, "2024-01-15T10:30:00.123456-05:00"),
				},
				IsRead:   true,
				ReadDate: new(parseTime(t, "2024-03-01T12:00:00Z")),
				Rating:   new(4),
				Created:  true,
				Modified: parseTime(t, "2024-03-01T12:00:00.654321-05:00"),
			},
		},
		{
			name:    "with options",
			fixture: "fixtures/scrobble_2.json",
			issueID: 12345,
			opts:    []metron.ScrobbleOption{metron.WithReadDate(dateRead), metron.WithRating(4)},
			validateBody: func(t *testing.T, body []byte) {
				t.Helper()

				var req struct {
					IssueID  int    `json:"issue_id"`
					DateRead string `json:"date_read"`
					Rating   int    `json:"rating"`
				}

				require.NoError(t, json.Unmarshal(body, &req))
				require.Equal(t, 12345, req.IssueID)
				require.Equal(t, 4, req.Rating)
				require.NotEmpty(t, req.DateRead)
			},
			expected: &metron.ScrobbleResult{
				ID: 42,
				Issue: metron.ScrobbleIssue{
					ID:        12345,
					Number:    "1",
					CoverDate: parseDate(t, "2023-06-15"),
					StoreDate: new(parseDate(t, "2023-06-13")),
					Series: metron.ScrobbleIssueSeries{
						Name:      "Amazing Spider-Man",
						Volume:    6,
						YearBegan: 2022,
					},
					Modified: parseTime(t, "2024-01-15T10:30:00.123456-05:00"),
				},
				IsRead:   true,
				ReadDate: new(parseTime(t, "2024-03-01T12:00:00Z")),
				Rating:   new(4),
				Created:  false,
				Modified: parseTime(t, "2024-03-01T12:00:00.654321-05:00"),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			c := newTestClient(t, []requestMock{
				{
					expectedURL:         "https://metron.cloud/api/collection/scrobble/",
					expectedMethod:      http.MethodPost,
					responseBodyFixture: tc.fixture,
					validateBody:        tc.validateBody,
				},
			})

			result, err := c.Scrobble(context.Background(), tc.issueID, tc.opts...)

			require.NoError(t, err)
			require.Equal(t, tc.expected, result)
		})
	}
}

// scrobbleResponse renders a minimal scrobble response with the given raw JSON
// as the rating, so each form of the RatingEnum|NullEnum union can be checked.
func scrobbleResponse(rating string) string {
	return `{"id":42,"created":true,"is_read":true,` +
		`"modified":"2024-01-15T10:30:00.123456-05:00","issue":{"id":12345,"number":"1",` +
		`"cover_date":"2023-06-15","modified":"2024-01-15T10:30:00.123456-05:00",` +
		`"series":{"name":"Amazing Spider-Man","volume":6,"year_began":2022}},` +
		`"rating":` + rating + `}`
}

// TestScrobbleRating pins down how each rating form round-trips. A null or
// absent rating is legitimate and yields nil; a malformed one must error
// rather than silently arriving as nil.
func TestScrobbleRating(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name           string
		rating         string
		expectedRating *int
		expectedError  string
	}{
		{name: "in range", rating: `4`, expectedRating: new(4)},
		{name: "null", rating: `null`, expectedRating: nil},
		{name: "string", rating: `"4"`, expectedError: "scrobble: rating:"},
		{name: "object", rating: `{"value":4}`, expectedError: "scrobble: rating:"},
		{name: "boolean", rating: `true`, expectedError: "scrobble: rating:"},
		// Metron documents 1-5, but out-of-range numbers are passed through
		// rather than rejected — the API is trusted on its own enum.
		{name: "below range", rating: `0`, expectedRating: new(0)},
		{name: "above range", rating: `9`, expectedRating: new(9)},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			c := newTestClient(t, []requestMock{
				{
					expectedURL:    "https://metron.cloud/api/collection/scrobble/",
					expectedMethod: http.MethodPost,
					responseStatus: http.StatusCreated,
					responseBody:   scrobbleResponse(tc.rating),
				},
			})

			result, err := c.Scrobble(context.Background(), 12345)

			if tc.expectedError != "" {
				require.ErrorContains(t, err, tc.expectedError)
				require.Nil(t, result)

				return
			}

			require.NoError(t, err)
			require.Equal(t, tc.expectedRating, result.Rating)
		})
	}
}
