package metron

import (
	"net/url"
	"strconv"
	"time"
)

type Filter func(q *url.Values)

func FilterByComicVineID(comicVineID int) Filter {
	return func(q *url.Values) {
		q.Set("cv_id", strconv.Itoa(comicVineID))
	}
}

func FilterByModifiedGreaterThan(modified time.Time) Filter {
	return func(q *url.Values) {
		q.Set("modified_gt", modified.Format(time.RFC3339))
	}
}

func FilterByName(name string) Filter {
	return func(q *url.Values) {
		q.Set("name", name)
	}
}
