package ratelimit_test

import (
	"errors"
	"net/http"
	"strconv"
	"testing"
	"time"

	"codeberg.org/jyggen/go-httpkit"
	"github.com/jyggen/go-metron/internal/ratelimit"
	"github.com/stretchr/testify/require"
)

func withHeaders(headers map[string]string) *http.Response {
	res := &http.Response{Header: make(http.Header)}

	for k, v := range headers {
		res.Header.Set(k, v)
	}

	return res
}

func TestMiddlewareNoPriorState(t *testing.T) {
	t.Parallel()

	calls := 0
	next := func(*http.Request) (*http.Response, error) {
		calls++

		return withHeaders(nil), nil
	}

	fn := ratelimit.Middleware()(next)

	_, err := fn(&http.Request{})

	require.NoError(t, err)
	require.Equal(t, 1, calls)
}

func TestMiddlewareBlocksOnExhaustedWindow(t *testing.T) {
	t.Parallel()

	reset := time.Now().Add(time.Hour)

	calls := 0
	next := func(*http.Request) (*http.Response, error) {
		calls++

		return withHeaders(map[string]string{
			"X-RateLimit-Burst-Remaining": "0",
			"X-RateLimit-Burst-Reset":     strconv.FormatInt(reset.Unix(), 10),
		}), nil
	}

	fn := ratelimit.Middleware()(next)

	_, err := fn(&http.Request{}) // observes the exhausted window
	require.NoError(t, err)

	_, err = fn(&http.Request{}) // should be pre-empted

	var retryAfter *httpkit.RetryAfterError

	require.ErrorAs(t, err, &retryAfter)
	require.InDelta(t, time.Hour, retryAfter.RetryAfter(), float64(time.Second))
	require.Equal(t, 1, calls, "next must not be called once the window is exhausted")
}

func TestMiddlewareAllowsRemainingWindow(t *testing.T) {
	t.Parallel()

	calls := 0
	next := func(*http.Request) (*http.Response, error) {
		calls++

		return withHeaders(map[string]string{
			"X-RateLimit-Burst-Remaining": "5",
			"X-RateLimit-Burst-Reset":     strconv.FormatInt(time.Now().Add(time.Hour).Unix(), 10),
		}), nil
	}

	fn := ratelimit.Middleware()(next)

	_, err := fn(&http.Request{})
	require.NoError(t, err)

	_, err = fn(&http.Request{})
	require.NoError(t, err)
	require.Equal(t, 2, calls)
}

func TestMiddlewareMostRestrictiveWindowWins(t *testing.T) {
	t.Parallel()

	burstReset := time.Now().Add(time.Minute)
	sustainedReset := time.Now().Add(time.Hour)

	calls := 0
	next := func(*http.Request) (*http.Response, error) {
		calls++

		return withHeaders(map[string]string{
			"X-RateLimit-Burst-Remaining":     "0",
			"X-RateLimit-Burst-Reset":         strconv.FormatInt(burstReset.Unix(), 10),
			"X-RateLimit-Sustained-Remaining": "0",
			"X-RateLimit-Sustained-Reset":     strconv.FormatInt(sustainedReset.Unix(), 10),
		}), nil
	}

	fn := ratelimit.Middleware()(next)

	_, err := fn(&http.Request{})
	require.NoError(t, err)

	_, err = fn(&http.Request{})

	var retryAfter *httpkit.RetryAfterError

	require.ErrorAs(t, err, &retryAfter)
	require.InDelta(t, time.Hour, retryAfter.RetryAfter(), float64(time.Second))
	require.Equal(t, 1, calls)
}

func TestMiddlewareIgnoresAbsentHeaders(t *testing.T) {
	t.Parallel()

	reset := time.Now().Add(time.Hour)

	calls := 0
	next := func(*http.Request) (*http.Response, error) {
		calls++

		if calls == 1 {
			return withHeaders(map[string]string{
				"X-RateLimit-Burst-Remaining": "0",
				"X-RateLimit-Burst-Reset":     strconv.FormatInt(reset.Unix(), 10),
			}), nil
		}

		// A later response with no rate-limit headers must not erase the
		// previously observed exhausted state.
		return withHeaders(nil), nil
	}

	fn := ratelimit.Middleware()(next)

	_, err := fn(&http.Request{})
	require.NoError(t, err)

	_, err = fn(&http.Request{})

	var retryAfter *httpkit.RetryAfterError

	require.ErrorAs(t, err, &retryAfter)
	require.Equal(t, 1, calls)
}

func TestMiddlewarePassesThroughNextError(t *testing.T) {
	t.Parallel()

	wantErr := errors.New("boom")
	next := func(*http.Request) (*http.Response, error) {
		return nil, wantErr
	}

	fn := ratelimit.Middleware()(next)

	_, err := fn(&http.Request{})

	require.ErrorIs(t, err, wantErr)
}
