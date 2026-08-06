package metron_test

import (
	"testing"

	"github.com/jyggen/go-metron"
)

func TestTeamByID(t *testing.T) {
	t.Parallel()
	testByID(t, "team", (*metron.Client).TeamByID, []testCase[*metron.Team]{
		{
			id: 180,
			expected: &metron.Team{
				ID:          180,
				Name:        "A-Force",
				Description: new("Marvel's first all-female team of Avengers."),
				ImageURL:    new(parseURL(t, "https://static.metron.cloud/media/team/2019/07/27/aforce.jpg")),
				Creators:    []metron.CreatorList{},
				Universes:   []metron.UniverseList{},
				ComicVineID: nil,
				ResourceURL: parseURL(t, "https://metron.cloud/team/a-force/"),
				Modified:    parseTime(t, "2019-07-27T16:52:41.547904-04:00"),
			},
		},
		{
			id: 930,
			expected: &metron.Team{
				ID:          930,
				Name:        "Infinites",
				Description: new(""),
				ImageURL: new(
					parseURL(
						t,
						"https://static.metron.cloud/media/team/2023/03/06/14bda442d94345548e47683884914484.jpg",
					),
				),
				Creators: []metron.CreatorList{},
				Universes: []metron.UniverseList{
					{
						ID:       110,
						Name:     "Age of Apocalypse",
						Modified: parseTime(t, "2024-04-17T10:31:52.624470-04:00"),
					},
				},
				ComicVineID: new(41013),
				ResourceURL: parseURL(t, "https://metron.cloud/team/infinites/"),
				Modified:    parseTime(t, "2025-02-19T01:42:38.767067-05:00"),
			},
		},
	})
}

func TestTeams(t *testing.T) {
	t.Parallel()
	testList(t, "team", (*metron.Client).Teams, []testCase[*metron.TeamList]{
		{
			id: 2225,
			expected: &metron.TeamList{
				ID:       2225,
				Name:     "13th Floor Witches",
				Modified: parseTime(t, "2025-09-10T09:28:21.026379-04:00"),
			},
		},
		{
			id: 2772,
			expected: &metron.TeamList{
				ID:       2772,
				Name:     "181st Imperial Fighter Wing",
				Modified: parseTime(t, "2026-01-06T19:28:27.534918-05:00"),
			},
		},
		{
			id: 1806,
			expected: &metron.TeamList{
				ID:       1806,
				Name:     "22 Brides",
				Modified: parseTime(t, "2025-03-16T09:32:56.095096-04:00"),
			},
		},
		{
			id: 2711,
			expected: &metron.TeamList{
				ID:       2711,
				Name:     "327th Star Corps",
				Modified: parseTime(t, "2026-01-04T13:32:44.102791-05:00"),
			},
		},
	})
}
