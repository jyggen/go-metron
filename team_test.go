package metron_test

import (
	"testing"

	"github.com/jyggen/go-metron"
)

func TestTeamByID(t *testing.T) {
	t.Parallel()
	testByID(t, "team", (*metron.Client).TeamByID, []testCase[metron.Team]{
		{
			id: 180,
			expected: metron.Team{
				ID:          180,
				Name:        "A-Force",
				Description: asReference("Marvel's first all-female team of Avengers."),
				ImageURL:    asReference(parseURL(t, "https://static.metron.cloud/media/team/2019/07/27/aforce.jpg")),
				Creators:    []metron.CreatorList{},
				Universes:   []metron.UniverseList{},
				ComicVineID: nil,
				ResourceURL: parseURL(t, "https://metron.cloud/team/a-force/"),
				Modified:    parseTime(t, "2019-07-27T16:52:41.547904-04:00"),
			},
		},
		{
			id: 930,
			expected: metron.Team{
				ID:          930,
				Name:        "Infinites",
				Description: asReference(""),
				ImageURL:    asReference(parseURL(t, "https://static.metron.cloud/media/team/2023/03/06/14bda442d94345548e47683884914484.jpg")),
				Creators:    []metron.CreatorList{},
				Universes: []metron.UniverseList{
					{
						ID:       110,
						Name:     "Age of Apocalypse",
						Modified: parseTime(t, "2024-04-17T10:31:52.624470-04:00"),
					},
				},
				ComicVineID: nil,
				ResourceURL: parseURL(t, "https://metron.cloud/team/infinites/"),
				Modified:    parseTime(t, "2024-08-06T09:23:45.515089-04:00"),
			},
		},
	})
}

func TestTeams(t *testing.T) {
	t.Parallel()
	testList(t, "team", (*metron.Client).Teams, []testCase[metron.TeamList]{
		{
			id: 1089,
			expected: metron.TeamList{
				ID:       1089,
				Name:     "501st Legion",
				Modified: parseTime(t, "2023-04-04T10:12:13.619621-04:00"),
			},
		},
		{
			id: 1498,
			expected: metron.TeamList{
				ID:       1498,
				Name:     "5th Dimensional Imps",
				Modified: parseTime(t, "2024-03-16T09:28:28.725699-04:00"),
			},
		},
		{
			id: 180,
			expected: metron.TeamList{
				ID:       180,
				Name:     "A-Force",
				Modified: parseTime(t, "2019-07-27T16:52:41.547904-04:00"),
			},
		},
		{
			id: 523,
			expected: metron.TeamList{
				ID:       523,
				Name:     "A-Next",
				Modified: parseTime(t, "2021-11-17T16:14:53.029611-05:00"),
			},
		},
	})
}
