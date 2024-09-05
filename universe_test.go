package metron_test

import (
	"github.com/jyggen/go-metron"
	"testing"
)

func TestUniverseByID(t *testing.T) {
	t.Parallel()
	testByID(t, "universe", (*metron.Client).UniverseByID, []testCase[metron.Universe]{
		{
			id: 24,
			expected: metron.Universe{
				ID: 24,
				Publisher: metron.Reference{
					ID:   2,
					Name: "DC Comics",
				},
				Name:        "ABC",
				Designation: asReference("Earth 25"),
				Description: asReference("Home to the characters from Alan Moore's America's Best Comics imprint."),
				ImageURL:    asReference(parseURL(t, "https://static.metron.cloud/media/universe/2024/01/25/earth-abc.webp")),
				ResourceURL: parseURL(t, "https://metron.cloud/universe/abc/"),
				Modified:    parseTime(t, "2024-01-25T09:24:33.271598-05:00"),
			},
		},
	})
}

func TestUniverses(t *testing.T) {
	t.Parallel()
	testList(t, "universe", (*metron.Client).Universes, []testCase[metron.UniverseList]{
		{
			id: 24,
			expected: metron.UniverseList{
				ID:       24,
				Name:     "ABC",
				Modified: parseTime(t, "2024-01-25T09:24:33.271598-05:00"),
			},
		},
		{
			id: 110,
			expected: metron.UniverseList{
				ID:       110,
				Name:     "Age of Apocalypse",
				Modified: parseTime(t, "2024-04-17T10:31:52.624470-04:00"),
			},
		},
		{
			id: 69,
			expected: metron.UniverseList{
				ID:       69,
				Name:     "Amalgam",
				Modified: parseTime(t, "2024-01-25T09:22:08.422067-05:00"),
			},
		},
		{
			id: 75,
			expected: metron.UniverseList{
				ID:       75,
				Name:     "Amazonia",
				Modified: parseTime(t, "2024-01-25T09:23:39.686552-05:00"),
			},
		},
	})
}
