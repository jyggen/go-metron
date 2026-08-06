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
	c, err := metron.NewClient(os.Getenv("METRON_TOKEN"), metron.WithCaching())
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()

	// Every Marvel issue that reached stores in the week of 2021-06-07.
	for issue, err := range c.Issues(ctx,
		metron.ByStoreDateRangeAfter(civil.Date{Year: 2021, Month: time.June, Day: 7}),
		metron.ByStoreDateRangeBefore(civil.Date{Year: 2021, Month: time.June, Day: 13}),
		metron.ByPublisherName("marvel"),
	) {
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
for issue, err := range c.Issues(ctx) {
	if err != nil {
		log.Fatal(err)
	}

	if issue.Name == "Amazing Fantasy" {
		break // no further pages are fetched
	}
}
```

### Filtering

Filters are applied at runtime, not checked by the compiler. A filter the endpoint does not accept returns a `FilterError` on the first iteration instead of being ignored, naming both the filter and the method that rejected it.

```go
for _, err := range c.Roles(ctx, metron.ByPublisherID(2)) {
	fmt.Println(err) // metron: ByPublisherID does not apply to Roles
}
```

### Caching

```go
c, err := metron.NewClient(token, metron.WithCaching())
```

Responses go under your user cache directory and are revalidated with `If-Modified-Since`. A response without a `Last-Modified` header is not cached, since there is nothing to derive a freshness window from.

### Retries

Nothing is retried unless you ask. `WithRetry(n)` allows up to n further attempts: a rate-limited request waits for `Retry-After`, while 502, 503, 504 and network failures back off with jitter.

Writes are exempt from the latter, since a request that failed in transit may already have been applied. They are still retried when rate-limited, because a rejected request never reached the endpoint.

```go
c, err := metron.NewClient(token, metron.WithRetry(3))
```

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
for _, err := range c.Issues(ctx) {
	if errors.Is(err, &metron.MapError{}) {
		fmt.Println(err) // metron: issue 2558: nil PriceCurrency
	}
}
```
