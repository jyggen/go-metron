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
					StoreDate: asReference(parseDate(t, "2023-06-13")),
					Series: struct {
						Name      string
						Volume    int
						YearBegan int
					}{
						Name:      "Amazing Spider-Man",
						Volume:    6,
						YearBegan: 2022,
					},
					Modified: parseTime(t, "2024-01-15T10:30:00.123456-05:00"),
				},
				IsRead:   true,
				ReadDate: asReference(parseTime(t, "2024-03-01T12:00:00Z")),
				Rating:   asReference(4),
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
					StoreDate: asReference(parseDate(t, "2023-06-13")),
					Series: struct {
						Name      string
						Volume    int
						YearBegan int
					}{
						Name:      "Amazing Spider-Man",
						Volume:    6,
						YearBegan: 2022,
					},
					Modified: parseTime(t, "2024-01-15T10:30:00.123456-05:00"),
				},
				IsRead:   true,
				ReadDate: asReference(parseTime(t, "2024-03-01T12:00:00Z")),
				Rating:   asReference(4),
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
