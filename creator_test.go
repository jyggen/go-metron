package metron_test

import (
	"testing"

	"github.com/jyggen/go-metron"
)

func TestCreatorByID(t *testing.T) {
	t.Parallel()
	testByID(t, "creator", (*metron.Client).CreatorByID, []testCase[metron.Creator]{
		{
			id: 5958,
			expected: metron.Creator{
				ID:          5958,
				Name:        "Frank Godwin",
				Birth:       asReference(parseDate(t, "1889-10-20")),
				Death:       asReference(parseDate(t, "1959-08-05")),
				Description: asReference("An American illustrator and comic strip artist, notable for his strip Connie and his book illustrations for Treasure Island, Kidnapped, Robinson Crusoe, Robin Hood and King Arthur. He also was a prolific editorial and advertising illustrator."),
				ImageURL:    asReference(parseURL(t, "https://static.metron.cloud/media/creator/2022/03/06/f-godwin.jpg")),
				Alias: asReference([]string{
					"Francis Godwin",
				}),
				ComicVineID: nil,
				ResourceURL: parseURL(t, "https://metron.cloud/creator/frank-godwin/"),
				Modified:    parseTime(t, "2022-03-06T10:43:26.423424-05:00"),
			},
		},
		{
			id: 11237,
			expected: metron.Creator{
				ID:          11237,
				Name:        "A. C. Farley",
				Birth:       nil,
				Death:       nil,
				Description: asReference(""),
				ImageURL:    nil,
				Alias:       nil,
				ComicVineID: asReference(49964),
				ResourceURL: parseURL(t, "https://metron.cloud/creator/a-c-farley/"),
				Modified:    parseTime(t, "2024-05-18T11:25:36.333655-04:00"),
			},
		},
	})
}

func TestCreators(t *testing.T) {
	t.Parallel()
	testList(t, "creator", (*metron.Client).Creators, []testCase[metron.CreatorList]{
		{
			id: 5466,
			expected: metron.CreatorList{
				ID:       5466,
				Name:     "A D'Amico",
				Modified: parseTime(t, "2024-02-25T15:53:40.615428-05:00"),
			},
		},
		{
			id: 7174,
			expected: metron.CreatorList{
				ID:       7174,
				Name:     "A Larger World Studios",
				Modified: parseTime(t, "2023-02-21T11:29:18.690700-05:00"),
			},
		},
		{
			id: 5848,
			expected: metron.CreatorList{
				ID:       5848,
				Name:     "A. A. Rubin",
				Modified: parseTime(t, "2021-12-29T10:10:02.832069-05:00"),
			},
		},
		{
			id: 11237,
			expected: metron.CreatorList{
				ID:       11237,
				Name:     "A. C. Farley",
				Modified: parseTime(t, "2024-05-18T11:25:36.333655-04:00"),
			},
		},
	})
}
