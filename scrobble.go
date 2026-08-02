package metron

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"cloud.google.com/go/civil"
	"github.com/jyggen/go-metron/internal"
	"github.com/oapi-codegen/nullable"
)

// ScrobbleIssueSeries is the series an issue belongs to, as embedded in a
// scrobble response.
type ScrobbleIssueSeries struct {
	Name      string
	Volume    int
	YearBegan int
}

// ScrobbleIssue is the issue summary returned in a scrobble response.
type ScrobbleIssue struct {
	ID        int
	Number    string
	CoverDate civil.Date
	StoreDate *civil.Date
	Series    ScrobbleIssueSeries
	Modified  time.Time
}

// ScrobbleResult is the result of marking an issue as read.
type ScrobbleResult struct {
	ID       int
	Issue    ScrobbleIssue
	IsRead   bool
	ReadDate *time.Time
	Rating   *int
	Created  bool
	Modified time.Time
}

// scrobbleOptions mirrors the optional fields of internal.ScrobbleRequest so
// ScrobbleOption never names a type from the internal package. Scrobble
// translates it into the request just before the call.
type scrobbleOptions struct {
	readDate *time.Time
	rating   *int
}

// ScrobbleOption configures an optional parameter for a scrobble request.
type ScrobbleOption func(*scrobbleOptions)

// WithReadDate sets the read date for a scrobble request.
func WithReadDate(t time.Time) ScrobbleOption {
	return func(o *scrobbleOptions) {
		o.readDate = &t
	}
}

// WithRating sets the rating (1-5) for a scrobble request.
func WithRating(rating int) ScrobbleOption {
	return func(o *scrobbleOptions) {
		o.rating = &rating
	}
}

// Scrobble marks an issue as read, optionally setting a read date and rating.
func (c *Client) Scrobble(ctx context.Context, issueID int, opts ...ScrobbleOption) (*ScrobbleResult, error) {
	var o scrobbleOptions

	for _, opt := range opts {
		opt(&o)
	}

	req := internal.ScrobbleRequest{
		IssueId: issueID,
	}

	if o.readDate != nil {
		req.DateRead = nullable.NewNullableWithValue(*o.readDate)
	}

	if o.rating != nil {
		req.Rating = nullable.NewNullableWithValue(*o.rating)
	}

	body, _, err := call(ctx, c.maxRetries, func(ctx context.Context, fn ...internal.RequestEditorFn) (*http.Response, error) {
		return c.client.ApiCollectionScrobbleCreate(ctx, req, fn...)
	})(nil)
	if err != nil {
		return nil, err
	}

	var v internal.ScrobbleResponse

	if err = json.NewDecoder(body).Decode(&v); err != nil {
		return nil, errors.Join(err, body.Close())
	}

	if err = body.Close(); err != nil {
		return nil, err
	}

	return scrobbleMapper(v)
}

func scrobbleMapper(in internal.ScrobbleResponse) (*ScrobbleResult, error) {
	if in.Id == nil {
		return nil, fmt.Errorf("scrobble: nil Id")
	}

	if in.Created == nil {
		return nil, fmt.Errorf("scrobble: nil Created")
	}

	if in.IsRead == nil {
		return nil, fmt.Errorf("scrobble: nil IsRead")
	}

	if in.Modified == nil {
		return nil, fmt.Errorf("scrobble: nil Modified")
	}

	if in.Issue == nil {
		return nil, fmt.Errorf("scrobble: nil Issue")
	}

	issue, err := collectionIssueMapper(*in.Issue)
	if err != nil {
		return nil, err
	}

	var rating *int

	if ratingVal, rErr := in.Rating.Get(); rErr == nil {
		if enumVal, eErr := ratingVal.AsRatingEnum(); eErr == nil {
			rating = new(int(enumVal))
		}
	}

	return &ScrobbleResult{
		ID:       *in.Id,
		Issue:    *issue,
		IsRead:   *in.IsRead,
		ReadDate: nullableToPtr(in.DateRead),
		Rating:   rating,
		Created:  *in.Created,
		Modified: *in.Modified,
	}, nil
}

func collectionIssueMapper(in internal.CollectionIssue) (*ScrobbleIssue, error) {
	if in.Id == nil {
		return nil, fmt.Errorf("scrobble: nil Issue.Id")
	}

	if in.Modified == nil {
		return nil, fmt.Errorf("scrobble: nil Issue.Modified")
	}

	if in.Series == nil {
		return nil, fmt.Errorf("scrobble: nil Issue.Series")
	}

	coverDate := civil.DateOf(in.CoverDate.Time)

	var maybeStoreDate *civil.Date

	if d, err := in.StoreDate.Get(); err == nil {
		maybeStoreDate = new(civil.DateOf(d.Time))
	}

	return &ScrobbleIssue{
		ID:        *in.Id,
		Number:    in.Number,
		CoverDate: coverDate,
		StoreDate: maybeStoreDate,
		Series: ScrobbleIssueSeries{
			Name:      in.Series.Name,
			Volume:    in.Series.Volume,
			YearBegan: in.Series.YearBegan,
		},
		Modified: *in.Modified,
	}, nil
}
