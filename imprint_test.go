package metron_test

import (
	"testing"

	"github.com/jyggen/go-metron"
)

func TestImprintByID(t *testing.T) {
	t.Parallel()
	testByID(t, "imprint", (*metron.Client).ImprintByID, []testCase[metron.Imprint]{
		{
			id: 0,
			expected: metron.Imprint{
				ID:          0,
				Name:        "Foobar",
				Founded:     nil,
				Description: asReference(""),
				ImageURL:    nil,
				Publisher:   metron.Reference{ID: 0, Name: "Foobar"},
				ComicVineID: asReference(0),
				ResourceURL: parseURL(t, "https://example.com/"),
				Modified:    parseTime(t, "1970-01-01T00:00:00.000000Z"),
			},
		},
		{
			id: 1,
			expected: metron.Imprint{
				ID:          1,
				Name:        "Vertigo Comics",
				Founded:     asReference(1993),
				Description: asReference("Vertigo Comics (also known as DC Vertigo) is an imprint of the American comic book publisher DC Comics. It was created in 1993 to publish stories with more graphic or adult content that could not fit within the restrictions of the Comics Code Authority, thus allowing more creative freedom than DC's main imprint. These comics were free to contain explicit violence, substance and drug abuse, sexuality, nudity, profanity, and other controversial subjects, similar to the content of R-rated films.\r\n\r\nAlthough its initial publications were primarily in the horror and fantasy genres, it has also published works dealing with crime, social commentary, speculative fiction, biography, and other genres. Originally publishing a mix of company- and creator-owned work, its current focus is on the latter. It pioneered in North America an increasingly common publishing model, in which monthly series are periodically comprised into collected editions which are kept in print for bookstore sale.\r\n\r\nVertigo series have won the comics industry's Eisner Award, including the \"best continuing series\" of various years (The Sandman, Preacher, 100 Bullets, Y: The Last Man and Fables). Several of its publications have been adapted to film (such as Constantine, A History of Violence, Stardust, and V for Vendetta) and episodic television (such as Constantine, iZombie, Lucifer, and Preacher)."),
				ImageURL:    asReference(parseURL(t, "https://static.metron.cloud/media/imprint/2024/08/12/vertigo.jpg")),
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
	testList(t, "imprint", (*metron.Client).Imprints, []testCase[metron.ImprintList]{
		{
			id: 14,
			expected: metron.ImprintList{
				ID:       14,
				Name:     "Amalgam Comics",
				Modified: parseTime(t, "2024-08-18T10:09:51.151912-04:00"),
			},
		},
		{
			id: 4,
			expected: metron.ImprintList{
				ID:       4,
				Name:     "Archie Horror",
				Modified: parseTime(t, "2024-08-12T15:07:30.049388-04:00"),
			},
		},
		{
			id: 13,
			expected: metron.ImprintList{
				ID:       13,
				Name:     "Boom! Box",
				Modified: parseTime(t, "2024-08-13T15:43:21.628316-04:00"),
			},
		},
		{
			id: 2,
			expected: metron.ImprintList{
				ID:       2,
				Name:     "DC Black Label",
				Modified: parseTime(t, "2024-08-12T12:14:36.732542-04:00"),
			},
		},
	})
}
