package metron

import (
	"fmt"
	"net/http"
	"strconv"
	"time"
)

// maxErrorBodyBytes caps how much of an error response body is retained on an
// APIError, so a large or malformed body can't be held in memory indefinitely.
const maxErrorBodyBytes = 4 << 10

// APIError is returned when the Metron API responds with an unexpected status
// code. Body holds the response body as returned by the API, truncated to 4KiB.
type APIError struct {
	StatusCode int
	Body       []byte
}

// Error implements the error interface.
func (e *APIError) Error() string {
	if len(e.Body) == 0 {
		return fmt.Sprintf("metron: unexpected status code: %d", e.StatusCode)
	}

	return fmt.Sprintf("metron: unexpected status code: %d: %s", e.StatusCode, e.Body)
}

// Is reports whether target is an APIError with the same status code, so
// callers can match on status alone:
//
//	errors.Is(err, &metron.APIError{StatusCode: http.StatusNotFound})
func (e *APIError) Is(target error) bool {
	t, ok := target.(*APIError)

	return ok && t.StatusCode == e.StatusCode
}

// parseRetryAfter parses a Retry-After header value. RFC 9110 allows both a
// delta-seconds and an HTTP-date form, so both are accepted; the date form is
// resolved against now. It reports false if the value is absent, malformed or
// negative, leaving the caller to treat the response as a plain error.
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
