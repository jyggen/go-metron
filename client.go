// Package metron is a Go client library for the Metron comic book database API.
package metron

import (
	"context"
	"crypto/sha1"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"iter"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"

	"codeberg.org/jyggen/go-filecache"
	"codeberg.org/jyggen/go-httpkit"
	"codeberg.org/jyggen/go-httpkit/middleware/throttle"
	"github.com/jyggen/go-metron/internal"
	"github.com/oapi-codegen/nullable"
)

const (
	baseURL   = "https://metron.cloud"
	userAgent = "go-metron/0.1.5"
)

// Reference identifies a related resource by ID and display name.
type Reference struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// Client is a Metron API client.
type Client struct {
	cache         *filecache.FileCache
	client        internal.ClientInterface
	enableCaching bool
	httpClient    *http.Client
	maxRetries    uint
	rateLimiter   *throttle.Throttle
	storagePath   string
}

// Option configures a Client.
type Option func(*Client)

// NewClient returns a Metron client authenticated with the given username and password.
func NewClient(username, password string, options ...Option) (*Client, error) {
	storagePath, err := os.UserCacheDir()
	if err != nil {
		storagePath = os.TempDir()
	}

	c := &Client{
		enableCaching: false,
		httpClient:    &http.Client{},
		storagePath:   filepath.Join(storagePath, "go-metron"),
	}

	for _, option := range options {
		option(c)
	}

	if err = os.MkdirAll(c.storagePath, 0o700); err != nil {
		return nil, err
	}

	if c.enableCaching {
		c.cache, err = filecache.New(filecache.WithBasePath(c.storagePath), filecache.WithCompression())
		if err != nil {
			return nil, err
		}
	}

	rl, err := throttle.New(
		filepath.Join(c.storagePath, fmt.Sprintf("throttle_%x.gob", sha1.Sum([]byte(username)))),
		throttle.Limit{Count: 20, Window: time.Minute},
		throttle.Limit{Count: 5000, Window: time.Hour * 24},
	)
	if err != nil {
		return nil, err
	}

	c.rateLimiter = rl

	c.httpClient = httpkit.NewFromClient(
		c.httpClient,
		httpkit.WithBasicAuth(username, password),
		httpkit.WithUserAgent(userAgent),
		httpkit.WithMiddleware(newBackOffMiddleware()),
		httpkit.WithMiddleware(rl.Middleware()),
	)

	internalClient, err := internal.NewClient(baseURL, internal.WithHTTPClient(c.httpClient))
	if err != nil {
		return nil, err
	}

	c.client = internalClient

	return c, nil
}

// Close persists rate limit state to disk.
func (c *Client) Close() error {
	return c.rateLimiter.Close()
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
				waitTime, err := strconv.Atoi(res.Header.Get("Retry-After"))
				if err == nil {
					d := time.Duration(waitTime) * time.Second

					m.Lock()
					if t := time.Now().Add(d); t.After(backOff) {
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
	return func(c *Client) {
		c.enableCaching = true
	}
}

// WithClient sets the underlying HTTP client.
func WithClient(client *http.Client) Option {
	return func(c *Client) {
		c.httpClient = client
	}
}

// WithRetry sets the maximum number of retries for rate-limited requests.
func WithRetry(maxRetries uint) Option {
	return func(c *Client) {
		c.maxRetries = maxRetries
	}
}

// WithStoragePath sets the directory used for cache and rate limiter state.
func WithStoragePath(storagePath string) Option {
	return func(c *Client) {
		c.storagePath = storagePath
	}
}

func byID[In, Out any](ctx context.Context, cache *filecache.FileCache, maxRetries uint, key string, f func(context.Context, int, ...internal.RequestEditorFn) (*http.Response, error), m func(In) (*Out, error), id int) (*Out, error) {
	var body io.ReadCloser
	var err error

	req := func(ctx context.Context, fn ...internal.RequestEditorFn) (*http.Response, error) {
		return f(ctx, id, fn...)
	}

	if cache != nil {
		body, err = cache.Get(key, call(ctx, maxRetries, req))
	} else {
		body, _, err = call(ctx, maxRetries, req)(nil)
	}

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

func call(ctx context.Context, maxRetries uint, f func(ctx context.Context, fn ...internal.RequestEditorFn) (*http.Response, error)) func(header *filecache.Header) (io.ReadCloser, time.Time, error) {
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
		return nil, ttl, errors.Join(fmt.Errorf("unexpected status code: %d", res.StatusCode), res.Body.Close())
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

func idPaginate[Response paginatedResponse[In], In, Out any, Params paginatable](ctx context.Context, cache *filecache.FileCache, maxRetries uint, kind string, f func(context.Context, int, Params, ...internal.RequestEditorFn) (*http.Response, error), m func(In) (*Out, error), id int, params Params) iter.Seq2[*Out, error] {
	wrapped := func(ctx context.Context, p Params, fn ...internal.RequestEditorFn) (*http.Response, error) {
		return f(ctx, id, p, fn...)
	}

	return paginate[Response, In, Out](ctx, cache, maxRetries, fmt.Sprintf("%s/%d", kind, id), wrapped, m, params)
}

func paginate[Response paginatedResponse[In], In, Out any, Params paginatable](ctx context.Context, cache *filecache.FileCache, maxRetries uint, kind string, f func(context.Context, Params, ...internal.RequestEditorFn) (*http.Response, error), m func(In) (*Out, error), params Params) iter.Seq2[*Out, error] {
	return func(yield func(*Out, error) bool) {
		page := 1

		for {
			var res Response
			params.SetPage(page)
			req := func(ctx context.Context, fn ...internal.RequestEditorFn) (*http.Response, error) {
				return f(ctx, params, fn...)
			}
			var body io.ReadCloser
			var err error
			if cache != nil {
				body, err = cache.Get(cacheKey(kind, params), call(ctx, maxRetries, req))
			} else {
				body, _, err = call(ctx, maxRetries, req)(nil)
			}
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
