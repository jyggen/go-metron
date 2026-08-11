package backoff_test

import (
	"errors"
	"io"
	"net/http"
	"testing"
	"time"

	"codeberg.org/jyggen/go-httpkit"
	"github.com/jyggen/go-metron/internal/backoff"
	"github.com/stretchr/testify/require"
)

func TestMiddlewareNoPriorDeadline(t *testing.T) {
	t.Parallel()

	calls := 0
	next := func(*http.Request) (*http.Response, error) {
		calls++

		return &http.Response{StatusCode: http.StatusOK, Body: http.NoBody}, nil
	}

	fn := backoff.Middleware()(next)

	_, err := fn(&http.Request{})

	require.NoError(t, err)
	require.Equal(t, 1, calls)
}

func TestMiddlewareBlocksOnUpstreamRetryAfterError(t *testing.T) {
	t.Parallel()

	calls := 0
	next := func(*http.Request) (*http.Response, error) {
		calls++

		return nil, httpkit.NewRetryAfterError(time.Hour)
	}

	fn := backoff.Middleware()(next)

	_, err := fn(&http.Request{}) // observes the nested RetryAfterError

	var retryAfter *httpkit.RetryAfterError

	require.ErrorAs(t, err, &retryAfter)

	_, err = fn(&http.Request{}) // pre-empted before next is called again

	require.ErrorAs(t, err, &retryAfter)
	require.InDelta(t, time.Hour, retryAfter.RetryAfter(), float64(time.Second))
	require.Equal(t, 1, calls)
}

func TestMiddlewareBlocksOn429DeltaSeconds(t *testing.T) {
	t.Parallel()

	calls := 0
	next := func(*http.Request) (*http.Response, error) {
		calls++

		header := make(http.Header)
		header.Set("Retry-After", "30")

		return &http.Response{StatusCode: http.StatusTooManyRequests, Header: header, Body: http.NoBody}, nil
	}

	fn := backoff.Middleware()(next)

	_, err := fn(&http.Request{})

	var retryAfter *httpkit.RetryAfterError

	require.ErrorAs(t, err, &retryAfter)
	require.Equal(t, 30*time.Second, retryAfter.RetryAfter())

	_, err = fn(&http.Request{})

	require.ErrorAs(t, err, &retryAfter)
	require.InDelta(t, 30*time.Second, retryAfter.RetryAfter(), float64(time.Second))
	require.Equal(t, 1, calls)
}

func TestMiddlewareBlocksOn429HTTPDate(t *testing.T) {
	t.Parallel()

	header := make(http.Header)
	header.Set("Retry-After", time.Now().Add(time.Hour).UTC().Format(http.TimeFormat))

	next := func(*http.Request) (*http.Response, error) {
		return &http.Response{StatusCode: http.StatusTooManyRequests, Header: header, Body: http.NoBody}, nil
	}

	fn := backoff.Middleware()(next)

	_, err := fn(&http.Request{})

	var retryAfter *httpkit.RetryAfterError

	require.ErrorAs(t, err, &retryAfter)
	require.InDelta(t, time.Hour, retryAfter.RetryAfter(), float64(time.Second))
}

func TestMiddlewareIgnores429WithoutUsableRetryAfter(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name       string
		retryAfter string
	}{
		{name: "absent", retryAfter: ""},
		{name: "malformed", retryAfter: "later please"},
		{name: "negative", retryAfter: "-5"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			calls := 0
			next := func(*http.Request) (*http.Response, error) {
				calls++

				header := make(http.Header)
				if tc.retryAfter != "" {
					header.Set("Retry-After", tc.retryAfter)
				}

				return &http.Response{StatusCode: http.StatusTooManyRequests, Header: header, Body: http.NoBody}, nil
			}

			fn := backoff.Middleware()(next)

			res, err := fn(&http.Request{})

			require.NoError(t, err)
			require.Equal(t, http.StatusTooManyRequests, res.StatusCode)

			// Not treated as a deadline: a second call still reaches next.
			_, err = fn(&http.Request{})

			require.NoError(t, err)
			require.Equal(t, 2, calls)
		})
	}
}

func TestMiddlewarePassesThroughOtherNextError(t *testing.T) {
	t.Parallel()

	wantErr := errors.New("boom")
	next := func(*http.Request) (*http.Response, error) {
		return nil, wantErr
	}

	fn := backoff.Middleware()(next)

	_, err := fn(&http.Request{})

	require.ErrorIs(t, err, wantErr)
}

func TestMiddlewarePassesThroughSuccessfulResponse(t *testing.T) {
	t.Parallel()

	next := func(*http.Request) (*http.Response, error) {
		return &http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(nil)}, nil
	}

	fn := backoff.Middleware()(next)

	res, err := fn(&http.Request{})

	require.NoError(t, err)
	require.Equal(t, http.StatusOK, res.StatusCode)
}
