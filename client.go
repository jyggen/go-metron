package metron

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/AliRizaAynaci/gorl"
	"github.com/AliRizaAynaci/gorl/core"
	"github.com/cenkalti/backoff/v5"
)

const (
	baseURL   = "https://metron.cloud/api/"
	userAgent = "go-metron/0.1.5"
)

type listTypes interface {
	ArcList | CharacterList | CreatorList | ImprintList | IssueList | PublisherList | RoleList | SeriesList | SeriesTypeList | TeamList | UniverseList
}

type paginatedList[T listTypes] struct {
	Count    int     `json:"count"`
	Next     *string `json:"next"`
	Previous *string `json:"previous"`
	Results  []T     `json:"results"`
}

type Reference struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type URL struct {
	*url.URL
}

func (u URL) MarshalJSON() ([]byte, error) {
	return json.Marshal(u.String())
}

func (u *URL) UnmarshalJSON(b []byte) error {
	var err error
	var s string

	if err = json.Unmarshal(b, &s); err != nil {
		return err
	}

	u.URL, err = url.Parse(s)

	return err
}

type Client struct {
	baseURL          *url.URL
	cacheDir         string
	client           *http.Client
	enableCaching    bool
	limiterBurst     core.Limiter
	limiterSustained core.Limiter
	password         string
	username         string
}

type Option func(*Client)

func NewClient(client *http.Client, options ...Option) (*Client, error) {
	burst, err := gorl.New(core.Config{
		Strategy: core.SlidingWindow,
		KeyBy:    core.KeyByAPIKey,
		Limit:    30,
		Window:   1 * time.Minute,
	})
	if err != nil {
		return nil, err
	}

	sustained, err := gorl.New(core.Config{
		Strategy: core.SlidingWindow,
		KeyBy:    core.KeyByAPIKey,
		Limit:    10_000,
		Window:   24 * time.Hour,
	})
	if err != nil {
		return nil, err
	}

	b, err := url.Parse(baseURL)
	if err != nil {
		return nil, err
	}

	c := &Client{
		baseURL:          b,
		cacheDir:         "",
		client:           client,
		enableCaching:    false,
		limiterBurst:     burst,
		limiterSustained: sustained,
		password:         "",
		username:         "",
	}

	for _, option := range options {
		option(c)
	}

	return c, nil
}

// WithAuthentication sets the username and password to be used for authentication.
func WithAuthentication(username, password string) Option {
	return func(c *Client) {
		c.password = password
		c.username = username
	}
}

// WithCaching enables heuristic caching. If cacheDir is empty, it will use the operating system's default location
// for user-specific cached data.
func WithCaching(cacheDir string) Option {
	return func(c *Client) {
		c.cacheDir = cacheDir
		c.enableCaching = true
	}
}

func WithClient(client *http.Client) Option {
	return func(c *Client) {
		c.client = client
	}
}

func WithTimeout(timeout time.Duration) Option {
	return func(c *Client) {
		c.client.Timeout = timeout
	}
}

func paginate[T listTypes](ctx context.Context, c *Client, path string, filters ...Filter) func(func(T, error) bool) {
	return func(yield func(T, error) bool) {
		u, err := c.baseURL.Parse(path)

		var v T

		if err != nil {
			yield(v, err)
			return
		}

		var vList paginatedList[T]

		for {
			vList, err = request[paginatedList[T]](ctx, c, u.String(), filters...)
			if err != nil {
				yield(v, err)
				return
			}

			for _, v = range vList.Results {
				if !yield(v, err) {
					return
				}
			}

			if vList.Next == nil {
				break
			}

			u, err = c.baseURL.Parse(*vList.Next)
			if err != nil {
				yield(v, err)
				return
			}

			filters = nil
		}
	}
}

func limitWithRetry(ctx context.Context, limiter core.Limiter, username string) error {
	_, err := backoff.Retry(ctx, func() (bool, error) {
		allowed, innerErr := limiter.Allow(username)

		if innerErr != nil {
			return false, innerErr
		}

		if !allowed {
			return false, errors.New("rate limit exceeded")
		}

		return true, nil
	}, backoff.WithBackOff(backoff.NewExponentialBackOff()))

	return err
}

func limit(ctx context.Context, c *Client) error {
	return errors.Join(
		limitWithRetry(ctx, c.limiterBurst, c.username),
		limitWithRetry(ctx, c.limiterSustained, c.username),
	)
}

func do[T any](ctx context.Context, c *Client, req *http.Request) (T, error) {
	var v T

	err := limit(ctx, c)
	if err != nil {
		return v, err
	}

	res, err := c.client.Do(req)
	if err != nil {
		return v, err
	}

	if res.StatusCode == http.StatusTooManyRequests {
		if err = res.Body.Close(); err != nil {
			return v, err
		}

		waitTime, innerErr := strconv.Atoi(res.Header.Get("Retry-After"))

		if innerErr != nil {
			return v, innerErr
		}

		select {
		case <-ctx.Done():
			return v, ctx.Err()
		case <-time.After(time.Duration(waitTime) * time.Second):
			return do[T](ctx, c, req)
		}
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return v, fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}

	return v, json.NewDecoder(res.Body).Decode(&v)
}

func request[T any](ctx context.Context, c *Client, path string, filters ...Filter) (T, error) {
	var v T

	u, err := c.baseURL.Parse(path)
	if err != nil {
		return v, err
	}

	q := u.Query()

	for _, f := range filters {
		f(&q)
	}

	u.RawQuery = q.Encode()

	var req *http.Request

	req, err = http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return v, err
	}

	req.SetBasicAuth(c.username, c.password)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", userAgent)

	return cache[T](ctx, c, req)
}
