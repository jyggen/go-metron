package metron

import (
	"context"
	"errors"
	"math/rand/v2"
	"net/http"
	"strconv"
	"time"
)

// Bounds for the exponential backoff on transient failures.
const (
	backOffBase = 250 * time.Millisecond
	backOffCap  = 8 * time.Second
)

// retryableStatus reports whether a status is worth retrying. 500 is excluded:
// Metron is Django, where it means an unhandled exception, usually deterministic.
func retryableStatus(status int) bool {
	switch status {
	case http.StatusBadGateway, http.StatusServiceUnavailable, http.StatusGatewayTimeout:
		return true
	default:
		return false
	}
}

// retryable reports whether err is worth another attempt. Transient failures
// apply only to idempotent requests.
func retryable(err error, idempotent bool) bool {
	if !idempotent || errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
		return false
	}

	var apiErr *APIError
	if errors.As(err, &apiErr) {
		return retryableStatus(apiErr.StatusCode)
	}

	// Anything else is a network failure; the response never arrived.
	return true
}

// backOff returns the wait before attempt, jittered so retries do not synchronise.
func backOff(attempt int) time.Duration {
	d := min(backOffBase<<attempt, backOffCap)

	return time.Duration(rand.Int64N(int64(d)) + 1)
}

// parseRetryAfter parses either RFC 9110 form of Retry-After, resolving the date
// form against now. Reports false if absent, malformed or negative.
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
