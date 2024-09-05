package metron_test

import (
	"testing"

	"github.com/jyggen/go-metron"
)

func TestSeriesByID(t *testing.T) {
	t.Parallel()
	testByID(t, "series", (*metron.Client).SeriesByID, []testCase[metron.Series]{
		{
			id: 3371,
			expected: metron.Series{
				ID:       3371,
				Name:     "Batman 2022 Annual",
				SortName: "Batman 2022 Annual",
				Volume:   1,
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
				Description: asReference(""),
				IssueCount:  1,
				Genres: []metron.Reference{
					{
						ID:   10,
						Name: "Super-Hero",
					},
				},
				Associated: []struct {
					ID   int    `json:"id"`
					Name string `json:"series"`
				}{
					{
						ID:   93,
						Name: "Batman (2016)",
					},
				},
				ComicVineID: asReference(143255),
				ResourceURL: parseURL(t, "https://metron.cloud/series/batman-2022-annual-2022/"),
				Modified:    parseTime(t, "2023-10-14T15:41:45.214400-04:00"),
			},
		},
		{
			id: 793,
			expected: metron.Series{
				ID:       793,
				Name:     "Fables",
				SortName: "Fables",
				Volume:   1,
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
				YearBegan:   2002,
				YearEnded:   asReference(2015),
				Description: asReference("When the Adversary conquered the lands of legends, the inhabitants were forced into exile. They form a secret society, a hidden enclave in modern-day New York. Bill Willingham's award-winning \"fractured-fairy-tale\" series explores the world of these beloved fables...one that exists within our own."),
				IssueCount:  149,
				Genres:      []metron.Reference{},
				Associated: []struct {
					ID   int    `json:"id"`
					Name string `json:"series"`
				}{
					{
						ID:   3396,
						Name: "Fables (2022)",
					},
				},
				ComicVineID: asReference(9723),
				ResourceURL: parseURL(t, "https://metron.cloud/series/fables-2002/"),
				Modified:    parseTime(t, "2024-06-20T12:32:08.926341-04:00"),
			},
		},
	})
}

func TestSeriesByPublisherID(t *testing.T) {
	t.Parallel()
	testListByID(t, "publisher", 1, "series", (*metron.Client).SeriesByPublisherID, seriesListTestCases(t))
}

func TestSeries(t *testing.T) {
	t.Parallel()
	testList(t, "series", (*metron.Client).Series, seriesListTestCases(t))
}

func seriesListTestCases(t *testing.T) []testCase[metron.SeriesList] {
	return []testCase[metron.SeriesList]{
		{
			id: 6227,
			expected: metron.SeriesList{
				ID:         6227,
				Name:       "'68 (2006)",
				YearBegan:  2006,
				Volume:     1,
				IssueCount: 1,
				Modified:   parseTime(t, "2023-12-02T10:10:47.573380-05:00"),
			},
		},
		{
			id: 6228,
			expected: metron.SeriesList{
				ID:         6228,
				Name:       "'68 (2011)",
				YearBegan:  2011,
				Volume:     2,
				IssueCount: 4,
				Modified:   parseTime(t, "2023-12-02T10:10:18.550260-05:00"),
			},
		},
		{
			id: 6229,
			expected: metron.SeriesList{
				ID:         6229,
				Name:       "'68 Hallowed Ground (2013)",
				YearBegan:  2013,
				Volume:     1,
				IssueCount: 1,
				Modified:   parseTime(t, "2023-12-02T10:13:50.130689-05:00"),
			},
		},
		{
			id: 6236,
			expected: metron.SeriesList{
				ID:         6236,
				Name:       "'68 Hardship (2011)",
				YearBegan:  2011,
				Volume:     1,
				IssueCount: 1,
				Modified:   parseTime(t, "2023-12-02T10:53:11.674254-05:00"),
			},
		},
	}
}
