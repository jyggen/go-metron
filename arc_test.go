package metron_test

import (
	"context"
	"net/http"
	"sync/atomic"
	"testing"

	"github.com/jyggen/go-metron"
	"github.com/stretchr/testify/require"
)

func TestArcByID(t *testing.T) {
	t.Parallel()
	testByID(t, "arc", (*metron.Client).ArcByID, []testCase[*metron.Arc]{
		{
			id: 659,
			expected: &metron.Arc{
				ID:          659,
				Name:        "'Til Death Do Us...",
				Description: asReference(""),
				ImageURL:    nil,
				ComicVineID: nil,
				ResourceURL: parseURL(t, "https://metron.cloud/arc/til-death-do-us/"),
				Modified:    parseTime(t, "2022-01-16T10:12:00.426525-05:00"),
			},
		},
		{
			id: 1493,
			expected: &metron.Arc{
				ID:          1493,
				Name:        "Crisis on Infinite Darkwings",
				Description: asReference(""),
				ImageURL: asReference(
					parseURL(
						t,
						"https://static.metron.cloud/media/arc/2024/03/07/3069ed71448e4cf3bdca8f808fc596ce.jpg",
					),
				),
				ComicVineID: asReference(56259),
				ResourceURL: parseURL(t, "https://metron.cloud/arc/darkwing-duck-crisis-on-infinite-darkwings/"),
				Modified:    parseTime(t, "2024-03-07T12:36:10.633143-05:00"),
			},
		},
	})
}

func TestArcsCached(t *testing.T) {
	t.Parallel()

	var requestCount atomic.Int64

	mocks := []requestMock{
		{"https://metron.cloud/api/arc/?page=1", "fixtures/arc_list_1.json"},
		{"https://metron.cloud/api/arc/?page=2", "fixtures/arc_list_2.json"},
	}

	c, err := metron.NewClient("username", "password",
		metron.WithCaching(),
		metron.WithStoragePath(t.TempDir()),
		metron.WithClient(&http.Client{
			Transport: roundTripFunc(func(req *http.Request) *http.Response {
				idx := int(requestCount.Add(1)) - 1
				require.Less(t, idx, len(mocks), "unexpected extra HTTP request: %s", req.URL)
				m := mocks[idx]

				require.Equal(t, m.expectedURL, req.URL.String())

				f, openErr := fs.Open(m.responseBodyFixture)
				require.NoError(t, openErr)

				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       f,
					Header: http.Header{
						"Last-Modified": {"Mon, 01 Jan 2024 00:00:00 GMT"},
					},
				}
			}),
		}),
	)
	require.NoError(t, err)
	defer c.Close()

	// First iteration — hits the network.
	var firstResults []*metron.ArcList
	for res, iterErr := range c.Arcs(context.Background()) {
		require.NoError(t, iterErr)
		firstResults = append(firstResults, res)
	}
	require.Len(t, firstResults, 4)
	require.Equal(t, int64(2), requestCount.Load(), "expected 2 HTTP requests on first iteration")

	// Second iteration — should be served from cache (no additional HTTP requests).
	var secondResults []*metron.ArcList
	for res, iterErr := range c.Arcs(context.Background()) {
		require.NoError(t, iterErr)
		secondResults = append(secondResults, res)
	}
	require.Len(t, secondResults, 4)
	require.Equal(t, int64(2), requestCount.Load(), "expected no additional HTTP requests on cached iteration")
	require.Equal(t, firstResults, secondResults)
}

func TestArcs(t *testing.T) {
	t.Parallel()
	testList(t, "arc", (*metron.Client).Arcs, []testCase[*metron.ArcList]{
		{
			id: 659,
			expected: &metron.ArcList{
				ID:       659,
				Name:     "'Til Death Do Us...",
				Modified: parseTime(t, "2022-01-16T10:12:00.426525-05:00"),
			},
		},
		{
			id: 931,
			expected: &metron.ArcList{
				ID:       931,
				Name:     "(She) Drunk History",
				Modified: parseTime(t, "2023-02-15T11:47:59.483664-05:00"),
			},
		},
		{
			id: 871,
			expected: &metron.ArcList{
				ID:       871,
				Name:     "1+2 = Fantastic Three",
				Modified: parseTime(t, "2022-09-08T09:53:30.626809-04:00"),
			},
		},
		{
			id: 1031,
			expected: &metron.ArcList{
				ID:       1031,
				Name:     "1602",
				Modified: parseTime(t, "2023-03-12T15:41:46.220978-04:00"),
			},
		},
	})
}
