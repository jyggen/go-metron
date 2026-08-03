package metron

import (
	"time"

	"cloud.google.com/go/civil"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

// Filter applies a query filter to a request parameter object. Filters are
// dispatched at runtime: passing one to a list method whose endpoint does not
// support it returns a FilterError on the first iteration rather than failing
// to compile.
type Filter func(v any) error

// ByAlternativeNumber filters by alternative issue number.
func ByAlternativeNumber(v string) Filter {
	return func(f any) error {
		p, ok := f.(alternativenumberFilterable)
		if !ok {
			return &FilterError{Filter: "ByAlternativeNumber"}
		}

		p.SetAltNumber(v)

		return nil
	}
}

// ByComicVineID filters by Comic Vine ID.
func ByComicVineID(v int) Filter {
	return func(f any) error {
		p, ok := f.(comicVineIDFilterable)
		if !ok {
			return &FilterError{Filter: "ByComicVineID"}
		}

		p.SetCvId(v)

		return nil
	}
}

// ByCoverHash filters by cover image hash.
func ByCoverHash(v string) Filter {
	return func(f any) error {
		p, ok := f.(coverHashFilterable)
		if !ok {
			return &FilterError{Filter: "ByCoverHash"}
		}

		p.SetCoverHash(v)

		return nil
	}
}

// ByCoverMonth filters by cover month (1-12).
func ByCoverMonth(v int) Filter {
	return func(f any) error {
		p, ok := f.(coverMonthFilterable)
		if !ok {
			return &FilterError{Filter: "ByCoverMonth"}
		}

		p.SetCoverMonth(float32(v))

		return nil
	}
}

// ByCoverYear filters by cover year.
func ByCoverYear(v int) Filter {
	return func(f any) error {
		p, ok := f.(coverYearFilterable)
		if !ok {
			return &FilterError{Filter: "ByCoverYear"}
		}

		p.SetCoverYear(float32(v))

		return nil
	}
}

// ByDesignation filters by universe designation.
func ByDesignation(v string) Filter {
	return func(f any) error {
		p, ok := f.(designationFilterable)
		if !ok {
			return &FilterError{Filter: "ByDesignation"}
		}

		p.SetDesignation(v)

		return nil
	}
}

// ByFinalOrderCutoffDate filters by final order cutoff date.
func ByFinalOrderCutoffDate(v civil.Date) Filter {
	return func(f any) error {
		p, ok := f.(finalOrderCutoffDateFilterable)
		if !ok {
			return &FilterError{Filter: "ByFinalOrderCutoffDate"}
		}

		p.SetFocDate(openapi_types.Date{Time: v.In(time.UTC)})

		return nil
	}
}

// ByFinalOrderCutoffDateRangeAfter filters to results with FOC date after v.
func ByFinalOrderCutoffDateRangeAfter(v civil.Date) Filter {
	return func(f any) error {
		p, ok := f.(finalOrderCutoffDateRangeAfterFilterable)
		if !ok {
			return &FilterError{Filter: "ByFinalOrderCutoffDateRangeAfter"}
		}

		p.SetFocDateRangeAfter(openapi_types.Date{Time: v.In(time.UTC)})

		return nil
	}
}

// ByFinalOrderCutoffDateRangeBefore filters to results with FOC date before v.
func ByFinalOrderCutoffDateRangeBefore(v civil.Date) Filter {
	return func(f any) error {
		p, ok := f.(finalOrderCutoffDateRangeBeforeFilterable)
		if !ok {
			return &FilterError{Filter: "ByFinalOrderCutoffDateRangeBefore"}
		}

		p.SetFocDateRangeBefore(openapi_types.Date{Time: v.In(time.UTC)})

		return nil
	}
}

// ByGrandComicsDatabaseID filters by Grand Comics Database ID.
func ByGrandComicsDatabaseID(v int) Filter {
	return func(f any) error {
		p, ok := f.(grandComicsDatabaseIDFilterable)
		if !ok {
			return &FilterError{Filter: "ByGrandComicsDatabaseID"}
		}

		p.SetGcdId(v)

		return nil
	}
}

// ByImprintID filters by imprint ID.
func ByImprintID(v int) Filter {
	return func(f any) error {
		p, ok := f.(imprintIDFilterable)
		if !ok {
			return &FilterError{Filter: "ByImprintID"}
		}

		p.SetImprintId(v)

		return nil
	}
}

// ByImprintName filters by imprint name.
func ByImprintName(v string) Filter {
	return func(f any) error {
		p, ok := f.(imprintnameFilterable)
		if !ok {
			return &FilterError{Filter: "ByImprintName"}
		}

		p.SetImprintName(v)

		return nil
	}
}

// ByMissingComicVineID filters to results that are missing a Comic Vine ID.
func ByMissingComicVineID() Filter {
	return func(f any) error {
		p, ok := f.(missingComicVineIDFilterable)
		if !ok {
			return &FilterError{Filter: "ByMissingComicVineID"}
		}

		p.SetMissingCvId(true)

		return nil
	}
}

// ByMissingGrandComicsDatabaseID filters to results that are missing a Grand Comics Database ID.
func ByMissingGrandComicsDatabaseID() Filter {
	return func(f any) error {
		p, ok := f.(missingGrandComicsDatabaseIDFilterable)
		if !ok {
			return &FilterError{Filter: "ByMissingGrandComicsDatabaseID"}
		}

		p.SetMissingGcdId(true)

		return nil
	}
}

// ByModifiedGreaterThan filters to results modified after v.
func ByModifiedGreaterThan(v time.Time) Filter {
	return func(f any) error {
		p, ok := f.(modifiedGreaterThanFilterable)
		if !ok {
			return &FilterError{Filter: "ByModifiedGreaterThan"}
		}

		p.SetModifiedGt(v)

		return nil
	}
}

// ByName filters by name.
func ByName(v string) Filter {
	return func(f any) error {
		p, ok := f.(nameFilterable)
		if !ok {
			return &FilterError{Filter: "ByName"}
		}

		p.SetName(v)

		return nil
	}
}

// ByNumber filters by issue number.
func ByNumber(v string) Filter {
	return func(f any) error {
		p, ok := f.(numberFilterable)
		if !ok {
			return &FilterError{Filter: "ByNumber"}
		}

		p.SetNumber(v)

		return nil
	}
}

// ByPublisherID filters by publisher ID.
func ByPublisherID(v int) Filter {
	return func(f any) error {
		p, ok := f.(publisherIDFilterable)
		if !ok {
			return &FilterError{Filter: "ByPublisherID"}
		}

		p.SetPublisherId(v)

		return nil
	}
}

// ByPublisherName filters by publisher name.
func ByPublisherName(v string) Filter {
	return func(f any) error {
		p, ok := f.(publishernameFilterable)
		if !ok {
			return &FilterError{Filter: "ByPublisherName"}
		}

		p.SetPublisherName(v)

		return nil
	}
}

// ByRating filters by content rating.
func ByRating(v string) Filter {
	return func(f any) error {
		p, ok := f.(ratingFilterable)
		if !ok {
			return &FilterError{Filter: "ByRating"}
		}

		p.SetRating(v)

		return nil
	}
}

// BySeriesID filters by series ID.
func BySeriesID(v int) Filter {
	return func(f any) error {
		p, ok := f.(seriesIDFilterable)
		if !ok {
			return &FilterError{Filter: "BySeriesID"}
		}

		p.SetSeriesId(v)

		return nil
	}
}

// BySeriesName filters by series name.
func BySeriesName(v string) Filter {
	return func(f any) error {
		p, ok := f.(seriesnameFilterable)
		if !ok {
			return &FilterError{Filter: "BySeriesName"}
		}

		p.SetSeriesName(v)

		return nil
	}
}

// BySeriesType filters by series type name.
func BySeriesType(v string) Filter {
	return func(f any) error {
		p, ok := f.(seriesTypeFilterable)
		if !ok {
			return &FilterError{Filter: "BySeriesType"}
		}

		p.SetSeriesType(v)

		return nil
	}
}

// BySeriesTypeID filters by series type ID.
func BySeriesTypeID(v int) Filter {
	return func(f any) error {
		p, ok := f.(seriesTypeIDFilterable)
		if !ok {
			return &FilterError{Filter: "BySeriesTypeID"}
		}

		p.SetSeriesTypeId(v)

		return nil
	}
}

// BySeriesVolume filters by series volume number.
func BySeriesVolume(v int) Filter {
	return func(f any) error {
		p, ok := f.(seriesvolumeFilterable)
		if !ok {
			return &FilterError{Filter: "BySeriesVolume"}
		}

		p.SetSeriesVolume(v)

		return nil
	}
}

// BySeriesYearBegan filters by the year a series began.
func BySeriesYearBegan(v int) Filter {
	return func(f any) error {
		p, ok := f.(seriesYearBeganFilterable)
		if !ok {
			return &FilterError{Filter: "BySeriesYearBegan"}
		}

		p.SetSeriesYearBegan(v)

		return nil
	}
}

// ByStatus filters by series status.
func ByStatus(v int) Filter {
	return func(f any) error {
		p, ok := f.(statusFilterable)
		if !ok {
			return &FilterError{Filter: "ByStatus"}
		}

		p.SetStatus(v)

		return nil
	}
}

// ByStockKeepingUnit filters by SKU.
func ByStockKeepingUnit(v string) Filter {
	return func(f any) error {
		p, ok := f.(stockKeepingUnitFilterable)
		if !ok {
			return &FilterError{Filter: "ByStockKeepingUnit"}
		}

		p.SetSku(v)

		return nil
	}
}

// ByStoreDate filters by in-store date.
func ByStoreDate(v civil.Date) Filter {
	return func(f any) error {
		p, ok := f.(storeDateFilterable)
		if !ok {
			return &FilterError{Filter: "ByStoreDate"}
		}

		p.SetStoreDate(openapi_types.Date{Time: v.In(time.UTC)})

		return nil
	}
}

// ByStoreDateRangeAfter filters to results with store date after v.
func ByStoreDateRangeAfter(v civil.Date) Filter {
	return func(f any) error {
		p, ok := f.(storeDateRangeAfterFilterable)
		if !ok {
			return &FilterError{Filter: "ByStoreDateRangeAfter"}
		}

		p.SetStoreDateRangeAfter(openapi_types.Date{Time: v.In(time.UTC)})

		return nil
	}
}

// ByStoreDateRangeBefore filters to results with store date before v.
func ByStoreDateRangeBefore(v civil.Date) Filter {
	return func(f any) error {
		p, ok := f.(storeDateRangeBeforeFilterable)
		if !ok {
			return &FilterError{Filter: "ByStoreDateRangeBefore"}
		}

		p.SetStoreDateRangeBefore(openapi_types.Date{Time: v.In(time.UTC)})

		return nil
	}
}

// ByUniversalProductCode filters by UPC.
func ByUniversalProductCode(v string) Filter {
	return func(f any) error {
		p, ok := f.(universalProductCodeFilterable)
		if !ok {
			return &FilterError{Filter: "ByUniversalProductCode"}
		}

		p.SetUpc(v)

		return nil
	}
}

// ByVolume filters by volume number.
func ByVolume(v int) Filter {
	return func(f any) error {
		p, ok := f.(volumeFilterable)
		if !ok {
			return &FilterError{Filter: "ByVolume"}
		}

		p.SetVolume(v)

		return nil
	}
}

// ByYearBegan filters by the year a series began.
func ByYearBegan(v int) Filter {
	return func(f any) error {
		p, ok := f.(yearBeganFilterable)
		if !ok {
			return &FilterError{Filter: "ByYearBegan"}
		}

		p.SetYearBegan(v)

		return nil
	}
}

// ByYearEnd filters by the year a series ended.
func ByYearEnd(v int) Filter {
	return func(f any) error {
		p, ok := f.(yearEndFilterable)
		if !ok {
			return &FilterError{Filter: "ByYearEnd"}
		}

		p.SetYearEnd(v)

		return nil
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
