package metron

import (
	"time"

	"cloud.google.com/go/civil"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

// Filter applies a query filter to a request parameter object. Filters are
// dispatched at runtime: an unsupported filter returns a FilterError.
type Filter func(v any) error

// newFilter builds a Filter for any params implementing I, inferred from apply.
// name is explicit because runtime.Caller would not survive inlining.
func newFilter[I any](name string, apply func(I)) Filter {
	return func(f any) error {
		p, ok := f.(I)
		if !ok {
			return &FilterError{Filter: name}
		}

		apply(p)

		return nil
	}
}

// ByAlternativeNumber filters by alternative issue number.
func ByAlternativeNumber(v string) Filter {
	return newFilter("ByAlternativeNumber", func(p interface{ SetAltNumber(string) }) {
		p.SetAltNumber(v)
	})
}

// ByComicVineID filters by Comic Vine ID.
func ByComicVineID(v int) Filter {
	return newFilter("ByComicVineID", func(p interface{ SetCvId(int) }) {
		p.SetCvId(v)
	})
}

// ByCoverHash filters by cover image hash.
func ByCoverHash(v string) Filter {
	return newFilter("ByCoverHash", func(p interface{ SetCoverHash(string) }) {
		p.SetCoverHash(v)
	})
}

// ByCoverMonth filters by cover month (1-12).
func ByCoverMonth(v int) Filter {
	return newFilter("ByCoverMonth", func(p interface{ SetCoverMonth(float32) }) {
		p.SetCoverMonth(float32(v))
	})
}

// ByCoverYear filters by cover year.
func ByCoverYear(v int) Filter {
	return newFilter("ByCoverYear", func(p interface{ SetCoverYear(float32) }) {
		p.SetCoverYear(float32(v))
	})
}

// ByDesignation filters by universe designation.
func ByDesignation(v string) Filter {
	return newFilter("ByDesignation", func(p interface{ SetDesignation(string) }) {
		p.SetDesignation(v)
	})
}

// ByFinalOrderCutoffDate filters by final order cutoff date.
func ByFinalOrderCutoffDate(v civil.Date) Filter {
	return newFilter("ByFinalOrderCutoffDate", func(p interface{ SetFocDate(openapi_types.Date) }) {
		p.SetFocDate(openapi_types.Date{Time: v.In(time.UTC)})
	})
}

// ByFinalOrderCutoffDateRangeAfter filters to results with FOC date after v.
func ByFinalOrderCutoffDateRangeAfter(v civil.Date) Filter {
	return newFilter("ByFinalOrderCutoffDateRangeAfter", func(p interface{ SetFocDateRangeAfter(openapi_types.Date) }) {
		p.SetFocDateRangeAfter(openapi_types.Date{Time: v.In(time.UTC)})
	})
}

// ByFinalOrderCutoffDateRangeBefore filters to results with FOC date before v.
func ByFinalOrderCutoffDateRangeBefore(v civil.Date) Filter {
	return newFilter("ByFinalOrderCutoffDateRangeBefore", func(p interface{ SetFocDateRangeBefore(openapi_types.Date) }) {
		p.SetFocDateRangeBefore(openapi_types.Date{Time: v.In(time.UTC)})
	})
}

// ByGrandComicsDatabaseID filters by Grand Comics Database ID.
func ByGrandComicsDatabaseID(v int) Filter {
	return newFilter("ByGrandComicsDatabaseID", func(p interface{ SetGcdId(int) }) {
		p.SetGcdId(v)
	})
}

// ByImprintID filters by imprint ID.
func ByImprintID(v int) Filter {
	return newFilter("ByImprintID", func(p interface{ SetImprintId(int) }) {
		p.SetImprintId(v)
	})
}

// ByImprintName filters by imprint name.
func ByImprintName(v string) Filter {
	return newFilter("ByImprintName", func(p interface{ SetImprintName(string) }) {
		p.SetImprintName(v)
	})
}

// ByMissingComicVineID filters to results that are missing a Comic Vine ID.
func ByMissingComicVineID() Filter {
	return newFilter("ByMissingComicVineID", func(p interface{ SetMissingCvId(bool) }) {
		p.SetMissingCvId(true)
	})
}

// ByMissingGrandComicsDatabaseID filters to results that are missing a Grand Comics Database ID.
func ByMissingGrandComicsDatabaseID() Filter {
	return newFilter("ByMissingGrandComicsDatabaseID", func(p interface{ SetMissingGcdId(bool) }) {
		p.SetMissingGcdId(true)
	})
}

// ByModifiedGreaterThan filters to results modified after v.
func ByModifiedGreaterThan(v time.Time) Filter {
	return newFilter("ByModifiedGreaterThan", func(p interface{ SetModifiedGt(time.Time) }) {
		p.SetModifiedGt(v)
	})
}

// ByName filters by name.
func ByName(v string) Filter {
	return newFilter("ByName", func(p interface{ SetName(string) }) {
		p.SetName(v)
	})
}

// ByNumber filters by issue number.
func ByNumber(v string) Filter {
	return newFilter("ByNumber", func(p interface{ SetNumber(string) }) {
		p.SetNumber(v)
	})
}

// ByPublisherID filters by publisher ID.
func ByPublisherID(v int) Filter {
	return newFilter("ByPublisherID", func(p interface{ SetPublisherId(int) }) {
		p.SetPublisherId(v)
	})
}

// ByPublisherName filters by publisher name.
func ByPublisherName(v string) Filter {
	return newFilter("ByPublisherName", func(p interface{ SetPublisherName(string) }) {
		p.SetPublisherName(v)
	})
}

// ByRating filters by content rating.
func ByRating(v string) Filter {
	return newFilter("ByRating", func(p interface{ SetRating(string) }) {
		p.SetRating(v)
	})
}

// BySeriesID filters by series ID.
func BySeriesID(v int) Filter {
	return newFilter("BySeriesID", func(p interface{ SetSeriesId(int) }) {
		p.SetSeriesId(v)
	})
}

// BySeriesName filters by series name.
func BySeriesName(v string) Filter {
	return newFilter("BySeriesName", func(p interface{ SetSeriesName(string) }) {
		p.SetSeriesName(v)
	})
}

// BySeriesType filters by series type name.
func BySeriesType(v string) Filter {
	return newFilter("BySeriesType", func(p interface{ SetSeriesType(string) }) {
		p.SetSeriesType(v)
	})
}

// BySeriesTypeID filters by series type ID.
func BySeriesTypeID(v int) Filter {
	return newFilter("BySeriesTypeID", func(p interface{ SetSeriesTypeId(int) }) {
		p.SetSeriesTypeId(v)
	})
}

// BySeriesVolume filters by series volume number.
func BySeriesVolume(v int) Filter {
	return newFilter("BySeriesVolume", func(p interface{ SetSeriesVolume(int) }) {
		p.SetSeriesVolume(v)
	})
}

// BySeriesYearBegan filters by the year a series began.
func BySeriesYearBegan(v int) Filter {
	return newFilter("BySeriesYearBegan", func(p interface{ SetSeriesYearBegan(int) }) {
		p.SetSeriesYearBegan(v)
	})
}

// ByStatus filters by series status.
func ByStatus(v int) Filter {
	return newFilter("ByStatus", func(p interface{ SetStatus(int) }) {
		p.SetStatus(v)
	})
}

// ByStockKeepingUnit filters by SKU.
func ByStockKeepingUnit(v string) Filter {
	return newFilter("ByStockKeepingUnit", func(p interface{ SetSku(string) }) {
		p.SetSku(v)
	})
}

// ByStoreDate filters by in-store date.
func ByStoreDate(v civil.Date) Filter {
	return newFilter("ByStoreDate", func(p interface{ SetStoreDate(openapi_types.Date) }) {
		p.SetStoreDate(openapi_types.Date{Time: v.In(time.UTC)})
	})
}

// ByStoreDateRangeAfter filters to results with store date after v.
func ByStoreDateRangeAfter(v civil.Date) Filter {
	return newFilter("ByStoreDateRangeAfter", func(p interface{ SetStoreDateRangeAfter(openapi_types.Date) }) {
		p.SetStoreDateRangeAfter(openapi_types.Date{Time: v.In(time.UTC)})
	})
}

// ByStoreDateRangeBefore filters to results with store date before v.
func ByStoreDateRangeBefore(v civil.Date) Filter {
	return newFilter("ByStoreDateRangeBefore", func(p interface{ SetStoreDateRangeBefore(openapi_types.Date) }) {
		p.SetStoreDateRangeBefore(openapi_types.Date{Time: v.In(time.UTC)})
	})
}

// ByUniversalProductCode filters by UPC.
func ByUniversalProductCode(v string) Filter {
	return newFilter("ByUniversalProductCode", func(p interface{ SetUpc(string) }) {
		p.SetUpc(v)
	})
}

// ByVolume filters by volume number.
func ByVolume(v int) Filter {
	return newFilter("ByVolume", func(p interface{ SetVolume(int) }) {
		p.SetVolume(v)
	})
}

// ByYearBegan filters by the year a series began.
func ByYearBegan(v int) Filter {
	return newFilter("ByYearBegan", func(p interface{ SetYearBegan(int) }) {
		p.SetYearBegan(v)
	})
}

// ByYearEnd filters by the year a series ended.
func ByYearEnd(v int) Filter {
	return newFilter("ByYearEnd", func(p interface{ SetYearEnd(int) }) {
		p.SetYearEnd(v)
	})
}
