// Package metron is a Go client library for the Metron comic book database API.
//
// ImageURL and ResourceURL fields are parsed from upstream strings and are not
// validated — url.Parse accepts any scheme. Check them before following or
// rendering them.
package metron

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"iter"
	"net/http"
	"net/url"
	"sync"
	"time"

	"codeberg.org/jyggen/go-httpkit"
	"github.com/jyggen/go-metron/internal"
	"github.com/oapi-codegen/nullable"
)

const defaultBaseURL = "https://metron.cloud"

// Reference identifies a related resource by ID and display name.
//
// ID is stable: pass it to the matching ByID method. Name is an upstream
// display string that may change and is not unique, so never key off it.
type Reference struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// Client is a Metron API client. Options are consumed during construction, so
// only request-serving state is retained.
type Client struct {
	client     internal.ClientInterface
	maxRetries int
}

// clientOptions is what an Option mutates, so options can be applied in any
// order before NewClient derives the Client from them.
type clientOptions struct {
	baseURL    string
	httpClient *http.Client
	maxRetries int
	userAgent  string
}

// Option configures a Client.
type Option func(*clientOptions)

func defaultOptions() clientOptions {
	return clientOptions{
		baseURL:    defaultBaseURL,
		httpClient: &http.Client{},
		userAgent:  defaultUserAgent,
	}
}

// NewClient returns a Metron client authenticated with the given API token.
func NewClient(apiToken string, options ...Option) (*Client, error) {
	o := defaultOptions()

	for _, option := range options {
		option(&o)
	}

	if err := validateBaseURL(o.baseURL); err != nil {
		return nil, err
	}

	c := &Client{maxRetries: o.maxRetries}

	httpClient := httpkit.NewFromClient(
		o.httpClient,
		httpkit.WithBearerToken(apiToken),
		httpkit.WithUserAgent(o.userAgent),
		httpkit.WithMiddleware(newBackOffMiddleware()),
		httpkit.WithMiddleware(newRateLimitMiddleware()),
	)

	internalClient, err := internal.NewClient(o.baseURL, internal.WithHTTPClient(httpClient))
	if err != nil {
		return nil, err
	}

	c.client = internalClient

	return c, nil
}

func newBackOffMiddleware() httpkit.Middleware {
	var m sync.RWMutex

	backOff := time.Now()

	return func(next httpkit.MiddlewareFunc) httpkit.MiddlewareFunc {
		return func(r *http.Request) (*http.Response, error) {
			m.RLock()
			wait := time.Until(backOff)
			m.RUnlock()

			if wait > 0 {
				return nil, httpkit.NewRetryAfterError(wait)
			}

			res, err := next(r)

			var retryAfterErr *httpkit.RetryAfterError
			if errors.As(err, &retryAfterErr) {
				m.Lock()
				if t := time.Now().Add(retryAfterErr.RetryAfter()); t.After(backOff) {
					backOff = t
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
					if t := now.Add(d); t.After(backOff) {
						backOff = t
					}
					m.Unlock()

					return nil, errors.Join(httpkit.NewRetryAfterError(d), res.Body.Close())
				}
			}

			return res, nil
		}
	}
}

// validateBaseURL rejects anything that is not an absolute http or https URL.
func validateBaseURL(rawURL string) error {
	u, err := url.Parse(rawURL)
	if err != nil {
		return fmt.Errorf("metron: base URL: %w", err)
	}

	if u.Scheme != "http" && u.Scheme != "https" {
		return fmt.Errorf("metron: base URL: scheme must be http or https, got %q", u.Scheme)
	}

	if u.Host == "" {
		return fmt.Errorf("metron: base URL: missing host in %q", rawURL)
	}

	return nil
}

// WithBaseURL points the client at another Metron deployment: a local test
// server, a proxy or a mirror. NewClient rejects a non-http(s) or hostless URL.
func WithBaseURL(rawURL string) Option {
	return func(o *clientOptions) {
		o.baseURL = rawURL
	}
}

// WithClient sets the underlying HTTP client.
func WithClient(client *http.Client) Option {
	return func(o *clientOptions) {
		o.httpClient = client
	}
}

// WithRetry sets the maximum number of retries for rate-limited requests.
// Negative values are clamped to zero.
func WithRetry(maxRetries int) Option {
	return func(o *clientOptions) {
		o.maxRetries = max(0, maxRetries)
	}
}

// WithUserAgent prefixes the given string onto the default user agent.
func WithUserAgent(userAgent string) Option {
	return func(o *clientOptions) {
		o.userAgent = fmt.Sprintf("%s %s", userAgent, defaultUserAgent)
	}
}

type reqFn func(ctx context.Context, fn ...internal.RequestEditorFn) (*http.Response, error)

func byID[In, Out any](ctx context.Context, c *Client, f func(context.Context, int, ...internal.RequestEditorFn) (*http.Response, error), m func(In) (*Out, error), id int, editors []internal.RequestEditorFn) (*Out, error) {
	body, err := call(ctx, c.maxRetries, true, func(ctx context.Context, fn ...internal.RequestEditorFn) (*http.Response, error) {
		return f(ctx, id, append(fn, editors...)...)
	})
	if err != nil {
		return nil, err
	}

	var v In

	if err = json.NewDecoder(body).Decode(&v); err != nil {
		return nil, errors.Join(err, body.Close())
	}

	if err = body.Close(); err != nil {
		return nil, err
	}

	return m(v)
}

// call runs f with up to maxRetries extra attempts. A RetryAfterError uses the
// duration it carries; transient failures back off, and only when idempotent.
func call(ctx context.Context, maxRetries int, idempotent bool, f reqFn) (io.ReadCloser, error) {
	var body io.ReadCloser
	var err error

	for attempt := range maxRetries + 1 {
		body, err = doCall(ctx, f)

		if err == nil || attempt == maxRetries {
			break
		}

		var wait time.Duration

		var retryErr *httpkit.RetryAfterError

		switch {
		case errors.As(err, &retryErr):
			wait = retryErr.RetryAfter()
		case retryable(err, idempotent):
			wait = backOff(attempt)
		default:
			return body, err
		}

		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(wait):
		}
	}

	return body, err
}

func doCall(ctx context.Context, f reqFn) (io.ReadCloser, error) {
	// Captured from the outgoing request rather than res.Request, which only the
	// stock Transport fills in — a caller's own RoundTripper leaves it nil.
	var method string
	var reqURL *url.URL

	res, err := f(ctx, func(_ context.Context, req *http.Request) error {
		method = req.Method
		reqURL = req.URL

		return nil
	})
	if err != nil {
		return nil, err
	}

	// Checked before the status below: a 304 is an answer, not an unexpected status.
	if res.StatusCode == http.StatusNotModified {
		return nil, errors.Join(ErrNotModified, res.Body.Close())
	}

	if res.StatusCode != http.StatusOK && res.StatusCode != http.StatusCreated {
		body, readErr := io.ReadAll(io.LimitReader(res.Body, maxErrorBodyBytes))

		apiErr := &APIError{StatusCode: res.StatusCode, Body: body, Method: method, URL: reqURL}

		return nil, errors.Join(apiErr, readErr, res.Body.Close())
	}

	return res.Body, nil
}

type paginatedResponse[T any] interface {
	GetNext() nullable.Nullable[string]
	GetResults() []T
}

// pageEditors decides which request edits a given page receives.
type pageEditors func(page int) []internal.RequestEditorFn

// everyPage applies editors to every page request.
func everyPage(editors []internal.RequestEditorFn) pageEditors {
	return func(int) []internal.RequestEditorFn { return editors }
}

// firstPageOnly applies editors to the first page request and nothing after it,
// so a conditional listing answers 304 before any result or not at all.
func firstPageOnly(editors []internal.RequestEditorFn) pageEditors {
	return func(page int) []internal.RequestEditorFn {
		if page == 1 {
			return editors
		}

		return nil
	}
}

func idPaginate[Response paginatedResponse[In], In, Out, Params any](ctx context.Context, c *Client, f func(context.Context, int, Params, ...internal.RequestEditorFn) (*http.Response, error), m func(In) (*Out, error), id int, params func(page int) Params, editors pageEditors) iter.Seq2[*Out, error] {
	wrapped := func(ctx context.Context, p Params, fn ...internal.RequestEditorFn) (*http.Response, error) {
		return f(ctx, id, p, fn...)
	}

	return paginate[Response, In, Out](ctx, c, wrapped, m, params, editors)
}

// paginate iterates a listing, building fresh params per page so the returned
// sequence is safe to range more than once, including concurrently.
func paginate[Response paginatedResponse[In], In, Out, Params any](ctx context.Context, c *Client, f func(context.Context, Params, ...internal.RequestEditorFn) (*http.Response, error), m func(In) (*Out, error), params func(page int) Params, editors pageEditors) iter.Seq2[*Out, error] {
	return func(yield func(*Out, error) bool) {
		page := 1

		for {
			var res Response

			p, pageEditors := params(page), editors(page)

			body, err := call(ctx, c.maxRetries, true, func(ctx context.Context, fn ...internal.RequestEditorFn) (*http.Response, error) {
				return f(ctx, p, append(fn, pageEditors...)...)
			})
			if err != nil {
				yield(nil, err)
				return
			}
			if err = json.NewDecoder(body).Decode(&res); err != nil {
				yield(nil, errors.Join(err, body.Close()))
				return
			}

			if err = body.Close(); err != nil {
				yield(nil, err)
				return
			}

			for _, v := range res.GetResults() {
				if !yield(m(v)) {
					return
				}
			}

			if _, err = res.GetNext().Get(); err != nil {
				break
			}

			page++
		}
	}
}
