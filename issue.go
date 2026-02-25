package metron

import (
	"context"
	"fmt"
	"iter"
	"net/url"
	"time"

	"cloud.google.com/go/civil"
	"github.com/jyggen/go-metron/internal"
)

type Issue struct {
	ID        int
	Publisher Reference
	Imprint   *Reference
	Series    struct {
		ID        int
		Name      string
		SortName  string
		Volume    int
		YearBegan int
		Type      Reference
		Genres    []Reference
	}
	Number               string
	AlternativeNumber    string
	Title                *string
	Name                 []string
	CoverDate            civil.Date
	StoreDate            *civil.Date
	FinalOrderCutoffDate *civil.Date
	Price                string
	Rating               Reference
	SKU                  *string
	ISBN                 *string
	UPC                  *string
	PageCount            *int
	Description          *string
	ImageURL             *url.URL
	CoverHash            *string
	Arcs                 []ArcList
	Credits              []struct {
		ID    int
		Name  string
		Roles []Reference
	}
	Characters []CharacterList
	Teams      []TeamList
	Universes  []UniverseList
	Reprints   []struct {
		ID    int
		Issue string
	}
	Variants []struct {
		Name     *string
		SKU      *string
		UPC      *string
		ImageURL url.URL
	}
	ComicVineID           *int
	GrandComicsDatabaseID *int
	ResourceURL           url.URL
	Modified              time.Time
}

type IssueList struct {
	ID     int
	Series struct {
		Name      string
		Volume    int
		YearBegan int
	}
	Name      string
	Number    string
	CoverDate civil.Date
	StoreDate *civil.Date
	ImageURL  *url.URL
	CoverHash *string
	Modified  time.Time
}

// IssueByID returns the information of an individual issue.
func (c *Client) IssueByID(ctx context.Context, id int) (*Issue, error) {
	return newByID(ctx, c.cache, fmt.Sprintf("issue/%d", id), c.client.ApiIssueRetrieve, issueMapper, id)
}

// Issues returns a list of all the issues.
func (c *Client) Issues(ctx context.Context, filters ...Filter) iter.Seq2[*IssueList, error] {
	params := &internal.ApiIssueListParams{}

	for _, f := range filters {
		f(params)
	}

	return newPaginate[internal.PaginatedIssueListList](ctx, c.client.ApiIssueList, issueListMapper, params)
}

// IssuesByArcID returns a list of all the issues for a story arc.
func (c *Client) IssuesByArcID(ctx context.Context, id int, filters ...Filter) iter.Seq2[*IssueList, error] {
	params := &internal.ApiArcIssueListListParams{}

	for _, f := range filters {
		f(params)
	}

	return newIDPaginate[internal.PaginatedIssueListList](ctx, c.client.ApiArcIssueListList, issueListMapper, id, params)
}

// IssuesByCharacterID returns a list of all the issues for a character.
func (c *Client) IssuesByCharacterID(ctx context.Context, id int, filters ...Filter) iter.Seq2[*IssueList, error] {
	params := &internal.ApiCharacterIssueListListParams{}

	for _, f := range filters {
		f(params)
	}

	return newIDPaginate[internal.PaginatedIssueListList](ctx, c.client.ApiCharacterIssueListList, issueListMapper, id, params)
}

// IssuesBySeriesID returns a list of all the issues for a series.
func (c *Client) IssuesBySeriesID(ctx context.Context, id int, filters ...Filter) iter.Seq2[*IssueList, error] {
	params := &internal.ApiSeriesIssueListListParams{}

	for _, f := range filters {
		f(params)
	}

	return newIDPaginate[internal.PaginatedIssueListList](ctx, c.client.ApiSeriesIssueListList, issueListMapper, id, params)
}

// IssuesByTeamID returns a list of all the issues for a team.
func (c *Client) IssuesByTeamID(ctx context.Context, id int, filters ...Filter) iter.Seq2[*IssueList, error] {
	params := &internal.ApiTeamIssueListListParams{}

	for _, f := range filters {
		f(params)
	}

	return newIDPaginate[internal.PaginatedIssueListList](ctx, c.client.ApiTeamIssueListList, issueListMapper, id, params)
}

func issueMapper(in internal.IssueRead) (*Issue, error) {
	var imageURL *url.URL
	var err error

	if in.Image != nil {
		imageURL, err = url.Parse(*in.Image)
		if err != nil {
			return nil, err
		}
	}

	resourceURL, err := url.Parse(*in.ResourceUrl)
	if err != nil {
		return nil, err
	}

	var imprint *Reference

	if in.Imprint != nil {
		imprint = &Reference{
			ID:   *in.Imprint.Id,
			Name: in.Imprint.Name,
		}
	}

	var genres []Reference

	if in.Series.Genres != nil {
		genres = make([]Reference, 0, len(*in.Series.Genres))

		for _, genre := range *in.Series.Genres {
			genres = append(genres, Reference{
				ID:   *genre.Id,
				Name: genre.Name,
			})
		}
	}

	coverDate := civil.DateOf(in.CoverDate.Time)

	var maybeFinalOrderCutoffDate *civil.Date

	if in.FocDate != nil {
		finalOrderCutoffDate := civil.DateOf(in.FocDate.Time)
		maybeFinalOrderCutoffDate = &finalOrderCutoffDate
	}

	var maybeStoreDate *civil.Date

	if in.StoreDate != nil {
		storeDate := civil.DateOf(in.StoreDate.Time)
		maybeStoreDate = &storeDate
	}

	arcs := make([]ArcList, 0, len(*in.Arcs))

	for _, arc := range *in.Arcs {
		a, innerErr := arcListMapper(arc)
		if innerErr != nil {
			return nil, innerErr
		}
		arcs = append(arcs, *a)
	}

	credits := make([]struct {
		ID    int
		Name  string
		Roles []Reference
	}, 0, len(*in.Credits))

	for _, credit := range *in.Credits {
		roles := make([]Reference, 0, len(credit.Role))

		for _, role := range credit.Role {
			roles = append(roles, Reference{
				ID:   *role.Id,
				Name: role.Name,
			})
		}

		credits = append(credits, struct {
			ID    int
			Name  string
			Roles []Reference
		}{
			ID:    *credit.Id,
			Name:  *credit.Creator,
			Roles: roles,
		})
	}

	characters := make([]CharacterList, 0, len(*in.Characters))

	for _, character := range *in.Characters {
		c, innerErr := characterListMapper(character)
		if innerErr != nil {
			return nil, innerErr
		}
		characters = append(characters, *c)
	}

	teams := make([]TeamList, 0, len(*in.Teams))

	for _, team := range *in.Teams {
		t, innerErr := teamListMapper(team)
		if innerErr != nil {
			return nil, innerErr
		}
		teams = append(teams, *t)
	}

	universes := make([]UniverseList, 0, len(*in.Universes))

	for _, universe := range *in.Universes {
		u, innerErr := universeListMapper(universe)
		if innerErr != nil {
			return nil, innerErr
		}
		universes = append(universes, *u)
	}

	reprints := make([]struct {
		ID    int
		Issue string
	}, 0, len(*in.Reprints))

	for _, reprint := range *in.Reprints {
		reprints = append(reprints, struct {
			ID    int
			Issue string
		}{
			ID:    *reprint.Id,
			Issue: reprint.Issue,
		})
	}

	variants := make([]struct {
		Name     *string
		SKU      *string
		UPC      *string
		ImageURL url.URL
	}, 0, len(*in.Variants))

	for _, variant := range *in.Variants {
		maybeVariantImageURL, innerErr := url.Parse(variant.Image)
		if innerErr != nil {
			return nil, innerErr
		}
		variantImageURL := *maybeVariantImageURL

		variants = append(variants, struct {
			Name     *string
			SKU      *string
			UPC      *string
			ImageURL url.URL
		}{
			Name:     variant.Name,
			SKU:      variant.Sku,
			UPC:      variant.Upc,
			ImageURL: variantImageURL,
		})
	}

	return &Issue{
		ID: *in.Id,
		Publisher: Reference{
			ID:   *in.Publisher.Id,
			Name: in.Publisher.Name,
		},
		Imprint: imprint,
		Series: struct {
			ID        int
			Name      string
			SortName  string
			Volume    int
			YearBegan int
			Type      Reference
			Genres    []Reference
		}{
			ID:        *in.Series.Id,
			Name:      in.Series.Name,
			SortName:  in.Series.SortName,
			Volume:    in.Series.Volume,
			YearBegan: in.Series.YearBegan,
			Type: Reference{
				ID:   *in.Series.SeriesType.Id,
				Name: in.Series.SeriesType.Name,
			},
			Genres: genres,
		},
		Number:               in.Number,
		AlternativeNumber:    *in.AltNumber,
		Title:                in.Title,
		Name:                 *in.Name,
		CoverDate:            coverDate,
		StoreDate:            maybeStoreDate,
		FinalOrderCutoffDate: maybeFinalOrderCutoffDate,
		Price:                *in.Price,
		Rating: Reference{
			ID:   *in.Rating.Id,
			Name: in.Rating.Name,
		},
		SKU:                   in.Sku,
		ISBN:                  in.Isbn,
		UPC:                   in.Upc,
		PageCount:             in.Page,
		Description:           in.Desc,
		ImageURL:              imageURL,
		CoverHash:             in.CoverHash,
		Arcs:                  arcs,
		Credits:               credits,
		Characters:            characters,
		Teams:                 teams,
		Universes:             universes,
		Reprints:              reprints,
		Variants:              variants,
		ComicVineID:           in.CvId,
		GrandComicsDatabaseID: in.GcdId,
		ResourceURL:           *resourceURL,
		Modified:              *in.Modified,
	}, nil
}

func issueListMapper(in internal.IssueList) (*IssueList, error) {
	coverDate := civil.DateOf(in.CoverDate.Time)

	var maybeStoreDate *civil.Date

	if in.StoreDate != nil {
		storeDate := civil.DateOf(in.StoreDate.Time)
		maybeStoreDate = &storeDate
	}

	var imageURL *url.URL
	var err error

	if in.Image != nil {
		imageURL, err = url.Parse(*in.Image)
		if err != nil {
			return nil, err
		}
	}

	return &IssueList{
		ID: *in.Id,
		Series: struct {
			Name      string
			Volume    int
			YearBegan int
		}{
			Name:      in.Series.Name,
			Volume:    in.Series.Volume,
			YearBegan: in.Series.YearBegan,
		},
		Name:      in.Issue,
		Number:    in.Number,
		CoverDate: coverDate,
		StoreDate: maybeStoreDate,
		ImageURL:  imageURL,
		CoverHash: in.CoverHash,
		Modified:  *in.Modified,
	}, nil
}
