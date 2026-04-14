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
	"time"

	"codeberg.org/jyggen/go-filecache"
	"codeberg.org/jyggen/go-httpkit"
	"github.com/jyggen/go-metron/internal"
	"github.com/oapi-codegen/nullable"
)

func nullableToPtr[T any](n nullable.Nullable[T]) *T {
	if v, err := n.Get(); err == nil {
		return &v
	}

	return nil
}

// Arc is a story arc.
type Arc struct {
	ID                    int
	Name                  string
	Description           *string
	ImageURL              *url.URL
	ComicVineID           *int
	GrandComicsDatabaseID *int
	ResourceURL           url.URL
	Modified              time.Time
}

// ArcList is a story arc as it appears in list responses.
type ArcList struct {
	ID       int
	Name     string
	Modified time.Time
}

// ArcByID returns a story arc by its ID.
func (c *Client) ArcByID(ctx context.Context, id int) (*Arc, error) {
	return newByID(ctx, c.cache, c.maxRetries, fmt.Sprintf("arc/%d", id), c.client.ApiArcRetrieve, arcMapper, id)
}

// Arcs returns an iterator over all story arcs.
func (c *Client) Arcs(ctx context.Context, filters ...Filter) iter.Seq2[*ArcList, error] {
	params := &internal.ApiArcListParams{}

	for _, f := range filters {
		f(params)
	}

	return newPaginate[internal.PaginatedArcListList](ctx, c.cache, c.maxRetries, "arc", c.client.ApiArcList, arcListMapper, params)
}

func arcMapper(in internal.Arc) (*Arc, error) {
	if in.Id == nil {
		return nil, fmt.Errorf("arc: nil Id")
	}

	if in.Modified == nil {
		return nil, fmt.Errorf("arc: nil Modified")
	}

	if in.ResourceUrl == nil {
		return nil, fmt.Errorf("arc: nil ResourceUrl")
	}

	var imageURL *url.URL
	var err error

	if in.Image != nil {
		imageURL, err = url.Parse(*in.Image)
		if err != nil {
			return nil, err
		}
	}

	resourceURL, err := url.Parse(*in.ResourceUrl)
	if err != nil {
		return nil, err
	}

	return &Arc{
		ID:                    *in.Id,
		Name:                  in.Name,
		Description:           in.Desc,
		ImageURL:              imageURL,
		ComicVineID:           nullableToPtr(in.CvId),
		GrandComicsDatabaseID: nullableToPtr(in.GcdId),
		ResourceURL:           *resourceURL,
		Modified:              *in.Modified,
	}, nil
}

func arcListMapper(in internal.ArcList) (*ArcList, error) {
	if in.Id == nil {
		return nil, fmt.Errorf("arc: nil Id")
	}

	if in.Modified == nil {
		return nil, fmt.Errorf("arc: nil Modified")
	}

	return &ArcList{
		ID:       *in.Id,
		Name:     in.Name,
		Modified: *in.Modified,
	}, nil
}

func newByID[In, Out any](ctx context.Context, cache *filecache.FileCache, maxRetries uint, key string, f func(context.Context, int, ...internal.RequestEditorFn) (*http.Response, error), m func(In) (*Out, error), id int) (*Out, error) {
	var body io.ReadCloser
	var err error

	req := func(ctx context.Context, fn ...internal.RequestEditorFn) (*http.Response, error) {
		return f(ctx, id, fn...)
	}

	if cache != nil {
		body, err = cache.Get(key, newCall(ctx, maxRetries, req))
	} else {
		body, _, err = newCall(ctx, maxRetries, req)(nil)
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

func newCall(ctx context.Context, maxRetries uint, f func(ctx context.Context, fn ...internal.RequestEditorFn) (*http.Response, error)) func(header *filecache.Header) (io.ReadCloser, time.Time, error) {
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

func cacheKeyWithID(kind string, id int, params any) string {
	b, _ := json.Marshal(params)
	return fmt.Sprintf("%s/%d/%s", kind, id, string(b))
}

func newIDPaginate[Response paginatedResponse[In], In, Out any, Params paginatable](ctx context.Context, cache *filecache.FileCache, maxRetries uint, kind string, call func(context.Context, int, Params, ...internal.RequestEditorFn) (*http.Response, error), m func(In) (*Out, error), id int, params Params) iter.Seq2[*Out, error] {
	return func(yield func(*Out, error) bool) {
		page := 1

		for {
			var res Response
			params.SetPage(page)
			req := func(ctx context.Context, fn ...internal.RequestEditorFn) (*http.Response, error) {
				return call(ctx, id, params, fn...)
			}
			var body io.ReadCloser
			var err error
			if cache != nil {
				body, err = cache.Get(cacheKeyWithID(kind, id, params), newCall(ctx, maxRetries, req))
			} else {
				body, _, err = newCall(ctx, maxRetries, req)(nil)
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

func newPaginate[Response paginatedResponse[In], In, Out any, Params paginatable](ctx context.Context, cache *filecache.FileCache, maxRetries uint, kind string, call func(context.Context, Params, ...internal.RequestEditorFn) (*http.Response, error), m func(In) (*Out, error), params Params) iter.Seq2[*Out, error] {
	return func(yield func(*Out, error) bool) {
		page := 1

		for {
			var res Response
			params.SetPage(page)
			req := func(ctx context.Context, fn ...internal.RequestEditorFn) (*http.Response, error) {
				return call(ctx, params, fn...)
			}
			var body io.ReadCloser
			var err error
			if cache != nil {
				body, err = cache.Get(cacheKey(kind, params), newCall(ctx, maxRetries, req))
			} else {
				body, _, err = newCall(ctx, maxRetries, req)(nil)
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
