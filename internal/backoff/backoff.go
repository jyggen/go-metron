// Package backoff implements a global backoff middleware for 429 responses,
// honouring the Retry-After header.
package backoff

import (
	"errors"
	"net/http"
	"strconv"
	"sync"
	"time"

	"codeberg.org/jyggen/go-httpkit"
)

// Middleware tracks a single deadline shared by every request: a 429 (or an
// upstream RetryAfterError) pushes it out, and a request made before it
// elapses is pre-empted rather than sent.
func Middleware() httpkit.Middleware {
	var m sync.RWMutex

	deadline := time.Now()

	return func(next httpkit.MiddlewareFunc) httpkit.MiddlewareFunc {
		return func(r *http.Request) (*http.Response, error) {
			m.RLock()
			wait := time.Until(deadline)
			m.RUnlock()

			if wait > 0 {
				return nil, httpkit.NewRetryAfterError(wait)
			}

			res, err := next(r)

			var retryAfterErr *httpkit.RetryAfterError
			if errors.As(err, &retryAfterErr) {
				m.Lock()
				if t := time.Now().Add(retryAfterErr.RetryAfter()); t.After(deadline) {
					deadline = t
				}
				m.Unlock()

				return nil, retryAfterErr
			}

			if err != nil {
				return nil, err
			}

			if res.StatusCode == http.StatusTooManyRequests {
				now := time.Now()

				if d, ok := parseRetryAfter(res.Header.Get("Retry-After"), now); ok {
					m.Lock()
					if t := now.Add(d); t.After(deadline) {
						deadline = t
					}
					m.Unlock()

					return nil, errors.Join(httpkit.NewRetryAfterError(d), res.Body.Close())
				}
			}

			return res, nil
		}
	}
}

// parseRetryAfter parses either RFC 9110 form of Retry-After, resolving the
// date form against now. Reports false if absent, malformed or negative.
func parseRetryAfter(value string, now time.Time) (time.Duration, bool) {
	if value == "" {
		return 0, false
	}

	if seconds, err := strconv.Atoi(value); err == nil {
		if seconds < 0 {
			return 0, false
		}

		return time.Duration(seconds) * time.Second, true
	}

	if t, err := http.ParseTime(value); err == nil {
		return max(0, t.Sub(now)), true
	}

	return 0, false
}
