# go-metron

A Go client for the [Metron](https://metron.cloud/) comic book database.

```
go get github.com/jyggen/go-metron
```

Requires Go 1.26. Full reference on [pkg.go.dev](https://pkg.go.dev/github.com/jyggen/go-metron).

## Usage

```go
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"cloud.google.com/go/civil"
	"github.com/jyggen/go-metron"
)

func main() {
	c, err := metron.NewClient(os.Getenv("METRON_TOKEN"))
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()

	// Every Marvel issue that reached stores in the week of 2021-06-07.
	// The date bounds are inclusive.
	for issue, err := range c.Issues(ctx, &metron.IssueFilters{
		StoreDateFrom: civil.Date{Year: 2021, Month: time.June, Day: 7},
		StoreDateTo:   civil.Date{Year: 2021, Month: time.June, Day: 13},
		PublisherName: "marvel",
	}) {
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("%d %s\n", issue.ID, issue.Name)
	}

	issue, err := c.IssueByID(ctx, 31660)
	if err != nil {
		log.Fatal(err)
	}

	// Optional fields are pointers.
	if issue.Description != nil {
		fmt.Println(*issue.Description)
	}
}
```

### Pagination

List methods return an iterator and fetch each page as you reach it. Stop early and the remaining pages are never requested. Detail methods return a single record.

```go
for issue, err := range c.Issues(ctx, nil) {
	if err != nil {
		log.Fatal(err)
	}

	if issue.Name == "Amazing Fantasy" {
		break // no further pages are fetched
	}
}
```

### Filtering

Each listing takes its own filters struct, so the compiler rejects a filter the endpoint does not accept. Pass `nil` for no filters.

```go
c.Issues(ctx, &metron.IssueFilters{PublisherName: "marvel", CoverYear: 2024})
c.Arcs(ctx, &metron.ArcFilters{Name: "Crisis"})
c.Roles(ctx, nil)

c.Arcs(ctx, &metron.ArcFilters{PublisherID: 2})
// unknown field PublisherID in struct literal of type metron.ArcFilters
```

A zero-valued field is omitted from the query, so a `nil` filters pointer and an empty struct are equivalent. The two `Missing*` fields are pointers because `false` is a distinct query — records that *have* the ID — rather than an absent one.

```go
// Issues with no Comic Vine ID, so they can be matched up and filled in.
c.Issues(ctx, &metron.IssueFilters{MissingComicVineID: new(true)})
```

Date bounds are inclusive on both ends, and `civil.Date` carries no time or zone.

```go
c.Issues(ctx, &metron.IssueFilters{
	StoreDateFrom: civil.Date{Year: 2021, Month: time.June, Day: 7},
	StoreDateTo:   civil.Date{Year: 2021, Month: time.June, Day: 13},
})
```

### Retries

Nothing is retried unless you ask. `WithRetry(n)` allows up to n further attempts: a rate-limited request waits for `Retry-After`, while 502, 503, 504 and network failures back off with jitter.

Writes are exempt from the latter, since a request that failed in transit may already have been applied. They are still retried when rate-limited, because a rejected request never reached the endpoint.

```go
c, err := metron.NewClient(token, metron.WithRetry(3))
```

### Conditional requests

Detail methods take `IfModifiedSince(t)`. The API answers an unchanged record with 304, which the client reports as `ErrNotModified` and no record — so a caller holding its own copy never pays for a body it already has.

A record's own `Modified` field is the timestamp to send back, so there is nothing extra to store.

```go
issue, err := c.IssueByID(ctx, 31660, metron.IfModifiedSince(cached.Modified))

switch {
case errors.Is(err, metron.ErrNotModified):
	issue = cached // unchanged upstream
case err != nil:
	log.Fatal(err)
}

fmt.Println(issue.Name)
```

Only detail methods and the per-parent issue listings accept it — upstream wraps exactly those in a conditional response. The filterable listings discard the header, so passing it to one is a compile error; use their `ModifiedAfter` filter instead.

### Errors

`APIError` carries the status, the body and the request it came from, whose URL holds the query your filters produced and the page the iterator had reached.

```go
var apiErr *metron.APIError

if errors.As(err, &apiErr) {
	fmt.Println(apiErr.Method, apiErr.URL, apiErr.StatusCode, string(apiErr.Body))
}

if errors.Is(err, &metron.APIError{StatusCode: http.StatusNotFound}) {
	// no such issue
}
```

`MapError` is a record the client could not convert. List methods report it per record and carry on to the next one.

```go
for _, err := range c.Issues(ctx, nil) {
	if errors.Is(err, &metron.MapError{}) {
		fmt.Println(err) // metron: issue 2558: nil PriceCurrency
	}
}
```
