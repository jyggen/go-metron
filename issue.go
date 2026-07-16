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

// Issue is a comic book issue.
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
	Price                *string
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
		Price    *string
		ImageURL url.URL
	}
	ComicVineID           *int
	GrandComicsDatabaseID *int
	ResourceURL           url.URL
	Modified              time.Time
}

// IssueList is an issue as it appears in list responses.
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

// IssueByID returns an issue by its ID.
func (c *Client) IssueByID(ctx context.Context, id int) (*Issue, error) {
	return byID(ctx, c, fmt.Sprintf("issue/%d", id), c.client.ApiIssueRetrieve, issueMapper, id)
}

// Issues returns an iterator over all issues.
func (c *Client) Issues(ctx context.Context, filters ...Filter) iter.Seq2[*IssueList, error] {
	params := &internal.ApiIssueListParams{}

	for _, f := range filters {
		f(params)
	}

	return paginate[internal.PaginatedIssueListList](ctx, c, "issue", c.client.ApiIssueList, issueListMapper, params)
}

// IssuesByArcID returns an iterator over all issues for a story arc.
func (c *Client) IssuesByArcID(ctx context.Context, id int, filters ...Filter) iter.Seq2[*IssueList, error] {
	params := &internal.ApiArcIssueListListParams{}

	for _, f := range filters {
		f(params)
	}

	return idPaginate[internal.PaginatedIssueListList](ctx, c, "arc/issue", c.client.ApiArcIssueListList, issueListMapper, id, params)
}

// IssuesByCharacterID returns an iterator over all issues for a character.
func (c *Client) IssuesByCharacterID(ctx context.Context, id int, filters ...Filter) iter.Seq2[*IssueList, error] {
	params := &internal.ApiCharacterIssueListListParams{}

	for _, f := range filters {
		f(params)
	}

	return idPaginate[internal.PaginatedIssueListList](ctx, c, "character/issue", c.client.ApiCharacterIssueListList, issueListMapper, id, params)
}

// IssuesBySeriesID returns an iterator over all issues for a series.
func (c *Client) IssuesBySeriesID(ctx context.Context, id int, filters ...Filter) iter.Seq2[*IssueList, error] {
	params := &internal.ApiSeriesIssueListListParams{}

	for _, f := range filters {
		f(params)
	}

	return idPaginate[internal.PaginatedIssueListList](ctx, c, "series/issue", c.client.ApiSeriesIssueListList, issueListMapper, id, params)
}

// IssuesByTeamID returns an iterator over all issues for a team.
func (c *Client) IssuesByTeamID(ctx context.Context, id int, filters ...Filter) iter.Seq2[*IssueList, error] {
	params := &internal.ApiTeamIssueListListParams{}

	for _, f := range filters {
		f(params)
	}

	return idPaginate[internal.PaginatedIssueListList](ctx, c, "team/issue", c.client.ApiTeamIssueListList, issueListMapper, id, params)
}

func issueMapper(in internal.IssueRead) (*Issue, error) {
	if in.Id == nil {
		return nil, fmt.Errorf("issue: nil Id")
	}

	if in.Modified == nil {
		return nil, fmt.Errorf("issue: nil Modified")
	}

	if in.ResourceUrl == nil {
		return nil, fmt.Errorf("issue: nil ResourceUrl")
	}

	if in.Publisher == nil {
		return nil, fmt.Errorf("issue: nil Publisher")
	}

	if in.Publisher.Id == nil {
		return nil, fmt.Errorf("issue: nil Publisher.Id")
	}

	if in.Series == nil {
		return nil, fmt.Errorf("issue: nil Series")
	}

	if in.Series.Id == nil {
		return nil, fmt.Errorf("issue: nil Series.Id")
	}

	if in.Series.SeriesType == nil {
		return nil, fmt.Errorf("issue: nil Series.SeriesType")
	}

	if in.Series.SeriesType.Id == nil {
		return nil, fmt.Errorf("issue: nil Series.SeriesType.Id")
	}

	if in.Rating == nil {
		return nil, fmt.Errorf("issue: nil Rating")
	}

	if in.AltNumber == nil {
		return nil, fmt.Errorf("issue: nil AltNumber")
	}

	if in.Name == nil {
		return nil, fmt.Errorf("issue: nil Name")
	}

	if in.Rating.Id == nil {
		return nil, fmt.Errorf("issue: nil Rating.Id")
	}

	var imageURL *url.URL
	var err error

	if image := nullableToPtr(in.Image); image != nil {
		imageURL, err = url.Parse(*image)
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
		if in.Imprint.Id == nil {
			return nil, fmt.Errorf("issue: nil Imprint.Id")
		}

		imprint = &Reference{
			ID:   *in.Imprint.Id,
			Name: in.Imprint.Name,
		}
	}

	var genres []Reference

	if in.Series.Genres != nil {
		genres = make([]Reference, 0, len(*in.Series.Genres))

		for _, genre := range *in.Series.Genres {
			if genre.Id == nil {
				return nil, fmt.Errorf("issue: nil Series.Genres[].Id")
			}

			genres = append(genres, Reference{
				ID:   *genre.Id,
				Name: genre.Name,
			})
		}
	}

	coverDate := civil.DateOf(in.CoverDate.Time)

	var maybeFinalOrderCutoffDate *civil.Date

	if d, err := in.FocDate.Get(); err == nil {
		finalOrderCutoffDate := civil.DateOf(d.Time)
		maybeFinalOrderCutoffDate = &finalOrderCutoffDate
	}

	var maybeStoreDate *civil.Date

	if d, err := in.StoreDate.Get(); err == nil {
		storeDate := civil.DateOf(d.Time)
		maybeStoreDate = &storeDate
	}

	var arcs []ArcList

	if in.Arcs != nil {
		arcs = make([]ArcList, 0, len(*in.Arcs))

		for _, arc := range *in.Arcs {
			a, innerErr := arcListMapper(arc)
			if innerErr != nil {
				return nil, innerErr
			}
			arcs = append(arcs, *a)
		}
	}

	var credits []struct {
		ID    int
		Name  string
		Roles []Reference
	}

	if in.Credits != nil {
		credits = make([]struct {
			ID    int
			Name  string
			Roles []Reference
		}, 0, len(*in.Credits))

		for _, credit := range *in.Credits {
			if credit.Id == nil {
				return nil, fmt.Errorf("issue: nil Credits[].Id")
			}

			if credit.Creator == nil {
				return nil, fmt.Errorf("issue: nil Credits[].Creator")
			}

			roles := make([]Reference, 0, len(credit.Role))

			for _, role := range credit.Role {
				if role.Id == nil {
					return nil, fmt.Errorf("issue: nil Credits[].Role[].Id")
				}

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
	}

	var characters []CharacterList

	if in.Characters != nil {
		characters = make([]CharacterList, 0, len(*in.Characters))

		for _, character := range *in.Characters {
			c, innerErr := characterListMapper(character)
			if innerErr != nil {
				return nil, innerErr
			}
			characters = append(characters, *c)
		}
	}

	var teams []TeamList

	if in.Teams != nil {
		teams = make([]TeamList, 0, len(*in.Teams))

		for _, team := range *in.Teams {
			t, innerErr := teamListMapper(team)
			if innerErr != nil {
				return nil, innerErr
			}
			teams = append(teams, *t)
		}
	}

	var universes []UniverseList

	if in.Universes != nil {
		universes = make([]UniverseList, 0, len(*in.Universes))

		for _, universe := range *in.Universes {
			u, innerErr := universeListMapper(universe)
			if innerErr != nil {
				return nil, innerErr
			}
			universes = append(universes, *u)
		}
	}

	var reprints []struct {
		ID    int
		Issue string
	}

	if in.Reprints != nil {
		reprints = make([]struct {
			ID    int
			Issue string
		}, 0, len(*in.Reprints))

		for _, reprint := range *in.Reprints {
			if reprint.Id == nil {
				return nil, fmt.Errorf("issue: nil Reprints[].Id")
			}

			reprints = append(reprints, struct {
				ID    int
				Issue string
			}{
				ID:    *reprint.Id,
				Issue: reprint.Issue,
			})
		}
	}

	var variants []struct {
		Name     *string
		SKU      *string
		UPC      *string
		Price    *string
		ImageURL url.URL
	}

	if in.Variants != nil {
		variants = make([]struct {
			Name     *string
			SKU      *string
			UPC      *string
			Price    *string
			ImageURL url.URL
		}, 0, len(*in.Variants))

		for _, variant := range *in.Variants {
			maybeVariantImageURL, innerErr := url.Parse(variant.Image)
			if innerErr != nil {
				return nil, innerErr
			}
			variantImageURL := *maybeVariantImageURL

			var variantPrice *string

			if p, priceErr := variant.Price.Get(); priceErr == nil {
				s, asErr := p.AsVariantsIssuePrice0()
				if asErr != nil {
					return nil, fmt.Errorf("issue: variants price: %w", asErr)
				}
				variantPrice = &s
			}

			variants = append(variants, struct {
				Name     *string
				SKU      *string
				UPC      *string
				Price    *string
				ImageURL url.URL
			}{
				Name:     variant.Name,
				SKU:      variant.Sku,
				UPC:      variant.Upc,
				Price:    variantPrice,
				ImageURL: variantImageURL,
			})
		}
	}

	var price *string

	if p, priceErr := in.Price.Get(); priceErr == nil {
		s, asErr := p.AsIssueReadPrice0()
		if asErr != nil {
			return nil, fmt.Errorf("issue: price: %w", asErr)
		}
		price = &s
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
		Price:                price,
		Rating: Reference{
			ID:   *in.Rating.Id,
			Name: in.Rating.Name,
		},
		SKU:                   in.Sku,
		ISBN:                  in.Isbn,
		UPC:                   in.Upc,
		PageCount:             nullableToPtr(in.Page),
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
		ComicVineID:           nullableToPtr(in.CvId),
		GrandComicsDatabaseID: nullableToPtr(in.GcdId),
		ResourceURL:           *resourceURL,
		Modified:              *in.Modified,
	}, nil
}

func issueListMapper(in internal.IssueList) (*IssueList, error) {
	if in.Id == nil {
		return nil, fmt.Errorf("issue: nil Id")
	}

	if in.Modified == nil {
		return nil, fmt.Errorf("issue: nil Modified")
	}

	if in.Series == nil {
		return nil, fmt.Errorf("issue: nil Series")
	}

	coverDate := civil.DateOf(in.CoverDate.Time)

	var maybeStoreDate *civil.Date

	if d, err := in.StoreDate.Get(); err == nil {
		storeDate := civil.DateOf(d.Time)
		maybeStoreDate = &storeDate
	}

	var imageURL *url.URL
	var err error

	if image := nullableToPtr(in.Image); image != nil {
		imageURL, err = url.Parse(*image)
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
