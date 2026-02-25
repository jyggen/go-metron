package metron

import (
	"time"

	"cloud.google.com/go/civil"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

// Filter is a function that applies a query filter to a request parameter object.
type Filter func(v any)

func ByAlternativeNumber(v string) Filter {
	return func(f any) {
		if p, ok := f.(AlternativeNumberFilterable); ok {
			p.SetAltNumber(v)
		}
	}
}

func ByComicVineID(v int) Filter {
	return func(f any) {
		if p, ok := f.(ComicVineIDFilterable); ok {
			p.SetCvId(v)
		}
	}
}

func ByCoverHash(v string) Filter {
	return func(f any) {
		if p, ok := f.(CoverHashFilterable); ok {
			p.SetCoverHash(v)
		}
	}
}

func ByCoverMonth(v int) Filter {
	return func(f any) {
		if p, ok := f.(CoverMonthFilterable); ok {
			p.SetCoverMonth(float32(v))
		}
	}
}

func ByCoverYear(v int) Filter {
	return func(f any) {
		if p, ok := f.(CoverYearFilterable); ok {
			p.SetCoverYear(float32(v))
		}
	}
}

func ByDesignation(v string) Filter {
	return func(f any) {
		if p, ok := f.(DesignationFilterable); ok {
			p.SetDesignation(v)
		}
	}
}

func ByFinalOrderCutoffDate(v civil.Date) Filter {
	return func(f any) {
		if p, ok := f.(FinalOrderCutoffDateFilterable); ok {
			p.SetFocDate(openapi_types.Date{Time: v.In(time.UTC)})
		}
	}
}

func ByFinalOrderCutoffDateRangeAfter(v civil.Date) Filter {
	return func(f any) {
		if p, ok := f.(FinalOrderCutoffDateRangeAfterFilterable); ok {
			p.SetFocDateRangeAfter(openapi_types.Date{Time: v.In(time.UTC)})
		}
	}
}

func ByFinalOrderCutoffDateRangeBefore(v civil.Date) Filter {
	return func(f any) {
		if p, ok := f.(FinalOrderCutoffDateRangeBeforeFilterable); ok {
			p.SetFocDateRangeBefore(openapi_types.Date{Time: v.In(time.UTC)})
		}
	}
}

func ByGrandComicsDatabaseID(v int) Filter {
	return func(f any) {
		if p, ok := f.(GrandComicsDatabaseIDFilterable); ok {
			p.SetGcdId(v)
		}
	}
}

func ByImprintID(v int) Filter {
	return func(f any) {
		if p, ok := f.(ImprintIDFilterable); ok {
			p.SetImprintId(v)
		}
	}
}

func ByImprintName(v string) Filter {
	return func(f any) {
		if p, ok := f.(ImprintNameFilterable); ok {
			p.SetImprintName(v)
		}
	}
}

func ByMissingComicVineID(v bool) Filter {
	return func(f any) {
		if p, ok := f.(MissingComicVineIDFilterable); ok {
			p.SetMissingCvId(v)
		}
	}
}

func ByMissingGrandComicsDatabaseID(v bool) Filter {
	return func(f any) {
		if p, ok := f.(MissingGrandComicsDatabaseIDFilterable); ok {
			p.SetMissingGcdId(v)
		}
	}
}

func ByModifiedGreaterThan(v time.Time) Filter {
	return func(f any) {
		if p, ok := f.(ModifiedGreaterThanFilterable); ok {
			p.SetModifiedGt(v)
		}
	}
}

func ByName(v string) Filter {
	return func(f any) {
		if p, ok := f.(NameFilterable); ok {
			p.SetName(v)
		}
	}
}

func ByNumber(v string) Filter {
	return func(f any) {
		if p, ok := f.(NumberFilterable); ok {
			p.SetNumber(v)
		}
	}
}

func ByPublisherID(v int) Filter {
	return func(f any) {
		if p, ok := f.(PublisherIDFilterable); ok {
			p.SetPublisherId(v)
		}
	}
}

func ByPublisherName(v string) Filter {
	return func(f any) {
		if p, ok := f.(PublisherNameFilterable); ok {
			p.SetPublisherName(v)
		}
	}
}

func ByRating(v string) Filter {
	return func(f any) {
		if p, ok := f.(RatingFilterable); ok {
			p.SetRating(v)
		}
	}
}

func BySeriesID(v int) Filter {
	return func(f any) {
		if p, ok := f.(SeriesIDFilterable); ok {
			p.SetSeriesId(v)
		}
	}
}

func BySeriesName(v string) Filter {
	return func(f any) {
		if p, ok := f.(SeriesNameFilterable); ok {
			p.SetSeriesName(v)
		}
	}
}

func BySeriesType(v string) Filter {
	return func(f any) {
		if p, ok := f.(SeriesTypeFilterable); ok {
			p.SetSeriesType(v)
		}
	}
}

func BySeriesTypeID(v int) Filter {
	return func(f any) {
		if p, ok := f.(SeriesTypeIDFilterable); ok {
			p.SetSeriesTypeId(v)
		}
	}
}

func BySeriesVolume(v int) Filter {
	return func(f any) {
		if p, ok := f.(SeriesVolumeFilterable); ok {
			p.SetSeriesVolume(v)
		}
	}
}

func BySeriesYearBegan(v int) Filter {
	return func(f any) {
		if p, ok := f.(SeriesYearBeganFilterable); ok {
			p.SetSeriesYearBegan(v)
		}
	}
}

func ByStatus(v int) Filter {
	return func(f any) {
		if p, ok := f.(StatusFilterable); ok {
			p.SetStatus(v)
		}
	}
}

func ByStockKeepingUnit(v string) Filter {
	return func(f any) {
		if p, ok := f.(StockKeepingUnitFilterable); ok {
			p.SetSku(v)
		}
	}
}

func ByStoreDate(v civil.Date) Filter {
	return func(f any) {
		if p, ok := f.(StoreDateFilterable); ok {
			p.SetStoreDate(openapi_types.Date{Time: v.In(time.UTC)})
		}
	}
}

func ByStoreDateRangeAfter(v civil.Date) Filter {
	return func(f any) {
		if p, ok := f.(StoreDateRangeAfterFilterable); ok {
			p.SetStoreDateRangeAfter(openapi_types.Date{Time: v.In(time.UTC)})
		}
	}
}

func ByStoreDateRangeBefore(v civil.Date) Filter {
	return func(f any) {
		if p, ok := f.(StoreDateRangeBeforeFilterable); ok {
			p.SetStoreDateRangeBefore(openapi_types.Date{Time: v.In(time.UTC)})
		}
	}
}

func ByUniversalProductCode(v string) Filter {
	return func(f any) {
		if p, ok := f.(UniversalProductCodeFilterable); ok {
			p.SetUpc(v)
		}
	}
}

func ByVolume(v int) Filter {
	return func(f any) {
		if p, ok := f.(VolumeFilterable); ok {
			p.SetVolume(v)
		}
	}
}

func ByYearBegan(v int) Filter {
	return func(f any) {
		if p, ok := f.(YearBeganFilterable); ok {
			p.SetYearBegan(v)
		}
	}
}

func ByYearEnd(v int) Filter {
	return func(f any) {
		if p, ok := f.(YearEndFilterable); ok {
			p.SetYearEnd(v)
		}
	}
}

type AlternativeNumberFilterable interface {
	SetAltNumber(string)
}

type ComicVineIDFilterable interface {
	SetCvId(int)
}

type CoverHashFilterable interface {
	SetCoverHash(string)
}

type CoverMonthFilterable interface {
	SetCoverMonth(float32)
}

type CoverYearFilterable interface {
	SetCoverYear(float32)
}

type DesignationFilterable interface {
	SetDesignation(string)
}

type FinalOrderCutoffDateFilterable interface {
	SetFocDate(openapi_types.Date)
}

type FinalOrderCutoffDateRangeAfterFilterable interface {
	SetFocDateRangeAfter(openapi_types.Date)
}

type FinalOrderCutoffDateRangeBeforeFilterable interface {
	SetFocDateRangeBefore(openapi_types.Date)
}

type GrandComicsDatabaseIDFilterable interface {
	SetGcdId(int)
}

type ImprintIDFilterable interface {
	SetImprintId(int)
}

type ImprintNameFilterable interface {
	SetImprintName(string)
}

type MissingComicVineIDFilterable interface {
	SetMissingCvId(bool)
}

type MissingGrandComicsDatabaseIDFilterable interface {
	SetMissingGcdId(bool)
}

type ModifiedGreaterThanFilterable interface {
	SetModifiedGt(time.Time)
}

type NameFilterable interface {
	SetName(string)
}

type NumberFilterable interface {
	SetNumber(string)
}

type PublisherIDFilterable interface {
	SetPublisherId(int)
}

type PublisherNameFilterable interface {
	SetPublisherName(string)
}

type RatingFilterable interface {
	SetRating(string)
}

type SeriesIDFilterable interface {
	SetSeriesId(int)
}

type SeriesNameFilterable interface {
	SetSeriesName(string)
}

type SeriesTypeFilterable interface {
	SetSeriesType(string)
}

type SeriesTypeIDFilterable interface {
	SetSeriesTypeId(int)
}

type SeriesVolumeFilterable interface {
	SetSeriesVolume(int)
}

type SeriesYearBeganFilterable interface {
	SetSeriesYearBegan(int)
}

type StatusFilterable interface {
	SetStatus(int)
}

type StockKeepingUnitFilterable interface {
	SetSku(string)
}

type StoreDateFilterable interface {
	SetStoreDate(openapi_types.Date)
}

type StoreDateRangeAfterFilterable interface {
	SetStoreDateRangeAfter(openapi_types.Date)
}

type StoreDateRangeBeforeFilterable interface {
	SetStoreDateRangeBefore(openapi_types.Date)
}

type UniversalProductCodeFilterable interface {
	SetUpc(string)
}

type VolumeFilterable interface {
	SetVolume(int)
}

type YearBeganFilterable interface {
	SetYearBegan(int)
}

type YearEndFilterable interface {
	SetYearEnd(int)
}
