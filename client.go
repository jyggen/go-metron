package metron

import (
	"context"
	"encoding/json"
	"fmt"
	"go.uber.org/ratelimit"
	"io"
	"net/http"
	"net/url"
	"time"
)

const baseURL = "https://metron.cloud/api/"
const userAgent = "go-metron/0.2-dev (+https://github.com/jyggen/go-metron)"

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
	limiter       ratelimit.Limiter
	password      string
	username      string
}

type Option func(*Client)

func NewClient(options ...Option) *Client {
	b, _ := url.Parse(baseURL)
	c := &Client{
		baseURL:       b,
		cacheDir:      "",
		client:        http.DefaultClient,
		enableCaching: false,
		limiter:       ratelimit.New(30, ratelimit.Per(time.Minute)),
		password:      "",
		username:      "",
	}

	for _, option := range options {
		option(c)
	}

	return c
}

// WithAuthentication sets the username and password to be used for authentication.
func WithAuthentication(username string, password string) Option {
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

func do[T any](c *Client, req *http.Request) (T, error) {
	var v T

	c.limiter.Take()

	res, err := c.client.Do(req)

	if err != nil {
		return v, err
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return v, fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}

	body, err := io.ReadAll(res.Body)

	if err != nil {
		return v, err
	}

	err = json.Unmarshal(body, &v)

	return v, err
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

	return cache[T](c, req)
}
