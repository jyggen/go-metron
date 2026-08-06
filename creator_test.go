package metron_test

import (
	"testing"

	"github.com/jyggen/go-metron"
)

func TestCreatorByID(t *testing.T) {
	t.Parallel()
	testByID(t, "creator", (*metron.Client).CreatorByID, []testCase[*metron.Creator]{
		{
			id: 5958,
			expected: &metron.Creator{
				ID:    5958,
				Name:  "Frank Godwin",
				Birth: new(parseDate(t, "1889-10-20")),
				Death: new(parseDate(t, "1959-08-05")),
				Description: new(
					"An American illustrator and comic strip artist, notable for his strip Connie and his book illustrations for Treasure Island, Kidnapped, Robinson Crusoe, Robin Hood and King Arthur. He also was a prolific editorial and advertising illustrator.",
				),
				ImageURL: new(
					parseURL(t, "https://static.metron.cloud/media/creator/2022/03/06/f-godwin.jpg"),
				),
				Alias: new([]string{
					"Francis Godwin",
				}),
				ComicVineID: nil,
				ResourceURL: parseURL(t, "https://metron.cloud/creator/frank-godwin/"),
				Modified:    parseTime(t, "2022-03-06T10:43:26.423424-05:00"),
			},
		},
		{
			id: 11237,
			expected: &metron.Creator{
				ID:          11237,
				Name:        "A. C. Farley",
				Birth:       nil,
				Death:       nil,
				Description: new(""),
				ImageURL:    nil,
				Alias:       new([]string{}),
				ComicVineID: new(49964),
				ResourceURL: parseURL(t, "https://metron.cloud/creator/a-c-farley/"),
				Modified:    parseTime(t, "2024-05-18T11:25:36.333655-04:00"),
			},
		},
	})
}

func TestCreators(t *testing.T) {
	t.Parallel()
	testList(t, "creator", (*metron.Client).Creators, []testCase[*metron.CreatorList]{
		{
			id: 13909,
			expected: &metron.CreatorList{
				ID:       13909,
				Name:     "Aadi Salman",
				Modified: parseTime(t, "2025-03-28T14:24:09.691934-04:00"),
			},
		},
		{
			id: 10900,
			expected: &metron.CreatorList{
				ID:       10900,
				Name:     "A.A. Milne",
				Modified: parseTime(t, "2024-03-07T11:02:46.113441-05:00"),
			},
		},
		{
			id: 3048,
			expected: &metron.CreatorList{
				ID:       3048,
				Name:     "Aaron Alexovich",
				Modified: parseTime(t, "2025-02-18T12:35:42.128935-05:00"),
			},
		},
		{
			id: 17721,
			expected: &metron.CreatorList{
				ID:       17721,
				Name:     "Aaron Allen",
				Modified: parseTime(t, "2026-03-12T16:18:11.116741-04:00"),
			},
		},
	})
}
