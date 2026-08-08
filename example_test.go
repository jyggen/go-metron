package metron_test

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"cloud.google.com/go/civil"
	"github.com/jyggen/go-metron"
)

// These examples have no Output comment, so `go test` compiles them without
// running them. That keeps them from drifting out of sync with the API.

func Example() {
	c, err := metron.NewClient(os.Getenv("METRON_TOKEN"))
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

// ExampleWithClient tunes the underlying transport. The defaults are fine for
// occasional calls, but an iterator walking thousands of pages benefits from a
// larger connection pool and an explicit timeout.
func ExampleWithClient() {
	transport := http.DefaultTransport.(*http.Transport).Clone()
	transport.MaxIdleConnsPerHost = 20

	c, err := metron.NewClient("your-api-token", metron.WithClient(&http.Client{
		Timeout:   30 * time.Second,
		Transport: transport,
	}))
	if err != nil {
		log.Fatal(err)
	}

	_ = c
}
