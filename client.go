// Package metron is a Go client library for the Metron comic book database API.
package metron

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"iter"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"

	"codeberg.org/jyggen/go-filecache"
	"codeberg.org/jyggen/go-httpkit"
	"github.com/jyggen/go-metron/internal"
	"github.com/oapi-codegen/nullable"
)

const (
	baseURL          = "https://metron.cloud"
	defaultUserAgent = "go-metron/0.1.5"
)

// Reference identifies a related resource by ID and display name. It is used
// for every kind of relation the API returns — publishers, roles, genres,
// ratings and so on — so the meaning of a Reference depends on the field
// holding it.
//
// ID is the resource's ID within its own collection and is the only stable
// identifier: pass it to the matching ByID method to resolve the full record.
// Name is the upstream display name, which may be edited or reformatted at any
// time, and is not unique. Do not key off Name.
type Reference struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// Client is a Metron API client. It holds only the state used to serve
// requests; everything supplied via Option is consumed during construction.
type Client struct {
	cache      *filecache.FileCache
	client     internal.ClientInterface
	maxRetries int
}

// clientOptions is the mutable configuration an Option acts on. It exists so
// options can be applied in any order before NewClient derives a Client from
// the final values.
type clientOptions struct {
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

	internalClient, err := internal.NewClient(baseURL, internal.WithHTTPClient(httpClient))
	if err != nil {
		return nil, err
	}

	c.client = internalClient

	return c, nil
}

// Close is a no-op kept for backward compatibility. Rate-limit state is now
// tracked in memory from response headers rather than persisted to disk.
func (c *Client) Close() error {
	return nil
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

type reqFn func(ctx context.Context, fn ...internal.RequestEditorFn) (*http.Response, error)

func (c *Client) fetch(ctx context.Context, key string, req reqFn) (io.ReadCloser, error) {
	if c.cache != nil {
		return c.cache.Get(key, call(ctx, c.maxRetries, req))
	}

	body, _, err := call(ctx, c.maxRetries, req)(nil)

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

func call(ctx context.Context, maxRetries int, f func(ctx context.Context, fn ...internal.RequestEditorFn) (*http.Response, error)) func(header *filecache.Header) (io.ReadCloser, time.Time, error) {
	return func(header *filecache.Header) (io.ReadCloser, time.Time, error) {
		var body io.ReadCloser
		var ttl time.Time
		var err error

		for attempt := range maxRetries + 1 {
			body, ttl, err = doCall(ctx, f, header)

			var retryErr *httpkit.RetryAfterError
			if attempt < maxRetries && errors.As(err, &retryErr) {
				select {
				case <-ctx.Done():
					return nil, time.Now(), ctx.Err()
				case <-time.After(retryErr.RetryAfter()):
					continue
				}
			}

			break
		}

		return body, ttl, err
	}
}

func doCall(ctx context.Context, f func(ctx context.Context, fn ...internal.RequestEditorFn) (*http.Response, error), header *filecache.Header) (io.ReadCloser, time.Time, error) {
	res, err := f(ctx, func(ctx context.Context, req *http.Request) error {
		if header != nil {
			req.Header.Set("If-Modified-Since", time.Unix(0, header.FetchedAt).Format(http.TimeFormat))
		}

		return nil
	})

	ttl := time.Now()

	if err != nil {
		return nil, ttl, err
	}

	lastModifiedHeader := res.Header.Get("Last-Modified")

	if lastModifiedHeader != "" {
		lastModified, err := http.ParseTime(lastModifiedHeader)
		if err != nil {
			return nil, ttl, errors.Join(err, res.Body.Close())
		}

		ttl = ttl.Add(ttl.Sub(lastModified) / 10)
	}

	if res.StatusCode == http.StatusNotModified {
		return nil, ttl, res.Body.Close()
	}

	if res.StatusCode != http.StatusOK && res.StatusCode != http.StatusCreated {
		body, readErr := io.ReadAll(io.LimitReader(res.Body, maxErrorBodyBytes))

		return nil, ttl, errors.Join(&APIError{StatusCode: res.StatusCode, Body: body}, readErr, res.Body.Close())
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

func cacheKey(kind string, params any) string {
	b, _ := json.Marshal(params)
	return kind + "/" + string(b)
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
			body, err := c.fetch(ctx, cacheKey(kind, params), func(ctx context.Context, fn ...internal.RequestEditorFn) (*http.Response, error) {
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
