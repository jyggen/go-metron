package metron

import (
	"context"
	"net/http"
	"time"

	"github.com/jyggen/go-metron/internal"
)

// requestOptions is what a RequestOption mutates, mirroring Option/clientOptions
// so a per-request setting never becomes a field on Client.
type requestOptions struct {
	ifModifiedSince *time.Time
}

// RequestOption configures a single detail request. Like ScrobbleOption it is a
// closed extension point: callers cannot write their own.
type RequestOption func(*requestOptions)

// IfModifiedSince makes the request conditional on the record having changed
// since t. The API answers an unchanged record with 304, which the client
// reports as ErrNotModified.
//
// A record's own Modified field is the timestamp to send back, so a caller
// holding a previous result needs to store nothing extra:
//
//	issue, err := c.IssueByID(ctx, id, metron.IfModifiedSince(cached.Modified))
//	if errors.Is(err, metron.ErrNotModified) {
//		// cached is still current
//	}
//
// List methods take ByModifiedGreaterThan instead, which filters server-side
// rather than answering all-or-nothing per page.
func IfModifiedSince(t time.Time) RequestOption {
	return func(o *requestOptions) {
		o.ifModifiedSince = &t
	}
}

// requestEditors folds opts and translates them into the request edits they
// describe. It returns nil for an unconditional request, since the generated
// client calls every editor it is handed and a no-op one would only cost.
func requestEditors(opts []RequestOption) []internal.RequestEditorFn {
	var o requestOptions

	for _, opt := range opts {
		opt(&o)
	}

	if o.ifModifiedSince == nil {
		return nil
	}

	// http.TimeFormat hardcodes GMT, so it renders any other zone as a lie.
	value := o.ifModifiedSince.UTC().Format(http.TimeFormat)

	return []internal.RequestEditorFn{
		func(_ context.Context, req *http.Request) error {
			req.Header.Set("If-Modified-Since", value)

			return nil
		},
	}
}
