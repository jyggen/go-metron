package metron

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"cloud.google.com/go/civil"
	"github.com/jyggen/go-metron/internal"
	"github.com/oapi-codegen/nullable"
)

// ScrobbleIssueSeries is the series an issue belongs to, as embedded in a
// scrobble response.
type ScrobbleIssueSeries struct {
	ID        int
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

// scrobbleOptions mirrors internal.ScrobbleRequest's optional fields so
// ScrobbleOption never names an internal type.
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

	// Not idempotent: a retried POST could scrobble the issue twice.
	body, _, err := call(ctx, c.maxRetries, false, func(ctx context.Context, fn ...internal.RequestEditorFn) (*http.Response, error) {
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
		return nil, &MapError{Kind: "scrobble", Field: "Id"}
	}

	id := *in.Id

	if in.Created == nil {
		return nil, &MapError{Kind: "scrobble", ID: id, Field: "Created"}
	}

	if in.IsRead == nil {
		return nil, &MapError{Kind: "scrobble", ID: id, Field: "IsRead"}
	}

	if in.Modified == nil {
		return nil, &MapError{Kind: "scrobble", ID: id, Field: "Modified"}
	}

	if in.Issue == nil {
		return nil, &MapError{Kind: "scrobble", ID: id, Field: "Issue"}
	}

	issue, err := collectionIssueMapper(*in.Issue)
	if err != nil {
		return nil, err
	}

	var rating *int

	// A null or absent rating is legitimate and Get reports both as an error, so
	// that one is ignored. A present but non-numeric rating is malformed.
	if ratingVal, rErr := in.Rating.Get(); rErr == nil {
		enumVal, eErr := ratingVal.AsRatingEnum()
		if eErr != nil {
			return nil, &MapError{Kind: "scrobble", ID: id, Field: "Rating", Err: eErr}
		}

		rating = new(int(enumVal))
	}

	return &ScrobbleResult{
		ID:       id,
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
		return nil, &MapError{Kind: "issue", Field: "Id"}
	}

	id := *in.Id

	if in.Modified == nil {
		return nil, &MapError{Kind: "issue", ID: id, Field: "Modified"}
	}

	if in.Series == nil {
		return nil, &MapError{Kind: "issue", ID: id, Field: "Series"}
	}

	if in.Series.Id == nil {
		return nil, &MapError{Kind: "issue", ID: id, Field: "Series.Id"}
	}

	coverDate := civil.DateOf(in.CoverDate.Time)

	var maybeStoreDate *civil.Date

	if d, err := in.StoreDate.Get(); err == nil {
		maybeStoreDate = new(civil.DateOf(d.Time))
	}

	return &ScrobbleIssue{
		ID:        id,
		Number:    in.Number,
		CoverDate: coverDate,
		StoreDate: maybeStoreDate,
		Series: ScrobbleIssueSeries{
			ID:        *in.Series.Id,
			Name:      in.Series.Name,
			Volume:    in.Series.Volume,
			YearBegan: in.Series.YearBegan,
		},
		Modified: *in.Modified,
	}, nil
}
