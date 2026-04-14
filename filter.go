package metron

import (
	"time"

	"cloud.google.com/go/civil"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

// Filter is a function that applies a query filter to a request parameter object.
type Filter func(v any)

// ByAlternativeNumber filters by alternative issue number.
func ByAlternativeNumber(v string) Filter {
	return func(f any) {
		if p, ok := f.(alternativenumberFilterable); ok {
			p.SetAltNumber(v)
		}
	}
}

// ByComicVineID filters by Comic Vine ID.
func ByComicVineID(v int) Filter {
	return func(f any) {
		if p, ok := f.(comicVineIDFilterable); ok {
			p.SetCvId(v)
		}
	}
}

// ByCoverHash filters by cover image hash.
func ByCoverHash(v string) Filter {
	return func(f any) {
		if p, ok := f.(coverHashFilterable); ok {
			p.SetCoverHash(v)
		}
	}
}

// ByCoverMonth filters by cover month (1-12).
func ByCoverMonth(v int) Filter {
	return func(f any) {
		if p, ok := f.(coverMonthFilterable); ok {
			p.SetCoverMonth(float32(v))
		}
	}
}

// ByCoverYear filters by cover year.
func ByCoverYear(v int) Filter {
	return func(f any) {
		if p, ok := f.(coverYearFilterable); ok {
			p.SetCoverYear(float32(v))
		}
	}
}

// ByDesignation filters by universe designation.
func ByDesignation(v string) Filter {
	return func(f any) {
		if p, ok := f.(designationFilterable); ok {
			p.SetDesignation(v)
		}
	}
}

// ByFinalOrderCutoffDate filters by final order cutoff date.
func ByFinalOrderCutoffDate(v civil.Date) Filter {
	return func(f any) {
		if p, ok := f.(finalOrderCutoffDateFilterable); ok {
			p.SetFocDate(openapi_types.Date{Time: v.In(time.UTC)})
		}
	}
}

// ByFinalOrderCutoffDateRangeAfter filters to results with FOC date after v.
func ByFinalOrderCutoffDateRangeAfter(v civil.Date) Filter {
	return func(f any) {
		if p, ok := f.(finalOrderCutoffDateRangeAfterFilterable); ok {
			p.SetFocDateRangeAfter(openapi_types.Date{Time: v.In(time.UTC)})
		}
	}
}

// ByFinalOrderCutoffDateRangeBefore filters to results with FOC date before v.
func ByFinalOrderCutoffDateRangeBefore(v civil.Date) Filter {
	return func(f any) {
		if p, ok := f.(finalOrderCutoffDateRangeBeforeFilterable); ok {
			p.SetFocDateRangeBefore(openapi_types.Date{Time: v.In(time.UTC)})
		}
	}
}

// ByGrandComicsDatabaseID filters by Grand Comics Database ID.
func ByGrandComicsDatabaseID(v int) Filter {
	return func(f any) {
		if p, ok := f.(grandComicsDatabaseIDFilterable); ok {
			p.SetGcdId(v)
		}
	}
}

// ByImprintID filters by imprint ID.
func ByImprintID(v int) Filter {
	return func(f any) {
		if p, ok := f.(imprintIDFilterable); ok {
			p.SetImprintId(v)
		}
	}
}

// ByImprintName filters by imprint name.
func ByImprintName(v string) Filter {
	return func(f any) {
		if p, ok := f.(imprintnameFilterable); ok {
			p.SetImprintName(v)
		}
	}
}

// ByMissingComicVineID filters to results that are missing a Comic Vine ID.
func ByMissingComicVineID() Filter {
	return func(f any) {
		if p, ok := f.(missingComicVineIDFilterable); ok {
			p.SetMissingCvId(true)
		}
	}
}

// ByMissingGrandComicsDatabaseID filters to results that are missing a Grand Comics Database ID.
func ByMissingGrandComicsDatabaseID() Filter {
	return func(f any) {
		if p, ok := f.(missingGrandComicsDatabaseIDFilterable); ok {
			p.SetMissingGcdId(true)
		}
	}
}

// ByModifiedGreaterThan filters to results modified after v.
func ByModifiedGreaterThan(v time.Time) Filter {
	return func(f any) {
		if p, ok := f.(modifiedGreaterThanFilterable); ok {
			p.SetModifiedGt(v)
		}
	}
}

// ByName filters by name.
func ByName(v string) Filter {
	return func(f any) {
		if p, ok := f.(nameFilterable); ok {
			p.SetName(v)
		}
	}
}

// ByNumber filters by issue number.
func ByNumber(v string) Filter {
	return func(f any) {
		if p, ok := f.(numberFilterable); ok {
			p.SetNumber(v)
		}
	}
}

// ByPublisherID filters by publisher ID.
func ByPublisherID(v int) Filter {
	return func(f any) {
		if p, ok := f.(publisherIDFilterable); ok {
			p.SetPublisherId(v)
		}
	}
}

// ByPublisherName filters by publisher name.
func ByPublisherName(v string) Filter {
	return func(f any) {
		if p, ok := f.(publishernameFilterable); ok {
			p.SetPublisherName(v)
		}
	}
}

// ByRating filters by content rating.
func ByRating(v string) Filter {
	return func(f any) {
		if p, ok := f.(ratingFilterable); ok {
			p.SetRating(v)
		}
	}
}

// BySeriesID filters by series ID.
func BySeriesID(v int) Filter {
	return func(f any) {
		if p, ok := f.(seriesIDFilterable); ok {
			p.SetSeriesId(v)
		}
	}
}

// BySeriesName filters by series name.
func BySeriesName(v string) Filter {
	return func(f any) {
		if p, ok := f.(seriesnameFilterable); ok {
			p.SetSeriesName(v)
		}
	}
}

// BySeriesType filters by series type name.
func BySeriesType(v string) Filter {
	return func(f any) {
		if p, ok := f.(seriesTypeFilterable); ok {
			p.SetSeriesType(v)
		}
	}
}

// BySeriesTypeID filters by series type ID.
func BySeriesTypeID(v int) Filter {
	return func(f any) {
		if p, ok := f.(seriesTypeIDFilterable); ok {
			p.SetSeriesTypeId(v)
		}
	}
}

// BySeriesVolume filters by series volume number.
func BySeriesVolume(v int) Filter {
	return func(f any) {
		if p, ok := f.(seriesvolumeFilterable); ok {
			p.SetSeriesVolume(v)
		}
	}
}

// BySeriesYearBegan filters by the year a series began.
func BySeriesYearBegan(v int) Filter {
	return func(f any) {
		if p, ok := f.(seriesYearBeganFilterable); ok {
			p.SetSeriesYearBegan(v)
		}
	}
}

// ByStatus filters by series status.
func ByStatus(v int) Filter {
	return func(f any) {
		if p, ok := f.(statusFilterable); ok {
			p.SetStatus(v)
		}
	}
}

// ByStockKeepingUnit filters by SKU.
func ByStockKeepingUnit(v string) Filter {
	return func(f any) {
		if p, ok := f.(stockKeepingUnitFilterable); ok {
			p.SetSku(v)
		}
	}
}

// ByStoreDate filters by in-store date.
func ByStoreDate(v civil.Date) Filter {
	return func(f any) {
		if p, ok := f.(storeDateFilterable); ok {
			p.SetStoreDate(openapi_types.Date{Time: v.In(time.UTC)})
		}
	}
}

// ByStoreDateRangeAfter filters to results with store date after v.
func ByStoreDateRangeAfter(v civil.Date) Filter {
	return func(f any) {
		if p, ok := f.(storeDateRangeAfterFilterable); ok {
			p.SetStoreDateRangeAfter(openapi_types.Date{Time: v.In(time.UTC)})
		}
	}
}

// ByStoreDateRangeBefore filters to results with store date before v.
func ByStoreDateRangeBefore(v civil.Date) Filter {
	return func(f any) {
		if p, ok := f.(storeDateRangeBeforeFilterable); ok {
			p.SetStoreDateRangeBefore(openapi_types.Date{Time: v.In(time.UTC)})
		}
	}
}

// ByUniversalProductCode filters by UPC.
func ByUniversalProductCode(v string) Filter {
	return func(f any) {
		if p, ok := f.(universalProductCodeFilterable); ok {
			p.SetUpc(v)
		}
	}
}

// ByVolume filters by volume number.
func ByVolume(v int) Filter {
	return func(f any) {
		if p, ok := f.(volumeFilterable); ok {
			p.SetVolume(v)
		}
	}
}

// ByYearBegan filters by the year a series began.
func ByYearBegan(v int) Filter {
	return func(f any) {
		if p, ok := f.(yearBeganFilterable); ok {
			p.SetYearBegan(v)
		}
	}
}

// ByYearEnd filters by the year a series ended.
func ByYearEnd(v int) Filter {
	return func(f any) {
		if p, ok := f.(yearEndFilterable); ok {
			p.SetYearEnd(v)
		}
	}
}

type alternativenumberFilterable interface {
	SetAltNumber(string)
}

type comicVineIDFilterable interface {
	SetCvId(int)
}

type coverHashFilterable interface {
	SetCoverHash(string)
}

type coverMonthFilterable interface {
	SetCoverMonth(float32)
}

type coverYearFilterable interface {
	SetCoverYear(float32)
}

type designationFilterable interface {
	SetDesignation(string)
}

type finalOrderCutoffDateFilterable interface {
	SetFocDate(openapi_types.Date)
}

type finalOrderCutoffDateRangeAfterFilterable interface {
	SetFocDateRangeAfter(openapi_types.Date)
}

type finalOrderCutoffDateRangeBeforeFilterable interface {
	SetFocDateRangeBefore(openapi_types.Date)
}

type grandComicsDatabaseIDFilterable interface {
	SetGcdId(int)
}

type imprintIDFilterable interface {
	SetImprintId(int)
}

type imprintnameFilterable interface {
	SetImprintName(string)
}

type missingComicVineIDFilterable interface {
	SetMissingCvId(bool)
}

type missingGrandComicsDatabaseIDFilterable interface {
	SetMissingGcdId(bool)
}

type modifiedGreaterThanFilterable interface {
	SetModifiedGt(time.Time)
}

type nameFilterable interface {
	SetName(string)
}

type numberFilterable interface {
	SetNumber(string)
}

type publisherIDFilterable interface {
	SetPublisherId(int)
}

type publishernameFilterable interface {
	SetPublisherName(string)
}

type ratingFilterable interface {
	SetRating(string)
}

type seriesIDFilterable interface {
	SetSeriesId(int)
}

type seriesnameFilterable interface {
	SetSeriesName(string)
}

type seriesTypeFilterable interface {
	SetSeriesType(string)
}

type seriesTypeIDFilterable interface {
	SetSeriesTypeId(int)
}

type seriesvolumeFilterable interface {
	SetSeriesVolume(int)
}

type seriesYearBeganFilterable interface {
	SetSeriesYearBegan(int)
}

type statusFilterable interface {
	SetStatus(int)
}

type stockKeepingUnitFilterable interface {
	SetSku(string)
}

type storeDateFilterable interface {
	SetStoreDate(openapi_types.Date)
}

type storeDateRangeAfterFilterable interface {
	SetStoreDateRangeAfter(openapi_types.Date)
}

type storeDateRangeBeforeFilterable interface {
	SetStoreDateRangeBefore(openapi_types.Date)
}

type universalProductCodeFilterable interface {
	SetUpc(string)
}

type volumeFilterable interface {
	SetVolume(int)
}

type yearBeganFilterable interface {
	SetYearBegan(int)
}

type yearEndFilterable interface {
	SetYearEnd(int)
}
