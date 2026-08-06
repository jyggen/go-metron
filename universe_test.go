package metron_test

import (
	"testing"

	"github.com/jyggen/go-metron"
)

func TestUniverseByID(t *testing.T) {
	t.Parallel()
	testByID(t, "universe", (*metron.Client).UniverseByID, []testCase[*metron.Universe]{
		{
			id: 24,
			expected: &metron.Universe{
				ID: 24,
				Publisher: metron.Reference{
					ID:   2,
					Name: "DC Comics",
				},
				Name:        "ABC",
				Designation: "Earth 25",
				Description: new("Home to the characters from Alan Moore's America's Best Comics imprint."),
				ImageURL: new(
					parseURL(t, "https://static.metron.cloud/media/universe/2024/01/25/earth-abc.webp"),
				),
				ResourceURL: parseURL(t, "https://metron.cloud/universe/abc/"),
				Modified:    parseTime(t, "2024-01-25T09:24:33.271598-05:00"),
			},
		},
	})
}

func TestUniverses(t *testing.T) {
	t.Parallel()
	testList(t, "universe", (*metron.Client).Universes, []testCase[*metron.UniverseList]{
		{
			id: 153,
			expected: &metron.UniverseList{
				ID:       153,
				Name:     "2099 AD - Marvel Knights",
				Modified: parseTime(t, "2025-01-30T14:13:19.980781-05:00"),
			},
		},
		{
			id: 24,
			expected: &metron.UniverseList{
				ID:       24,
				Name:     "ABC",
				Modified: parseTime(t, "2024-01-25T09:24:33.271598-05:00"),
			},
		},
		{
			id: 157,
			expected: &metron.UniverseList{
				ID:       157,
				Name:     "Absolute Universe",
				Modified: parseTime(t, "2025-04-16T09:21:47.838281-04:00"),
			},
		},
		{
			id: 110,
			expected: &metron.UniverseList{
				ID:       110,
				Name:     "Age of Apocalypse",
				Modified: parseTime(t, "2024-04-17T10:31:52.624470-04:00"),
			},
		},
	})
}
