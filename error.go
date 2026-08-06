package metron

import (
	"fmt"
	"net/url"
	"strings"
)

// maxErrorBodyBytes caps how much of an error body is retained on an APIError.
const maxErrorBodyBytes = 4 << 10

// APIError is returned when the API responds with an unexpected status code.
// Body is the response body, truncated to 4KiB.
//
// Method and URL describe the request that failed. They name the endpoint, the
// query the filters produced and the page the iterator had reached — none of
// which a caller can reconstruct, since the client builds them.
type APIError struct {
	StatusCode int
	Body       []byte
	Method     string
	URL        *url.URL
}

// Error implements the error interface.
func (e *APIError) Error() string {
	var b strings.Builder

	b.WriteString("metron: ")

	if e.Method != "" && e.URL != nil {
		fmt.Fprintf(&b, "%s %s: ", e.Method, e.URL)
	}

	fmt.Fprintf(&b, "unexpected status code: %d", e.StatusCode)

	if len(e.Body) > 0 {
		fmt.Fprintf(&b, ": %s", e.Body)
	}

	return b.String()
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
	// Endpoint is the list method that rejected it, e.g. "Roles".
	Endpoint string
}

// Error implements the error interface.
func (e *FilterError) Error() string {
	return fmt.Sprintf("metron: %s does not apply to %s", e.Filter, e.Endpoint)
}

// Is matches any FilterError, narrowed by whichever of Filter and Endpoint the
// target sets. Zero fields match anything:
//
//	errors.Is(err, &metron.FilterError{Filter: "ByPublisherID"})
func (e *FilterError) Is(target error) bool {
	t, ok := target.(*FilterError)

	return ok &&
		(t.Filter == "" || t.Filter == e.Filter) &&
		(t.Endpoint == "" || t.Endpoint == e.Endpoint)
}

// MapError is returned when a record cannot be converted into its public type,
// either because a field the client needs is absent or because its contents are
// unusable. List methods report it per record and keep iterating, so a single
// malformed record does not end the sequence.
//
// The record is identified here because the caller cannot identify it: a list
// method yields a nil value alongside the error.
type MapError struct {
	// Kind is the resource, e.g. "issue" or "series".
	Kind string
	// ID is the record's Metron ID, or zero if the record had none.
	ID int
	// Field is the offending field, e.g. "RatingCount" or "Credits[].Creator".
	Field string
	// Err is the underlying cause, set when the field was present but unusable.
	Err error
}

// Error implements the error interface.
func (e *MapError) Error() string {
	subject := e.Kind

	if e.ID != 0 {
		subject = fmt.Sprintf("%s %d", e.Kind, e.ID)
	}

	if e.Err != nil {
		return fmt.Sprintf("metron: %s: %s: %v", subject, e.Field, e.Err)
	}

	return fmt.Sprintf("metron: %s: nil %s", subject, e.Field)
}

// Unwrap returns the underlying cause, or nil for a missing field.
func (e *MapError) Unwrap() error {
	return e.Err
}

// Is matches any MapError, narrowed by whichever of Kind, ID and Field the
// target sets. Zero fields match anything, so the common case of skipping every
// unmappable record is:
//
//	errors.Is(err, &metron.MapError{})
func (e *MapError) Is(target error) bool {
	t, ok := target.(*MapError)

	return ok &&
		(t.Kind == "" || t.Kind == e.Kind) &&
		(t.ID == 0 || t.ID == e.ID) &&
		(t.Field == "" || t.Field == e.Field)
}
