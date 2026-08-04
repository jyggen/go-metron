package metron

import "fmt"

// maxErrorBodyBytes caps how much of an error body is retained on an APIError.
const maxErrorBodyBytes = 4 << 10

// APIError is returned when the API responds with an unexpected status code.
// Body is the response body, truncated to 4KiB.
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

// Is matches any APIError with the same status code:
//
//	errors.Is(err, &metron.APIError{StatusCode: http.StatusNotFound})
func (e *APIError) Is(target error) bool {
	t, ok := target.(*APIError)

	return ok && t.StatusCode == e.StatusCode
}

// FilterError is returned when a Filter is passed to a list method whose
// endpoint does not support it. It surfaces on the first iteration.
type FilterError struct {
	// Filter is the name of the offending filter, e.g. "ByPublisherID".
	Filter string
}

// Error implements the error interface.
func (e *FilterError) Error() string {
	return fmt.Sprintf("metron: %s does not apply to this endpoint", e.Filter)
}

// Is matches any FilterError for the same filter:
//
//	errors.Is(err, &metron.FilterError{Filter: "ByPublisherID"})
func (e *FilterError) Is(target error) bool {
	t, ok := target.(*FilterError)

	return ok && t.Filter == e.Filter
}
