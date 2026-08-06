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
	"os"
	"path/filepath"
	"sync"
	"time"

	"codeberg.org/jyggen/go-filecache"
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
	cache      *filecache.FileCache
	client     internal.ClientInterface
	maxRetries int
}

// clientOptions is what an Option mutates, so options can be applied in any
// order before NewClient derives the Client from them.
type clientOptions struct {
	baseURL     string
	caching     bool
	httpClient  *http.Client
	maxRetries  int
	storagePath string
	userAgent   string
}

// Option configures a Client.
type Option func(*clientOptions)

func defaultOptions() clientOptions {
	storagePath, err := os.UserCacheDir()
	if err != nil {
		storagePath = os.TempDir()
	}

	return clientOptions{
		baseURL:     defaultBaseURL,
		httpClient:  &http.Client{},
		storagePath: filepath.Join(storagePath, "go-metron"),
		userAgent:   defaultUserAgent,
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

	if err := os.MkdirAll(o.storagePath, 0o700); err != nil {
		return nil, err
	}

	c := &Client{maxRetries: o.maxRetries}

	if o.caching {
		cache, err := filecache.New(filecache.WithBasePath(o.storagePath), filecache.WithCompression())
		if err != nil {
			return nil, err
		}

		c.cache = cache
	}

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

// WithCaching enables on-disk response caching.
func WithCaching() Option {
	return func(o *clientOptions) {
		o.caching = true
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

// WithStoragePath sets the directory used for cache and rate limiter state.
func WithStoragePath(storagePath string) Option {
	return func(o *clientOptions) {
		o.storagePath = storagePath
	}
}

// WithUserAgent prefixes the given string onto the default user agent.
func WithUserAgent(userAgent string) Option {
	return func(o *clientOptions) {
		o.userAgent = fmt.Sprintf("%s %s", userAgent, defaultUserAgent)
	}
}

// applyFilters applies each filter, stopping at the first unsupported one.
// endpoint is the calling method's name, which only it knows and which a
// FilterError is meaningless without.
func applyFilters(endpoint string, params any, filters []Filter) error {
	for _, f := range filters {
		if err := f(params); err != nil {
			var filterErr *FilterError

			if errors.As(err, &filterErr) {
				filterErr.Endpoint = endpoint
			}

			return err
		}
	}

	return nil
}

// errIter returns an iterator yielding err once, letting a list method report a
// failure without changing its signature.
func errIter[T any](err error) iter.Seq2[T, error] {
	return func(yield func(T, error) bool) {
		var zero T

		yield(zero, err)
	}
}

type reqFn func(ctx context.Context, fn ...internal.RequestEditorFn) (*http.Response, error)

func (c *Client) fetch(ctx context.Context, key string, req reqFn) (io.ReadCloser, error) {
	if c.cache != nil {
		return c.cache.Get(key, call(ctx, c.maxRetries, true, req))
	}

	body, _, err := call(ctx, c.maxRetries, true, req)(nil)

	return body, err
}

func byID[In, Out any](ctx context.Context, c *Client, key string, f func(context.Context, int, ...internal.RequestEditorFn) (*http.Response, error), m func(In) (*Out, error), id int) (*Out, error) {
	body, err := c.fetch(ctx, key, func(ctx context.Context, fn ...internal.RequestEditorFn) (*http.Response, error) {
		return f(ctx, id, fn...)
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
func call(ctx context.Context, maxRetries int, idempotent bool, f func(ctx context.Context, fn ...internal.RequestEditorFn) (*http.Response, error)) func(header *filecache.Header) (io.ReadCloser, time.Time, error) {
	return func(header *filecache.Header) (io.ReadCloser, time.Time, error) {
		var body io.ReadCloser
		var ttl time.Time
		var err error

		for attempt := range maxRetries + 1 {
			body, ttl, err = doCall(ctx, f, header)

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
				return body, ttl, err
			}

			select {
			case <-ctx.Done():
				return nil, time.Now(), ctx.Err()
			case <-time.After(wait):
			}
		}

		return body, ttl, err
	}
}

func doCall(ctx context.Context, f func(ctx context.Context, fn ...internal.RequestEditorFn) (*http.Response, error), header *filecache.Header) (io.ReadCloser, time.Time, error) {
	// Captured from the outgoing request rather than res.Request, which only the
	// stock Transport fills in — a caller's own RoundTripper leaves it nil.
	var method string
	var reqURL *url.URL

	res, err := f(ctx, func(ctx context.Context, req *http.Request) error {
		method = req.Method
		reqURL = req.URL

		if header != nil {
			req.Header.Set("If-Modified-Since", time.Unix(0, header.FetchedAt).Format(http.TimeFormat))
		}

		return nil
	})

	// Without Last-Modified there is no basis for a TTL, so the entry is written
	// already expired and never served as fresh. It still earns its keep: the
	// stored FetchedAt drives If-Modified-Since, so the API can answer 304.
	now := time.Now()
	ttl := now

	if err != nil {
		return nil, ttl, err
	}

	lastModifiedHeader := res.Header.Get("Last-Modified")

	if lastModifiedHeader != "" {
		lastModified, err := http.ParseTime(lastModifiedHeader)
		if err != nil {
			return nil, ttl, errors.Join(err, res.Body.Close())
		}

		ttl = now.Add(now.Sub(lastModified) / 10)
	}

	if res.StatusCode == http.StatusNotModified {
		return nil, ttl, res.Body.Close()
	}

	if res.StatusCode != http.StatusOK && res.StatusCode != http.StatusCreated {
		body, readErr := io.ReadAll(io.LimitReader(res.Body, maxErrorBodyBytes))

		apiErr := &APIError{StatusCode: res.StatusCode, Body: body, Method: method, URL: reqURL}

		return nil, ttl, errors.Join(apiErr, readErr, res.Body.Close())
	}

	return res.Body, ttl, nil
}

type paginatable interface {
	SetPage(v int)
}

type paginatedResponse[T any] interface {
	GetNext() nullable.Nullable[string]
	GetResults() []T
}

// cacheKey serialises params into a key; go-filecache hashes it, so it need only
// be stable. The marshal error is propagated because time.Time params reject
// years outside [0,9999], and dropping it would collapse them onto one key.
func cacheKey(kind string, params any) (string, error) {
	b, err := json.Marshal(params)
	if err != nil {
		return "", fmt.Errorf("metron: cache key: %w", err)
	}

	return kind + "/" + string(b), nil
}

func idPaginate[Response paginatedResponse[In], In, Out any, Params paginatable](ctx context.Context, c *Client, kind string, f func(context.Context, int, Params, ...internal.RequestEditorFn) (*http.Response, error), m func(In) (*Out, error), id int, params Params) iter.Seq2[*Out, error] {
	wrapped := func(ctx context.Context, p Params, fn ...internal.RequestEditorFn) (*http.Response, error) {
		return f(ctx, id, p, fn...)
	}

	return paginate[Response, In, Out](ctx, c, fmt.Sprintf("%s/%d", kind, id), wrapped, m, params)
}

func paginate[Response paginatedResponse[In], In, Out any, Params paginatable](ctx context.Context, c *Client, kind string, f func(context.Context, Params, ...internal.RequestEditorFn) (*http.Response, error), m func(In) (*Out, error), params Params) iter.Seq2[*Out, error] {
	return func(yield func(*Out, error) bool) {
		page := 1

		for {
			var res Response
			params.SetPage(page)

			key, err := cacheKey(kind, params)
			if err != nil {
				yield(nil, err)
				return
			}

			body, err := c.fetch(ctx, key, func(ctx context.Context, fn ...internal.RequestEditorFn) (*http.Response, error) {
				return f(ctx, params, fn...)
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
