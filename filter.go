package metron

import (
	"slices"
	"time"

	"cloud.google.com/go/civil"
	"github.com/jyggen/go-metron/internal"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

// SeriesStatus is a series' publication status.
type SeriesStatus int

// The publication statuses Metron defines.
const (
	SeriesCancelled SeriesStatus = 1
	SeriesCompleted SeriesStatus = 2
	SeriesHiatus    SeriesStatus = 3
	SeriesOngoing   SeriesStatus = 4
)

// date converts a calendar date to the generated type, at midnight UTC.
func date(v civil.Date) *openapi_types.Date {
	return &openapi_types.Date{Time: v.In(time.UTC)}
}

// IssueFilters narrows an issue listing. A zero-valued field is omitted, so
// nil and an empty struct are equivalent.
//
// The From/To date pairs are inclusive at both ends.
type IssueFilters struct {
	AlternativeNumber        string
	CharacterID              int
	ComicVineID              int
	CoverHash                string
	CoverMonth               time.Month
	CoverYear                int
	CreatorID                int
	FinalOrderCutoffDate     civil.Date
	FinalOrderCutoffDateFrom civil.Date
	FinalOrderCutoffDateTo   civil.Date
	GrandComicsDatabaseID    int
	ImprintID                int
	ImprintName              string
	// MissingComicVineID selects records without a Comic Vine ID when true, and
	// records that have one when false. Nil does not filter.
	MissingComicVineID *bool
	// MissingGrandComicsDatabaseID behaves as MissingComicVineID.
	MissingGrandComicsDatabaseID *bool
	ModifiedAfter                time.Time
	Number                       string
	PublisherID                  int
	PublisherName                string
	Rating                       string
	RoleIDs                      []int
	SKU                          string
	SeriesAlternativeNames       string
	SeriesID                     int
	SeriesName                   string
	// SeriesSearch matches the series name and its alternative names at once.
	SeriesSearch    string
	SeriesVolume    int
	SeriesYearBegan int
	StoreDate       civil.Date
	StoreDateFrom   civil.Date
	StoreDateTo     civil.Date
	TeamID          int
	UPC             string
	// UPCStartsWith matches a UPC prefix, for scanners that drop the EAN supplement.
	UPCStartsWith string
	UniverseID    int
}

func (f *IssueFilters) params(page int) *internal.ApiIssueListParams {
	p := &internal.ApiIssueListParams{Page: &page}
	if f == nil {
		return p
	}

	if f.AlternativeNumber != "" {
		p.AltNumber = new(f.AlternativeNumber)
	}
	if f.CharacterID != 0 {
		p.CharacterId = new(f.CharacterID)
	}
	if f.ComicVineID != 0 {
		p.CvId = new(f.ComicVineID)
	}
	if f.CoverHash != "" {
		p.CoverHash = new(f.CoverHash)
	}
	if f.CoverMonth != 0 {
		p.CoverMonth = new(float32(f.CoverMonth))
	}
	if f.CoverYear != 0 {
		p.CoverYear = new(float32(f.CoverYear))
	}
	if f.CreatorID != 0 {
		p.CreatorId = new(f.CreatorID)
	}
	if f.FinalOrderCutoffDate.IsValid() {
		p.FocDate = date(f.FinalOrderCutoffDate)
	}
	if f.FinalOrderCutoffDateFrom.IsValid() {
		p.FocDateRangeAfter = date(f.FinalOrderCutoffDateFrom)
	}
	if f.FinalOrderCutoffDateTo.IsValid() {
		p.FocDateRangeBefore = date(f.FinalOrderCutoffDateTo)
	}
	if f.GrandComicsDatabaseID != 0 {
		p.GcdId = new(f.GrandComicsDatabaseID)
	}
	if f.ImprintID != 0 {
		p.ImprintId = new(f.ImprintID)
	}
	if f.ImprintName != "" {
		p.ImprintName = new(f.ImprintName)
	}
	if f.MissingComicVineID != nil {
		p.MissingCvId = new(*f.MissingComicVineID)
	}
	if f.MissingGrandComicsDatabaseID != nil {
		p.MissingGcdId = new(*f.MissingGrandComicsDatabaseID)
	}
	if !f.ModifiedAfter.IsZero() {
		p.ModifiedGt = new(f.ModifiedAfter)
	}
	if f.Number != "" {
		p.Number = new(f.Number)
	}
	if f.PublisherID != 0 {
		p.PublisherId = new(f.PublisherID)
	}
	if f.PublisherName != "" {
		p.PublisherName = new(f.PublisherName)
	}
	if f.Rating != "" {
		p.Rating = new(f.Rating)
	}
	if len(f.RoleIDs) > 0 {
		p.RoleId = new(slices.Clone(f.RoleIDs))
	}
	if f.SKU != "" {
		p.Sku = new(f.SKU)
	}
	if f.SeriesAlternativeNames != "" {
		p.SeriesAltNames = new(f.SeriesAlternativeNames)
	}
	if f.SeriesID != 0 {
		p.SeriesId = new(f.SeriesID)
	}
	if f.SeriesName != "" {
		p.SeriesName = new(f.SeriesName)
	}
	if f.SeriesSearch != "" {
		p.SeriesQ = new(f.SeriesSearch)
	}
	if f.SeriesVolume != 0 {
		p.SeriesVolume = new(f.SeriesVolume)
	}
	if f.SeriesYearBegan != 0 {
		p.SeriesYearBegan = new(f.SeriesYearBegan)
	}
	if f.StoreDate.IsValid() {
		p.StoreDate = date(f.StoreDate)
	}
	if f.StoreDateFrom.IsValid() {
		p.StoreDateRangeAfter = date(f.StoreDateFrom)
	}
	if f.StoreDateTo.IsValid() {
		p.StoreDateRangeBefore = date(f.StoreDateTo)
	}
	if f.TeamID != 0 {
		p.TeamId = new(f.TeamID)
	}
	if f.UPC != "" {
		p.Upc = new(f.UPC)
	}
	if f.UPCStartsWith != "" {
		p.UpcStartsWith = new(f.UPCStartsWith)
	}
	if f.UniverseID != 0 {
		p.UniverseId = new(f.UniverseID)
	}

	return p
}

// SeriesFilters narrows a series listing. A zero-valued field is omitted, so
// nil and an empty struct are equivalent.
type SeriesFilters struct {
	AlternativeNames      string
	CharacterID           int
	ComicVineID           int
	CreatorID             int
	GrandComicsDatabaseID int
	ImprintID             int
	ImprintName           string
	// MissingComicVineID selects records without a Comic Vine ID when true, and
	// records that have one when false. Nil does not filter.
	MissingComicVineID *bool
	// MissingGrandComicsDatabaseID behaves as MissingComicVineID.
	MissingGrandComicsDatabaseID *bool
	ModifiedAfter                time.Time
	Name                         string
	PublisherID                  int
	PublisherName                string
	RoleIDs                      []int
	// Search matches the series name and its alternative names at once.
	Search       string
	SeriesType   string
	SeriesTypeID int
	Status       SeriesStatus
	TeamID       int
	UniverseID   int
	Volume       int
	YearBegan    int
	YearEnd      int
}

func (f *SeriesFilters) params(page int) *internal.ApiSeriesListParams {
	p := &internal.ApiSeriesListParams{Page: &page}
	if f == nil {
		return p
	}

	if f.AlternativeNames != "" {
		p.AltNames = new(f.AlternativeNames)
	}
	if f.CharacterID != 0 {
		p.CharacterId = new(f.CharacterID)
	}
	if f.ComicVineID != 0 {
		p.CvId = new(f.ComicVineID)
	}
	if f.CreatorID != 0 {
		p.CreatorId = new(f.CreatorID)
	}
	if f.GrandComicsDatabaseID != 0 {
		p.GcdId = new(f.GrandComicsDatabaseID)
	}
	if f.ImprintID != 0 {
		p.ImprintId = new(f.ImprintID)
	}
	if f.ImprintName != "" {
		p.ImprintName = new(f.ImprintName)
	}
	if f.MissingComicVineID != nil {
		p.MissingCvId = new(*f.MissingComicVineID)
	}
	if f.MissingGrandComicsDatabaseID != nil {
		p.MissingGcdId = new(*f.MissingGrandComicsDatabaseID)
	}
	if !f.ModifiedAfter.IsZero() {
		p.ModifiedGt = new(f.ModifiedAfter)
	}
	if f.Name != "" {
		p.Name = new(f.Name)
	}
	if f.PublisherID != 0 {
		p.PublisherId = new(f.PublisherID)
	}
	if f.PublisherName != "" {
		p.PublisherName = new(f.PublisherName)
	}
	if len(f.RoleIDs) > 0 {
		p.RoleId = new(slices.Clone(f.RoleIDs))
	}
	if f.Search != "" {
		p.Q = new(f.Search)
	}
	if f.SeriesType != "" {
		p.SeriesType = new(f.SeriesType)
	}
	if f.SeriesTypeID != 0 {
		p.SeriesTypeId = new(f.SeriesTypeID)
	}
	if f.Status != 0 {
		p.Status = new(int(f.Status))
	}
	if f.TeamID != 0 {
		p.TeamId = new(f.TeamID)
	}
	if f.UniverseID != 0 {
		p.UniverseId = new(f.UniverseID)
	}
	if f.Volume != 0 {
		p.Volume = new(f.Volume)
	}
	if f.YearBegan != 0 {
		p.YearBegan = new(f.YearBegan)
	}
	if f.YearEnd != 0 {
		p.YearEnd = new(f.YearEnd)
	}

	return p
}

// ArcFilters narrows a story arc listing. A zero-valued field is omitted, so
// nil and an empty struct are equivalent.
type ArcFilters struct {
	ComicVineID           int
	GrandComicsDatabaseID int
	ModifiedAfter         time.Time
	Name                  string
}

func (f *ArcFilters) params(page int) *internal.ApiArcListParams {
	p := &internal.ApiArcListParams{Page: &page}
	if f == nil {
		return p
	}

	if f.ComicVineID != 0 {
		p.CvId = new(f.ComicVineID)
	}
	if f.GrandComicsDatabaseID != 0 {
		p.GcdId = new(f.GrandComicsDatabaseID)
	}
	if !f.ModifiedAfter.IsZero() {
		p.ModifiedGt = new(f.ModifiedAfter)
	}
	if f.Name != "" {
		p.Name = new(f.Name)
	}

	return p
}

// CharacterFilters narrows a character listing. A zero-valued field is
// omitted, so nil and an empty struct are equivalent.
type CharacterFilters struct {
	ComicVineID           int
	GrandComicsDatabaseID int
	ModifiedAfter         time.Time
	Name                  string
}

func (f *CharacterFilters) params(page int) *internal.ApiCharacterListParams {
	p := &internal.ApiCharacterListParams{Page: &page}
	if f == nil {
		return p
	}

	if f.ComicVineID != 0 {
		p.CvId = new(f.ComicVineID)
	}
	if f.GrandComicsDatabaseID != 0 {
		p.GcdId = new(f.GrandComicsDatabaseID)
	}
	if !f.ModifiedAfter.IsZero() {
		p.ModifiedGt = new(f.ModifiedAfter)
	}
	if f.Name != "" {
		p.Name = new(f.Name)
	}

	return p
}

// CreatorFilters narrows a creator listing. A zero-valued field is omitted, so
// nil and an empty struct are equivalent.
type CreatorFilters struct {
	ComicVineID           int
	GrandComicsDatabaseID int
	ModifiedAfter         time.Time
	Name                  string
}

func (f *CreatorFilters) params(page int) *internal.ApiCreatorListParams {
	p := &internal.ApiCreatorListParams{Page: &page}
	if f == nil {
		return p
	}

	if f.ComicVineID != 0 {
		p.CvId = new(f.ComicVineID)
	}
	if f.GrandComicsDatabaseID != 0 {
		p.GcdId = new(f.GrandComicsDatabaseID)
	}
	if !f.ModifiedAfter.IsZero() {
		p.ModifiedGt = new(f.ModifiedAfter)
	}
	if f.Name != "" {
		p.Name = new(f.Name)
	}

	return p
}

// ImprintFilters narrows an imprint listing. A zero-valued field is omitted,
// so nil and an empty struct are equivalent.
type ImprintFilters struct {
	ComicVineID           int
	GrandComicsDatabaseID int
	ModifiedAfter         time.Time
	Name                  string
}

func (f *ImprintFilters) params(page int) *internal.ApiImprintListParams {
	p := &internal.ApiImprintListParams{Page: &page}
	if f == nil {
		return p
	}

	if f.ComicVineID != 0 {
		p.CvId = new(f.ComicVineID)
	}
	if f.GrandComicsDatabaseID != 0 {
		p.GcdId = new(f.GrandComicsDatabaseID)
	}
	if !f.ModifiedAfter.IsZero() {
		p.ModifiedGt = new(f.ModifiedAfter)
	}
	if f.Name != "" {
		p.Name = new(f.Name)
	}

	return p
}

// PublisherFilters narrows a publisher listing. A zero-valued field is
// omitted, so nil and an empty struct are equivalent.
type PublisherFilters struct {
	ComicVineID           int
	GrandComicsDatabaseID int
	ModifiedAfter         time.Time
	Name                  string
}

func (f *PublisherFilters) params(page int) *internal.ApiPublisherListParams {
	p := &internal.ApiPublisherListParams{Page: &page}
	if f == nil {
		return p
	}

	if f.ComicVineID != 0 {
		p.CvId = new(f.ComicVineID)
	}
	if f.GrandComicsDatabaseID != 0 {
		p.GcdId = new(f.GrandComicsDatabaseID)
	}
	if !f.ModifiedAfter.IsZero() {
		p.ModifiedGt = new(f.ModifiedAfter)
	}
	if f.Name != "" {
		p.Name = new(f.Name)
	}

	return p
}

// TeamFilters narrows a team listing. A zero-valued field is omitted, so nil
// and an empty struct are equivalent.
type TeamFilters struct {
	ComicVineID           int
	GrandComicsDatabaseID int
	ModifiedAfter         time.Time
	Name                  string
}

func (f *TeamFilters) params(page int) *internal.ApiTeamListParams {
	p := &internal.ApiTeamListParams{Page: &page}
	if f == nil {
		return p
	}

	if f.ComicVineID != 0 {
		p.CvId = new(f.ComicVineID)
	}
	if f.GrandComicsDatabaseID != 0 {
		p.GcdId = new(f.GrandComicsDatabaseID)
	}
	if !f.ModifiedAfter.IsZero() {
		p.ModifiedGt = new(f.ModifiedAfter)
	}
	if f.Name != "" {
		p.Name = new(f.Name)
	}

	return p
}

// UniverseFilters narrows a universe listing. A zero-valued field is omitted,
// so nil and an empty struct are equivalent.
type UniverseFilters struct {
	Designation   string
	ModifiedAfter time.Time
	Name          string
}

func (f *UniverseFilters) params(page int) *internal.ApiUniverseListParams {
	p := &internal.ApiUniverseListParams{Page: &page}
	if f == nil {
		return p
	}

	if f.Designation != "" {
		p.Designation = new(f.Designation)
	}
	if !f.ModifiedAfter.IsZero() {
		p.ModifiedGt = new(f.ModifiedAfter)
	}
	if f.Name != "" {
		p.Name = new(f.Name)
	}

	return p
}

// RoleFilters narrows a creator role listing. A zero-valued field is omitted,
// so nil and an empty struct are equivalent.
type RoleFilters struct {
	ModifiedAfter time.Time
	Name          string
}

func (f *RoleFilters) params(page int) *internal.ApiRoleListParams {
	p := &internal.ApiRoleListParams{Page: &page}
	if f == nil {
		return p
	}

	if !f.ModifiedAfter.IsZero() {
		p.ModifiedGt = new(f.ModifiedAfter)
	}
	if f.Name != "" {
		p.Name = new(f.Name)
	}

	return p
}

// SeriesTypeFilters narrows a series type listing. A zero-valued field is
// omitted, so nil and an empty struct are equivalent.
type SeriesTypeFilters struct {
	ModifiedAfter time.Time
	Name          string
}

func (f *SeriesTypeFilters) params(page int) *internal.ApiSeriesTypeListParams {
	p := &internal.ApiSeriesTypeListParams{Page: &page}
	if f == nil {
		return p
	}

	if !f.ModifiedAfter.IsZero() {
		p.ModifiedGt = new(f.ModifiedAfter)
	}
	if f.Name != "" {
		p.Name = new(f.Name)
	}

	return p
}
