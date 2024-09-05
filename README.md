# go-metron

An API client allowing Go programs to interact with Metron.                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                          

## Example Usage

```go
package main

import (
	"cloud.google.com/go/civil"
	"context"
	"fmt"
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

```
