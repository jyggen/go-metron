package metron_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/jyggen/go-metron"
	"github.com/stretchr/testify/require"
)

func TestIfModifiedSince(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name           string
		opts           []metron.RequestOption
		expectedHeader string
	}{
		{
			name:           "absent without the option",
			opts:           nil,
			expectedHeader: "",
		},
		{
			name:           "UTC",
			opts:           []metron.RequestOption{metron.IfModifiedSince(time.Date(2024, time.March, 7, 17, 36, 10, 0, time.UTC))},
			expectedHeader: "Thu, 07 Mar 2024 17:36:10 GMT",
		},
		{
			// http.TimeFormat hardcodes GMT, so a non-UTC time must be converted
			// rather than formatted where it stands.
			name:           "converted from another zone",
			opts:           []metron.RequestOption{metron.IfModifiedSince(parseTime(t, "2024-03-07T12:36:10.633143-05:00"))},
			expectedHeader: "Thu, 07 Mar 2024 17:36:10 GMT",
		},
		{
			name:           "last option wins",
			opts:           []metron.RequestOption{metron.IfModifiedSince(time.Unix(0, 0)), metron.IfModifiedSince(time.Date(2024, time.March, 7, 17, 36, 10, 0, time.UTC))},
			expectedHeader: "Thu, 07 Mar 2024 17:36:10 GMT",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			c := newTestClient(t, []requestMock{
				{
					expectedURL:         "https://metron.cloud/api/arc/659/",
					responseBodyFixture: "fixtures/arc_659.json",
					expectedHeaders:     map[string]string{"If-Modified-Since": tc.expectedHeader},
				},
			})

			arc, err := c.ArcByID(context.Background(), 659, tc.opts...)

			require.NoError(t, err)
			require.Equal(t, 659, arc.ID)
		})
	}
}

func TestNotModified(t *testing.T) {
	t.Parallel()

	c := newTestClient(t, []requestMock{
		{
			expectedURL:     "https://metron.cloud/api/arc/659/",
			responseStatus:  http.StatusNotModified,
			expectedHeaders: map[string]string{"If-Modified-Since": "Thu, 07 Mar 2024 17:36:10 GMT"},
		},
	})

	arc, err := c.ArcByID(context.Background(), 659,
		metron.IfModifiedSince(time.Date(2024, time.March, 7, 17, 36, 10, 0, time.UTC)),
	)

	require.Nil(t, arc)
	require.ErrorIs(t, err, metron.ErrNotModified)

	// A 304 is the answer the request asked for, so it never becomes an APIError.
	var apiErr *metron.APIError

	require.NotErrorAs(t, err, &apiErr)
	require.NotErrorIs(t, err, &metron.APIError{StatusCode: http.StatusNotModified})
}

func TestNotModifiedNotRetried(t *testing.T) {
	t.Parallel()

	// ErrNotModified is not an APIError, so without an explicit exclusion the
	// retry loop would take it for a network failure and ask again.
	c, count := retryClient(t, 3, []*http.Response{status(http.StatusNotModified)})

	_, err := c.ArcByID(context.Background(), 659, metron.IfModifiedSince(time.Unix(0, 0)))

	require.ErrorIs(t, err, metron.ErrNotModified)
	require.Equal(t, int64(1), count.Load())
}

func TestNotModifiedWithoutRequestingIt(t *testing.T) {
	t.Parallel()

	// A proxy or a caller's own RoundTripper can answer 304 unprompted.
	c := newTestClient(t, []requestMock{
		{expectedURL: "https://metron.cloud/api/arc/659/", responseStatus: http.StatusNotModified},
	})

	_, err := c.ArcByID(context.Background(), 659)

	require.ErrorIs(t, err, metron.ErrNotModified)
}
