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
	"github.com/jyggen/go-metron/internal"
)

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

type ArcList struct {
	ID       int
	Name     string
	Modified time.Time
}

// ArcByID returns the information of an individual story arc.
func (c *Client) ArcByID(ctx context.Context, id int) (*Arc, error) {
	return newByID(ctx, c.cache, fmt.Sprintf("arc/%d", id), c.client.ApiArcRetrieve, arcMapper, id)
}

// Arcs returns a list of all the story arcs.
func (c *Client) Arcs(ctx context.Context, filters ...Filter) iter.Seq2[*ArcList, error] {
	params := &internal.ApiArcListParams{}

	for _, f := range filters {
		f(params)
	}

	return newPaginate[internal.PaginatedArcListList](ctx, c.client.ApiArcList, arcListMapper, params)
}

func arcMapper(in internal.Arc) (*Arc, error) {
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
		ComicVineID:           in.CvId,
		GrandComicsDatabaseID: in.GcdId,
		ResourceURL:           *resourceURL,
		Modified:              *in.Modified,
	}, nil
}

func arcListMapper(in internal.ArcList) (*ArcList, error) {
	return &ArcList{
		ID:       *in.Id,
		Name:     in.Name,
		Modified: *in.Modified,
	}, nil
}

func newByID[In, Out any](ctx context.Context, cache *filecache.FileCache, key string, f func(context.Context, int, ...internal.RequestEditorFn) (*http.Response, error), m func(In) (*Out, error), id int) (*Out, error) {
	var body io.ReadCloser
	var err error

	req := func(ctx context.Context, fn ...internal.RequestEditorFn) (*http.Response, error) {
		return f(ctx, id, fn...)
	}

	if cache != nil {
		body, err = cache.Get(key, newCall(ctx, req))
	} else {
		body, _, err = newCall(ctx, req)(nil)
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

func newCall(ctx context.Context, f func(ctx context.Context, fn ...internal.RequestEditorFn) (*http.Response, error)) func(header *filecache.Header) (io.ReadCloser, time.Time, error) {
	return func(header *filecache.Header) (io.ReadCloser, time.Time, error) {
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

		if res.StatusCode == http.StatusTooManyRequests {
			if err = res.Body.Close(); err != nil {
				return nil, ttl, err
			}

			return newCall(ctx, f)(header)
		}

		lastModifiedHeader := res.Header.Get("Last-Modified")

		if lastModifiedHeader != "" {
			lastModified, innerErr := http.ParseTime(lastModifiedHeader)
			if innerErr != nil {
				return nil, ttl, errors.Join(innerErr, res.Body.Close())
			}

			ttl = ttl.Add(ttl.Sub(lastModified) / 10)
		}

		if res.StatusCode == http.StatusNotModified {
			return nil, ttl, res.Body.Close()
		}

		if res.StatusCode != http.StatusOK {
			return nil, ttl, errors.Join(fmt.Errorf("unexpected status code: %d", res.StatusCode), res.Body.Close())
		}

		return res.Body, ttl, nil
	}
}

type paginatable interface {
	SetPage(v int)
}

type paginatedResponse[T any] interface {
	GetNext() *string
	GetResults() []T
}

func newIDPaginate[Response paginatedResponse[In], In, Out any, Params paginatable](ctx context.Context, call func(context.Context, int, Params, ...internal.RequestEditorFn) (*http.Response, error), m func(In) (*Out, error), id int, params Params) iter.Seq2[*Out, error] {
	return func(yield func(*Out, error) bool) {
		page := 1

		var res Response

		for {
			body, _, err := newCall(ctx, func(ctx context.Context, fn ...internal.RequestEditorFn) (*http.Response, error) {
				params.SetPage(page)
				return call(ctx, id, params, fn...)
			})(nil)
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

			if res.GetNext() == nil {
				break
			}

			page++
		}
	}
}

func newPaginate[Response paginatedResponse[In], In, Out any, Params paginatable](ctx context.Context, call func(context.Context, Params, ...internal.RequestEditorFn) (*http.Response, error), m func(In) (*Out, error), params Params) iter.Seq2[*Out, error] {
	return func(yield func(*Out, error) bool) {
		page := 1

		var res Response

		for {
			body, _, err := newCall(ctx, func(ctx context.Context, fn ...internal.RequestEditorFn) (*http.Response, error) {
				params.SetPage(page)
				return call(ctx, params, fn...)
			})(nil)
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

			if res.GetNext() == nil {
				break
			}

			page++
		}
	}
}
