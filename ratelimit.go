package metron

import (
	"net/http"
	"strconv"
	"sync"
	"time"

	"codeberg.org/jyggen/go-httpkit"
)

const (
	headerBurstRemaining     = "X-RateLimit-Burst-Remaining"
	headerBurstReset         = "X-RateLimit-Burst-Reset"
	headerSustainedRemaining = "X-RateLimit-Sustained-Remaining"
	headerSustainedReset     = "X-RateLimit-Sustained-Reset"
)

// rateLimitWindow holds the most recently observed state of a single
// rate-limit window (burst or sustained), as reported by Metron's
// X-RateLimit-* response headers. Fields are nil until observed.
type rateLimitWindow struct {
	remaining *int
	reset     *time.Time
}

// newRateLimitMiddleware tracks Metron's rate-limit state reactively.
// Metron's sustained (daily) limit varies per account tier and is only known
// once observed from response headers, so instead of enforcing a fixed local
// quota, this middleware trusts the most recently observed X-RateLimit-*
// headers and pre-empts a request only once they show a window is already
// exhausted.
func newRateLimitMiddleware() httpkit.Middleware {
	var m sync.RWMutex

	var burst, sustained rateLimitWindow

	return func(next httpkit.MiddlewareFunc) httpkit.MiddlewareFunc {
		return func(r *http.Request) (*http.Response, error) {
			m.RLock()
			wait := max(rateLimitWait(burst), rateLimitWait(sustained))
			m.RUnlock()

			if wait > 0 {
				return nil, httpkit.NewRetryAfterError(wait)
			}

			res, err := next(r)
			if err != nil {
				return res, err
			}

			newBurst, okBurst := parseRateLimitWindow(res.Header, headerBurstRemaining, headerBurstReset)
			newSustained, okSustained := parseRateLimitWindow(res.Header, headerSustainedRemaining, headerSustainedReset)

			if okBurst || okSustained {
				m.Lock()
				if okBurst {
					burst = newBurst
				}
				if okSustained {
					sustained = newSustained
				}
				m.Unlock()
			}

			return res, nil
		}
	}
}

// rateLimitWait returns how long to wait before w allows another request, or
// 0 if it already does (including when w hasn't been observed yet).
func rateLimitWait(w rateLimitWindow) time.Duration {
	if w.remaining == nil || *w.remaining > 0 || w.reset == nil {
		return 0
	}

	return max(0, time.Until(*w.reset))
}

// parseRateLimitWindow parses one rate-limit window's Remaining/Reset headers.
// It reports false when neither header is present, so a response that
// doesn't report this window leaves previously observed state intact.
func parseRateLimitWindow(header http.Header, remainingHeader, resetHeader string) (rateLimitWindow, bool) {
	remainingStr := header.Get(remainingHeader)
	resetStr := header.Get(resetHeader)

	if remainingStr == "" && resetStr == "" {
		return rateLimitWindow{}, false
	}

	var w rateLimitWindow

	if v, err := strconv.Atoi(remainingStr); err == nil {
		w.remaining = &v
	}

	if v, err := strconv.ParseInt(resetStr, 10, 64); err == nil {
		t := time.Unix(v, 0)
		w.reset = &t
	}

	return w, true
}
