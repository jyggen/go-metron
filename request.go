package metron

import (
	"context"
	"net/http"
	"time"

	"github.com/jyggen/go-metron/internal"
)

type requestOptions struct {
	ifModifiedSince *time.Time
}

// RequestOption configures a request to any read endpoint. None is currently
// defined; the parameter reserves the position on every read method.
type RequestOption interface {
	applyRequest(*requestOptions)
}

// ConditionalOption configures a request to an endpoint that honours
// If-Modified-Since: the detail methods and the per-parent issue listings.
// Other endpoints ignore the header and take RequestOption instead.
type ConditionalOption interface {
	applyConditional(*requestOptions)
}

type ifModifiedSince time.Time

func (o ifModifiedSince) applyConditional(r *requestOptions) {
	r.ifModifiedSince = new(time.Time(o))
}

// IfModifiedSince makes the request conditional on the record having changed
// since t. The API answers an unchanged record with 304, which the client
// reports as ErrNotModified.
//
// A record's own Modified field is the timestamp to send back:
//
//	issue, err := c.IssueByID(ctx, id, metron.IfModifiedSince(cached.Modified))
//	if errors.Is(err, metron.ErrNotModified) {
//		// cached is still current
//	}
//
// IssuesByArcID, IssuesByCharacterID and IssuesByTeamID compare against the
// parent's timestamp, so a 304 there means no issue joined or left the parent —
// an issue's own fields may still have changed. IssuesBySeriesID is exact.
func IfModifiedSince(t time.Time) ConditionalOption {
	return ifModifiedSince(t)
}

// editors returns nil when no option was set, since the generated client calls
// every editor it is handed.
func (o requestOptions) editors() []internal.RequestEditorFn {
	if o.ifModifiedSince == nil {
		return nil
	}

	// http.TimeFormat hardcodes GMT, so any other zone would render as a lie.
	value := o.ifModifiedSince.UTC().Format(http.TimeFormat)

	return []internal.RequestEditorFn{
		func(_ context.Context, req *http.Request) error {
			req.Header.Set("If-Modified-Since", value)

			return nil
		},
	}
}

func requestEditors(opts []RequestOption) []internal.RequestEditorFn {
	var o requestOptions

	for _, opt := range opts {
		opt.applyRequest(&o)
	}

	return o.editors()
}

func conditionalEditors(opts []ConditionalOption) []internal.RequestEditorFn {
	var o requestOptions

	for _, opt := range opts {
		opt.applyConditional(&o)
	}

	return o.editors()
}
