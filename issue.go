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

// IssueSeries is the series an issue belongs to, as embedded in an issue.
type IssueSeries struct {
	ID               int
	Name             string
	AlternativeNames []string
	SortName         string
	Volume           int
	YearBegan        int
	Type             Reference
	Genres           []Reference
}

// IssueCredit is a creator credited on an issue, together with their roles.
type IssueCredit struct {
	ID    int
	Name  string
	Roles []Reference
}

// IssueReprint is an issue reprinted by another issue.
type IssueReprint struct {
	ID    int
	Issue string
}

// IssueVariant is an alternative cover for an issue.
type IssueVariant struct {
	Name     *string
	SKU      *string
	UPC      *string
	Price    *string
	ImageURL url.URL
}

// IssueListSeries is the series an issue belongs to, as embedded in list responses.
type IssueListSeries struct {
	ID        int
	Name      string
	Volume    int
	YearBegan int
}

// Issue is a comic book issue.
type Issue struct {
	ID                    int
	Publisher             Reference
	Imprint               *Reference
	Series                IssueSeries
	Number                string
	AlternativeNumber     string
	Title                 *string
	Name                  []string
	CoverDate             civil.Date
	StoreDate             *civil.Date
	FinalOrderCutoffDate  *civil.Date
	Price                 *string
	PriceCurrency         string
	Rating                Reference
	SKU                   *string
	ISBN                  *string
	UPC                   *string
	PageCount             *int
	Description           *string
	ImageURL              *url.URL
	CoverHash             *string
	AverageRating         *float64
	RatingCount           int
	Arcs                  []ArcList
	Credits               []IssueCredit
	Characters            []CharacterList
	Teams                 []TeamList
	Universes             []UniverseList
	Reprints              []IssueReprint
	Variants              []IssueVariant
	ComicVineID           *int
	GrandComicsDatabaseID *int
	ResourceURL           url.URL
	Modified              time.Time
}

// IssueList is an issue as it appears in list responses.
type IssueList struct {
	ID        int
	Series    IssueListSeries
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

	if err := applyFilters("Issues", params, filters); err != nil {
		return errIter[*IssueList](err)
	}

	return paginate[internal.PaginatedIssueListList](ctx, c, "issue", c.client.ApiIssueList, issueListMapper, params)
}

// IssuesByArcID returns an iterator over all issues for a story arc.
func (c *Client) IssuesByArcID(ctx context.Context, id int, filters ...Filter) iter.Seq2[*IssueList, error] {
	params := &internal.ApiArcIssueListListParams{}

	if err := applyFilters("IssuesByArcID", params, filters); err != nil {
		return errIter[*IssueList](err)
	}

	return idPaginate[internal.PaginatedIssueListList](ctx, c, "arc/issue", c.client.ApiArcIssueListList, issueListMapper, id, params)
}

// IssuesByCharacterID returns an iterator over all issues for a character.
func (c *Client) IssuesByCharacterID(ctx context.Context, id int, filters ...Filter) iter.Seq2[*IssueList, error] {
	params := &internal.ApiCharacterIssueListListParams{}

	if err := applyFilters("IssuesByCharacterID", params, filters); err != nil {
		return errIter[*IssueList](err)
	}

	return idPaginate[internal.PaginatedIssueListList](ctx, c, "character/issue", c.client.ApiCharacterIssueListList, issueListMapper, id, params)
}

// IssuesBySeriesID returns an iterator over all issues for a series.
func (c *Client) IssuesBySeriesID(ctx context.Context, id int, filters ...Filter) iter.Seq2[*IssueList, error] {
	params := &internal.ApiSeriesIssueListListParams{}

	if err := applyFilters("IssuesBySeriesID", params, filters); err != nil {
		return errIter[*IssueList](err)
	}

	return idPaginate[internal.PaginatedIssueListList](ctx, c, "series/issue", c.client.ApiSeriesIssueListList, issueListMapper, id, params)
}

// IssuesByTeamID returns an iterator over all issues for a team.
func (c *Client) IssuesByTeamID(ctx context.Context, id int, filters ...Filter) iter.Seq2[*IssueList, error] {
	params := &internal.ApiTeamIssueListListParams{}

	if err := applyFilters("IssuesByTeamID", params, filters); err != nil {
		return errIter[*IssueList](err)
	}

	return idPaginate[internal.PaginatedIssueListList](ctx, c, "team/issue", c.client.ApiTeamIssueListList, issueListMapper, id, params)
}

func issueMapper(in internal.IssueRead) (*Issue, error) {
	if in.Id == nil {
		return nil, &MapError{Kind: "issue", Field: "Id"}
	}

	id := *in.Id

	if in.Modified == nil {
		return nil, &MapError{Kind: "issue", ID: id, Field: "Modified"}
	}

	if in.ResourceUrl == nil {
		return nil, &MapError{Kind: "issue", ID: id, Field: "ResourceUrl"}
	}

	if in.Publisher == nil {
		return nil, &MapError{Kind: "issue", ID: id, Field: "Publisher"}
	}

	if in.Publisher.Id == nil {
		return nil, &MapError{Kind: "issue", ID: id, Field: "Publisher.Id"}
	}

	if in.Series == nil {
		return nil, &MapError{Kind: "issue", ID: id, Field: "Series"}
	}

	if in.Series.Id == nil {
		return nil, &MapError{Kind: "issue", ID: id, Field: "Series.Id"}
	}

	if in.Series.SeriesType == nil {
		return nil, &MapError{Kind: "issue", ID: id, Field: "Series.SeriesType"}
	}

	if in.Series.SeriesType.Id == nil {
		return nil, &MapError{Kind: "issue", ID: id, Field: "Series.SeriesType.Id"}
	}

	if in.Rating == nil {
		return nil, &MapError{Kind: "issue", ID: id, Field: "Rating"}
	}

	if in.Rating.Id == nil {
		return nil, &MapError{Kind: "issue", ID: id, Field: "Rating.Id"}
	}

	if in.PriceCurrency == nil {
		return nil, &MapError{Kind: "issue", ID: id, Field: "PriceCurrency"}
	}

	// alt_number and name are absent from the spec's required list, and
	// rating_count is a readOnly aggregate that can come back null. Guarded
	// anyway: a record missing them is surfaced, not silently zeroed.
	if in.AltNumber == nil {
		return nil, &MapError{Kind: "issue", ID: id, Field: "AltNumber"}
	}

	if in.Name == nil {
		return nil, &MapError{Kind: "issue", ID: id, Field: "Name"}
	}

	if in.RatingCount == nil {
		return nil, &MapError{Kind: "issue", ID: id, Field: "RatingCount"}
	}

	var imageURL *url.URL
	var err error

	if image := nullableToPtr(in.Image); image != nil {
		imageURL, err = url.Parse(*image)
		if err != nil {
			return nil, &MapError{Kind: "issue", ID: id, Field: "Image", Err: err}
		}
	}

	resourceURL, err := url.Parse(*in.ResourceUrl)
	if err != nil {
		return nil, &MapError{Kind: "issue", ID: id, Field: "ResourceUrl", Err: err}
	}

	var imprint *Reference

	if in.Imprint != nil {
		if in.Imprint.Id == nil {
			return nil, &MapError{Kind: "issue", ID: id, Field: "Imprint.Id"}
		}

		imprint = &Reference{
			ID:   *in.Imprint.Id,
			Name: in.Imprint.Name,
		}
	}

	var seriesAltNames []string

	if in.Series.AltNames != nil {
		seriesAltNames = *in.Series.AltNames
	}

	var genres []Reference

	if in.Series.Genres != nil {
		genres = make([]Reference, 0, len(*in.Series.Genres))

		for _, genre := range *in.Series.Genres {
			if genre.Id == nil {
				return nil, &MapError{Kind: "issue", ID: id, Field: "Series.Genres[].Id"}
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
		maybeFinalOrderCutoffDate = new(civil.DateOf(d.Time))
	}

	var maybeStoreDate *civil.Date

	if d, err := in.StoreDate.Get(); err == nil {
		maybeStoreDate = new(civil.DateOf(d.Time))
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

	var credits []IssueCredit

	if in.Credits != nil {
		credits = make([]IssueCredit, 0, len(*in.Credits))

		for _, credit := range *in.Credits {
			if credit.Id == nil {
				return nil, &MapError{Kind: "issue", ID: id, Field: "Credits[].Id"}
			}

			if credit.Creator == nil {
				return nil, &MapError{Kind: "issue", ID: id, Field: "Credits[].Creator"}
			}

			roles := make([]Reference, 0, len(credit.Role))

			for _, role := range credit.Role {
				if role.Id == nil {
					return nil, &MapError{Kind: "issue", ID: id, Field: "Credits[].Role[].Id"}
				}

				roles = append(roles, Reference{
					ID:   *role.Id,
					Name: role.Name,
				})
			}

			credits = append(credits, IssueCredit{
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

	var reprints []IssueReprint

	if in.Reprints != nil {
		reprints = make([]IssueReprint, 0, len(*in.Reprints))

		for _, reprint := range *in.Reprints {
			if reprint.Id == nil {
				return nil, &MapError{Kind: "issue", ID: id, Field: "Reprints[].Id"}
			}

			reprints = append(reprints, IssueReprint{
				ID:    *reprint.Id,
				Issue: reprint.Issue,
			})
		}
	}

	var variants []IssueVariant

	if in.Variants != nil {
		variants = make([]IssueVariant, 0, len(*in.Variants))

		for _, variant := range *in.Variants {
			maybeVariantImageURL, innerErr := url.Parse(variant.Image)
			if innerErr != nil {
				return nil, &MapError{Kind: "issue", ID: id, Field: "Variants[].Image", Err: innerErr}
			}
			variantImageURL := *maybeVariantImageURL

			var variantPrice *string

			if p, priceErr := variant.Price.Get(); priceErr == nil {
				s, asErr := p.AsVariantsIssuePrice0()
				if asErr != nil {
					return nil, &MapError{Kind: "issue", ID: id, Field: "Variants[].Price", Err: asErr}
				}
				variantPrice = &s
			}

			variants = append(variants, IssueVariant{
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
			return nil, &MapError{Kind: "issue", ID: id, Field: "Price", Err: asErr}
		}
		price = &s
	}

	return &Issue{
		ID: id,
		Publisher: Reference{
			ID:   *in.Publisher.Id,
			Name: in.Publisher.Name,
		},
		Imprint: imprint,
		Series: IssueSeries{
			ID:               *in.Series.Id,
			Name:             in.Series.Name,
			AlternativeNames: seriesAltNames,
			SortName:         in.Series.SortName,
			Volume:           in.Series.Volume,
			YearBegan:        in.Series.YearBegan,
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
		PriceCurrency:        *in.PriceCurrency,
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
		AverageRating:         in.AverageRating,
		RatingCount:           *in.RatingCount,
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
		return nil, &MapError{Kind: "issue", Field: "Id"}
	}

	id := *in.Id

	if in.Modified == nil {
		return nil, &MapError{Kind: "issue", ID: id, Field: "Modified"}
	}

	if in.Series == nil {
		return nil, &MapError{Kind: "issue", ID: id, Field: "Series"}
	}

	if in.Series.Id == nil {
		return nil, &MapError{Kind: "issue", ID: id, Field: "Series.Id"}
	}

	coverDate := civil.DateOf(in.CoverDate.Time)

	var maybeStoreDate *civil.Date

	if d, err := in.StoreDate.Get(); err == nil {
		maybeStoreDate = new(civil.DateOf(d.Time))
	}

	var imageURL *url.URL
	var err error

	if image := nullableToPtr(in.Image); image != nil {
		imageURL, err = url.Parse(*image)
		if err != nil {
			return nil, &MapError{Kind: "issue", ID: id, Field: "Image", Err: err}
		}
	}

	return &IssueList{
		ID: id,
		Series: IssueListSeries{
			ID:        *in.Series.Id,
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
