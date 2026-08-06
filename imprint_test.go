package metron_test

import (
	"testing"

	"github.com/jyggen/go-metron"
)

func TestImprintByID(t *testing.T) {
	t.Parallel()
	testByID(t, "imprint", (*metron.Client).ImprintByID, []testCase[*metron.Imprint]{
		{
			id: 0,
			expected: &metron.Imprint{
				ID:          0,
				Name:        "Foobar",
				Founded:     nil,
				Description: new(""),
				ImageURL:    nil,
				Publisher:   metron.Reference{ID: 0, Name: "Foobar"},
				ComicVineID: new(0),
				ResourceURL: parseURL(t, "https://example.com/"),
				Modified:    parseTime(t, "1970-01-01T00:00:00.000000Z"),
			},
		},
		{
			id: 1,
			expected: &metron.Imprint{
				ID:      1,
				Name:    "Vertigo Comics",
				Founded: new(1993),
				Description: new(
					"Vertigo Comics (also known as DC Vertigo) is an imprint of the American comic book publisher DC Comics. It was created in 1993 to publish stories with more graphic or adult content that could not fit within the restrictions of the Comics Code Authority, thus allowing more creative freedom than DC's main imprint. These comics were free to contain explicit violence, substance and drug abuse, sexuality, nudity, profanity, and other controversial subjects, similar to the content of R-rated films.\r\n\r\nAlthough its initial publications were primarily in the horror and fantasy genres, it has also published works dealing with crime, social commentary, speculative fiction, biography, and other genres. Originally publishing a mix of company- and creator-owned work, its current focus is on the latter. It pioneered in North America an increasingly common publishing model, in which monthly series are periodically comprised into collected editions which are kept in print for bookstore sale.\r\n\r\nVertigo series have won the comics industry's Eisner Award, including the \"best continuing series\" of various years (The Sandman, Preacher, 100 Bullets, Y: The Last Man and Fables). Several of its publications have been adapted to film (such as Constantine, A History of Violence, Stardust, and V for Vendetta) and episodic television (such as Constantine, iZombie, Lucifer, and Preacher).",
				),
				ImageURL: new(
					parseURL(t, "https://static.metron.cloud/media/imprint/2024/08/12/vertigo.jpg"),
				),
				Publisher:   metron.Reference{ID: 2, Name: "DC Comics"},
				ComicVineID: nil,
				ResourceURL: parseURL(t, "https://metron.cloud/imprint/vertigo-comics/"),
				Modified:    parseTime(t, "2024-08-12T12:13:54.087728-04:00"),
			},
		},
	})
}

func TestImprints(t *testing.T) {
	t.Parallel()
	testList(t, "imprint", (*metron.Client).Imprints, []testCase[*metron.ImprintList]{
		{
			id: 31,
			expected: &metron.ImprintList{
				ID:       31,
				Name:     "Abrams Fanfare",
				Modified: parseTime(t, "2025-04-14T19:35:22.359198-04:00"),
			},
		},
		{
			id: 32,
			expected: &metron.ImprintList{
				ID:       32,
				Name:     "Action Lab Danger Zone",
				Modified: parseTime(t, "2025-09-12T11:31:32.032705-04:00"),
			},
		},
		{
			id: 42,
			expected: &metron.ImprintList{
				ID:       42,
				Name:     "Adventure Comics",
				Modified: parseTime(t, "2026-05-21T12:59:20.248379-04:00"),
			},
		},
		{
			id: 29,
			expected: &metron.ImprintList{
				ID:       29,
				Name:     "Aircel Publishing",
				Modified: parseTime(t, "2026-03-01T20:26:36.804924-05:00"),
			},
		},
	})
}
