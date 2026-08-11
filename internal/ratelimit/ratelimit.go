// Package ratelimit implements a reactive, header-driven rate-limit
// middleware for Metron's X-RateLimit-* response headers.
package ratelimit

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

// window is the last observed state of one rate-limit window (burst or
// sustained), from Metron's X-RateLimit-* headers. Fields are nil until first
// observed.
type window struct {
	remaining *int
	reset     *time.Time
}

// Middleware tracks rate-limit state reactively. The sustained (daily) limit
// varies per account tier, so rather than enforce a local quota it trusts the
// last observed headers and pre-empts only an exhausted window.
func Middleware() httpkit.Middleware {
	var m sync.RWMutex

	var burst, sustained window

	return func(next httpkit.MiddlewareFunc) httpkit.MiddlewareFunc {
		return func(r *http.Request) (*http.Response, error) {
			m.RLock()
			d := max(wait(burst), wait(sustained))
			m.RUnlock()

			if d > 0 {
				return nil, httpkit.NewRetryAfterError(d)
			}

			res, err := next(r)
			if err != nil {
				return res, err
			}

			newBurst, okBurst := parseWindow(res.Header, headerBurstRemaining, headerBurstReset)
			newSustained, okSustained := parseWindow(res.Header, headerSustainedRemaining, headerSustainedReset)

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

// wait returns how long until w allows another request, or 0 if it already
// does or has not been observed.
func wait(w window) time.Duration {
	if w.remaining == nil || *w.remaining > 0 || w.reset == nil {
		return 0
	}

	return max(0, time.Until(*w.reset))
}

// parseWindow parses one window's Remaining/Reset headers, reporting false
// when neither is present so previously observed state survives.
func parseWindow(header http.Header, remainingHeader, resetHeader string) (window, bool) {
	remainingStr := header.Get(remainingHeader)
	resetStr := header.Get(resetHeader)

	if remainingStr == "" && resetStr == "" {
		return window{}, false
	}

	var w window

	if v, err := strconv.Atoi(remainingStr); err == nil {
		w.remaining = &v
	}

	if v, err := strconv.ParseInt(resetStr, 10, 64); err == nil {
		w.reset = new(time.Unix(v, 0))
	}

	return w, true
}
