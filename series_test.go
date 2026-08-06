package metron_test

import (
	"testing"

	"github.com/jyggen/go-metron"
)

func TestSeriesByID(t *testing.T) {
	t.Parallel()
	testByID(t, "series", (*metron.Client).SeriesByID, []testCase[*metron.Series]{
		{
			id: 3371,
			expected: &metron.Series{
				ID:               3371,
				Name:             "Batman 2022 Annual",
				AlternativeNames: []string{},
				SortName:         "Batman 2022 Annual",
				Volume:           1,
				Type: metron.Reference{
					ID:   6,
					Name: "Annual",
				},
				Status: "Completed",
				Publisher: metron.Reference{
					ID:   2,
					Name: "DC Comics",
				},
				Imprint:     nil,
				YearBegan:   2022,
				YearEnded:   nil,
				Description: new(""),
				IssueCount:  1,
				Genres: []metron.Reference{
					{
						ID:   10,
						Name: "Super-Hero",
					},
				},
				Associated: []metron.Reference{
					{
						ID:   93,
						Name: "Batman (2016)",
					},
				},
				ComicVineID:           new(143255),
				GrandComicsDatabaseID: new(185580),
				ResourceURL:           parseURL(t, "https://metron.cloud/series/batman-2022-annual-2022/"),
				Modified:              parseTime(t, "2024-12-21T13:48:18.897702-05:00"),
			},
		},
		{
			id: 793,
			expected: &metron.Series{
				ID:               793,
				Name:             "Fables",
				AlternativeNames: []string{},
				SortName:         "Fables",
				Volume:           1,
				Type: metron.Reference{
					ID:   13,
					Name: "Single Issue",
				},
				Status: "Cancelled",
				Publisher: metron.Reference{
					ID:   2,
					Name: "DC Comics",
				},
				Imprint: &metron.Reference{
					ID:   1,
					Name: "Vertigo Comics",
				},
				YearBegan: 2002,
				YearEnded: new(2015),
				Description: new(
					"When the Adversary conquered the lands of legends, the inhabitants were forced into exile. They form a secret society, a hidden enclave in modern-day New York. Bill Willingham's award-winning \"fractured-fairy-tale\" series explores the world of these beloved fables...one that exists within our own.",
				),
				IssueCount: 149,
				Genres:     []metron.Reference{},
				Associated: []metron.Reference{
					{
						ID:   3396,
						Name: "Fables (2022)",
					},
				},
				ComicVineID:           new(9723),
				GrandComicsDatabaseID: new(10549),
				ResourceURL:           parseURL(t, "https://metron.cloud/series/fables-2002/"),
				Modified:              parseTime(t, "2026-05-29T18:38:10.592560-04:00"),
			},
		},
	})
}

func TestSeriesByPublisherID(t *testing.T) {
	t.Parallel()
	testListByID(t, "publisher", 1, "series", (*metron.Client).SeriesByPublisherID, publisherSeriesListTestCases(t))
}

func TestSeries(t *testing.T) {
	t.Parallel()
	testList(t, "series", (*metron.Client).Series, seriesListTestCases(t))
}

func seriesListTestCases(t *testing.T) []testCase[*metron.SeriesList] {
	return []testCase[*metron.SeriesList]{
		{
			id: 3856,
			expected: &metron.SeriesList{
				ID:         3856,
				Name:       "The 06 Protocol (2022)",
				YearBegan:  2022,
				YearEnded:  new(2022),
				Volume:     1,
				IssueCount: 3,
				Modified:   parseTime(t, "2024-12-23T16:13:24.914552-05:00"),
			},
		},
		{
			id: 10307,
			expected: &metron.SeriesList{
				ID:         10307,
				Name:       "8-Bit Zombie (2013)",
				YearBegan:  2013,
				YearEnded:  new(2013),
				Volume:     1,
				IssueCount: 1,
				Modified:   parseTime(t, "2025-02-25T22:49:53.851020-05:00"),
			},
		},
		{
			id: 10311,
			expected: &metron.SeriesList{
				ID:         10311,
				Name:       "8-Bit Zombie: The Full Byte TPB (2015)",
				YearBegan:  2015,
				YearEnded:  new(2015),
				Volume:     1,
				IssueCount: 1,
				Modified:   parseTime(t, "2025-02-25T23:02:02.115714-05:00"),
			},
		},
		{
			id: 6195,
			expected: &metron.SeriesList{
				ID:         6195,
				Name:       "100% (2002)",
				YearBegan:  2002,
				YearEnded:  new(2003),
				Volume:     1,
				IssueCount: 5,
				Modified:   parseTime(t, "2025-01-19T11:26:32.503473-05:00"),
			},
		},
	}
}

func publisherSeriesListTestCases(t *testing.T) []testCase[*metron.SeriesList] {
	return []testCase[*metron.SeriesList]{
		{
			id: 8531,
			expected: &metron.SeriesList{
				ID:         8531,
				Name:       "The 100 Greatest Marvels of All Time (2001)",
				YearBegan:  2001,
				YearEnded:  new(2001),
				Volume:     1,
				IssueCount: 10,
				Modified:   parseTime(t, "2024-12-26T16:46:05.803827-05:00"),
			},
		},
		{
			id: 1530,
			expected: &metron.SeriesList{
				ID:         1530,
				Name:       "100ᵗʰ Anniversary Special: Avengers (2014)",
				YearBegan:  2014,
				YearEnded:  nil,
				Volume:     1,
				IssueCount: 1,
				Modified:   parseTime(t, "2024-12-26T16:46:20.756441-05:00"),
			},
		},
		{
			id: 1531,
			expected: &metron.SeriesList{
				ID:         1531,
				Name:       "100ᵗʰ Anniversary Special: Fantastic Four (2014)",
				YearBegan:  2014,
				YearEnded:  nil,
				Volume:     1,
				IssueCount: 1,
				Modified:   parseTime(t, "2024-12-26T16:46:35.608465-05:00"),
			},
		},
		{
			id: 1532,
			expected: &metron.SeriesList{
				ID:         1532,
				Name:       "100ᵗʰ Anniversary Special: Guardians of the Galaxy (2014)",
				YearBegan:  2014,
				YearEnded:  nil,
				Volume:     1,
				IssueCount: 1,
				Modified:   parseTime(t, "2024-12-26T16:46:49.163869-05:00"),
			},
		},
	}
}
