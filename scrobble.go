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

// ScrobbleIssue is the issue summary returned in a scrobble response.
type ScrobbleIssue struct {
	ID        int
	Number    string
	CoverDate civil.Date
	StoreDate *civil.Date
	Series    struct {
		Name      string
		Volume    int
		YearBegan int
	}
	Modified time.Time
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

// ScrobbleOption configures an optional parameter for a scrobble request.
type ScrobbleOption func(*internal.ScrobbleRequest)

// WithReadDate sets the read date for a scrobble request.
func WithReadDate(t time.Time) ScrobbleOption {
	return func(r *internal.ScrobbleRequest) {
		r.DateRead = nullable.NewNullableWithValue(t)
	}
}

// WithRating sets the rating (1-5) for a scrobble request.
func WithRating(rating int) ScrobbleOption {
	return func(r *internal.ScrobbleRequest) {
		r.Rating = nullable.NewNullableWithValue(rating)
	}
}

// Scrobble marks an issue as read, optionally setting a read date and rating.
func (c *Client) Scrobble(ctx context.Context, issueID int, opts ...ScrobbleOption) (*ScrobbleResult, error) {
	req := internal.ScrobbleRequest{
		IssueId: issueID,
	}

	for _, opt := range opts {
		opt(&req)
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
			r := int(enumVal)
			rating = &r
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
		storeDate := civil.DateOf(d.Time)
		maybeStoreDate = &storeDate
	}

	return &ScrobbleIssue{
		ID:        *in.Id,
		Number:    in.Number,
		CoverDate: coverDate,
		StoreDate: maybeStoreDate,
		Series: struct {
			Name      string
			Volume    int
			YearBegan int
		}{
			Name:      in.Series.Name,
			Volume:    in.Series.Volume,
			YearBegan: in.Series.YearBegan,
		},
		Modified: *in.Modified,
	}, nil
}
