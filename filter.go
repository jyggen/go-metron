package metron

import (
	"cloud.google.com/go/civil"
	"fmt"
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

func FilterByCoverHash(coverHash string) Filter {
	return func(q *url.Values) {
		q.Set("cover_hash", coverHash)
	}
}

func FilterByCoverMonth(month time.Month) Filter {
	return func(q *url.Values) {
		q.Set("cover_month", strconv.Itoa(int(month)))
	}
}

func FilterByCoverYear(year int) Filter {
	return func(q *url.Values) {
		q.Set("cover_year", strconv.Itoa(year))
	}
}

func FilterByDesignation(designation string) Filter {
	return func(q *url.Values) {
		q.Set("designation", designation)
	}
}

func FilterByImprintID(id int) Filter {
	return func(q *url.Values) {
		q.Set("imprint_id", strconv.Itoa(id))
	}
}

func FilterByImprintName(name string) Filter {
	return func(q *url.Values) {
		q.Set("imprint_name", name)
	}
}

func FilterByMissingComicVineID(isMissing bool) Filter {
	return func(q *url.Values) {
		q.Set("missing_cv_id", fmt.Sprintf("%t", isMissing))
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

func FilterByNumber(number string) Filter {
	return func(q *url.Values) {
		q.Set("number", number)
	}
}

func FilterByPublisherID(id int) Filter {
	return func(q *url.Values) {
		q.Set("publisher_id", strconv.Itoa(id))
	}
}

func FilterByPublisherName(name string) Filter {
	return func(q *url.Values) {
		q.Set("publisher_name", name)
	}
}

func FilterByRating(rating string) Filter {
	return func(q *url.Values) {
		q.Set("rating", rating)
	}
}

func FilterBySeriesID(id int) Filter {
	return func(q *url.Values) {
		q.Set("series_id", strconv.Itoa(id))
	}
}

func FilterBySeriesName(name string) Filter {
	return func(q *url.Values) {
		q.Set("series_name", name)
	}
}

func FilterBySeriesType(name string) Filter {
	return func(q *url.Values) {
		q.Set("series_type", name)
	}
}

func FilterBySeriesTypeID(id int) Filter {
	return func(q *url.Values) {
		q.Set("series_type_id", strconv.Itoa(id))
	}
}

func FilterBySeriesVolume(volume int) Filter {
	return func(q *url.Values) {
		q.Set("series_volume", strconv.Itoa(volume))
	}
}

func FilterBySeriesYearBegan(year int) Filter {
	return func(q *url.Values) {
		q.Set("series_year_began", strconv.Itoa(year))
	}
}

func FilterBySKU(sku string) Filter {
	return func(q *url.Values) {
		q.Set("sku", sku)
	}
}

func FilterByStatus(status int) Filter {
	return func(q *url.Values) {
		q.Set("status", strconv.Itoa(status))
	}
}

func FilterByStoreDate(date civil.Date) Filter {
	return func(q *url.Values) {
		q.Set("store_date", date.String())
	}
}

func FilterByStoreDateRangeAfter(date civil.Date) Filter {
	return func(q *url.Values) {
		q.Set("store_date_range_after", date.String())
	}
}

func FilterByStoreDateRangeBefore(date civil.Date) Filter {
	return func(q *url.Values) {
		q.Set("store_date_range_before", date.String())
	}
}

func FilterByUPC(upc string) Filter {
	return func(q *url.Values) {
		q.Set("upc", upc)
	}
}

func FilterByVolume(volume int) Filter {
	return func(q *url.Values) {
		q.Set("volume", strconv.Itoa(volume))
	}
}

func FilterByYearBegan(year int) Filter {
	return func(q *url.Values) {
		q.Set("year_began", strconv.Itoa(year))
	}
}

func FilterByYearEnded(year int) Filter {
	return func(q *url.Values) {
		q.Set("year_end", strconv.Itoa(year))
	}
}
