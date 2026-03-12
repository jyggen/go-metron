package metron

import (
	"crypto/sha1"
	"errors"
	"fmt"
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
)

const (
	baseURL   = "https://metron.cloud"
	userAgent = "go-metron/0.1.5"
)

type Reference struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Client struct {
	cache         *filecache.FileCache
	client        internal.ClientInterface
	enableCaching bool
	httpClient    *http.Client
	rateLimiter   *throttle.Throttle
	storagePath   string
}

type Option func(*Client)

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

func WithCaching() Option {
	return func(c *Client) {
		c.enableCaching = true
	}
}

func WithClient(client *http.Client) Option {
	return func(c *Client) {
		c.httpClient = client
	}
}

func WithStoragePath(storagePath string) Option {
	return func(c *Client) {
		c.storagePath = storagePath
	}
}
