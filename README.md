# go-metron

`go-metron` is an API client enabling Go programs to interact with [Metron](https://metron.cloud/).

## Features

### Caching

`go-metron` can optionally use [heuristic caching](https://developer.mozilla.org/en-US/docs/Web/HTTP/Caching#heuristic_caching). 

### Filtering

`go-metron` supports all filtering capabilities of Metron's API.

### Pagination

`go-metron` will automatically fetch the next page while iterating over results.

### Rate Limiting

`go-metron` will at most make 30 requests per minute, adhering to [Metron's API guidelines](https://metron.cloud/pages/guidelines/api/).

## Example Usage

```go
package main

import (
	"context"
	"fmt"

	"cloud.google.com/go/civil"
	"github.com/jyggen/go-metron"
)

func main() {
	c := metron.NewClient(
		metron.WithAuthentication("username", "password"),
		metron.WithCaching(""),
	)

	after, _ := civil.ParseDate("2021-06-07")
	before, _ := civil.ParseDate("2021-06-13")

	// Get all Marvel comics for the week of 2021-06-07
	for i, err := range c.Issues(
		context.Background(),
		metron.FilterByStoreDateRangeAfter(after),
		metron.FilterByStoreDateRangeBefore(before),
		metron.FilterByPublisherName("marvel"),
	) {
		if err != nil {
			panic(err)
		}

		fmt.Printf("%d %s\n", i.ID, i.Name)
	}

	// Retrieve the detail for an individual issue
	asm68, err := c.IssueByID(context.Background(), 31660)

	if err != nil {
		panic(err)
	}

	fmt.Println(asm68.Description)
}
```