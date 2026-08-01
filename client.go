package metron

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"sync"
	"time"
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
	baseURL       *url.URL
	cacheDir      string
	client        *http.Client
	enableCaching bool
	rateLimit     rateLimitState
	password      string
	username      string
}

type Option func(*Client)

func NewClient(options ...Option) *Client {
	b, _ := url.Parse(baseURL)
	c := &Client{
		baseURL:       b,
		cacheDir:      "",
		client:        &http.Client{},
		enableCaching: false,
		password:      "",
		username:      "",
	}

	for _, option := range options {
		option(c)
	}

	return c
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

// rateLimitWindow tracks the most recently observed state of a single
// rate-limit window (e.g. burst or sustained), as reported by the API.
type rateLimitWindow struct {
	known     bool
	remaining int
	reset     time.Time
}

// update parses the "<prefix>Remaining" and "<prefix>Reset" headers and, if
// both are present and well-formed, records the window's state.
func (w *rateLimitWindow) update(header http.Header, prefix string) {
	remaining, err := strconv.Atoi(header.Get(prefix + "Remaining"))
	if err != nil {
		return
	}

	resetUnix, err := strconv.ParseInt(header.Get(prefix+"Reset"), 10, 64)
	if err != nil {
		return
	}

	w.known = true
	w.remaining = remaining
	w.reset = time.Unix(resetUnix, 0)
}

// wait returns how long to wait before the window allows another request.
func (w *rateLimitWindow) wait() time.Duration {
	if !w.known || w.remaining > 0 {
		return 0
	}

	return time.Until(w.reset)
}

// rateLimitState tracks the burst and sustained rate-limit windows reported
// by the metron.cloud API via its X-RateLimit-* response headers, as
// described at https://github.com/Metron-Project/metron/blob/master/api/RATELIMIT.md.
type rateLimitState struct {
	mu        sync.Mutex
	burst     rateLimitWindow
	sustained rateLimitWindow
}

// update records the rate-limit state reported by a response.
func (s *rateLimitState) update(header http.Header) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.burst.update(header, "X-RateLimit-Burst-")
	s.sustained.update(header, "X-RateLimit-Sustained-")
}

// wait returns how long to wait before the next request is allowed to proceed.
func (s *rateLimitState) wait() time.Duration {
	s.mu.Lock()
	defer s.mu.Unlock()

	return max(s.burst.wait(), s.sustained.wait())
}

func limit(ctx context.Context, c *Client) error {
	wait := c.rateLimit.wait()
	if wait <= 0 {
		return nil
	}

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(wait):
		return nil
	}
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

	c.rateLimit.update(res.Header)

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
