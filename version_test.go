package metron_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/jyggen/go-metron"
	"github.com/stretchr/testify/require"
)

// TestUserAgent checks the User-Agent built from build info. Deps is empty in
// our own test binary (golang/go#68045), so it falls back to "devel".
func TestUserAgent(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name     string
		options  []metron.Option
		expected string
	}{
		{
			name:     "default",
			expected: "go-metron/devel",
		},
		{
			name:     "with prefix",
			options:  []metron.Option{metron.WithUserAgent("myapp/1.2.3")},
			expected: "myapp/1.2.3 go-metron/devel",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			var got string

			options := append([]metron.Option{
				metron.WithClient(&http.Client{
					Transport: roundTripFunc(func(req *http.Request) *http.Response {
						got = req.Header.Get("User-Agent")

						f, err := fs.Open("fixtures/arc_659.json")
						require.NoError(t, err)

						return &http.Response{
							StatusCode: http.StatusOK,
							Body:       f,
							Header:     make(http.Header),
						}
					}),
				}),
			}, tc.options...)

			c, err := metron.NewClient("foobar", options...)
			require.NoError(t, err)

			_, err = c.ArcByID(context.Background(), 659)
			require.NoError(t, err)

			require.Equal(t, tc.expected, got)
		})
	}
}
