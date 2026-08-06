package metron_test

import (
	"testing"

	"github.com/jyggen/go-metron"
)

func TestPublisherByID(t *testing.T) {
	t.Parallel()
	testByID(t, "publisher", (*metron.Client).PublisherByID, []testCase[*metron.Publisher]{
		{
			id: 1,
			expected: &metron.Publisher{
				ID:          1,
				Name:        "Marvel",
				Founded:     new(1939),
				CountryCode: new("US"),
				Description: new(
					"Marvel Comics is the brand name and primary imprint of Marvel Worldwide Inc., formerly Marvel Publishing, Inc. and Marvel Comics Group, a publisher of American comic books and related media. In 2009, The Walt Disney Company acquired Marvel Entertainment, Marvel Worldwide's parent company.\r\n\r\nMarvel started in 1939 as Timely Publications, and by the early 1950s, had generally become known as Atlas Comics. The Marvel branding began in 1961, the year that the company launched The Fantastic Four and other superhero titles created by Steve Ditko, Stan Lee, Jack Kirby and many others.",
				),
				ImageURL: new(
					parseURL(t, "https://static.metron.cloud/media/publisher/2018/11/11/marvel.jpg"),
				),
				ComicVineID:           new(31),
				GrandComicsDatabaseID: new(78),
				ResourceURL:           parseURL(t, "https://metron.cloud/publisher/marvel/"),
				Modified:              parseTime(t, "2025-01-28T09:35:05.804360-05:00"),
			},
		},
		{
			id: 29,
			expected: &metron.Publisher{
				ID:          29,
				Name:        "12-Gauge Comics",
				Founded:     new(2004),
				CountryCode: new("US"),
				Description: new(
					"Establishing itself from day one as a company dedicated to creating quality material, 12-Gauge has continued to strengthen its reputation with each new comic series. With a focus on crime and action stories, 12-Gauge has carved out a special place in the industry, now in its 15th year of creating stories with a tradition of excellence that can’t be denied.",
				),
				ImageURL: new(
					parseURL(t, "https://static.metron.cloud/media/publisher/2019/11/03/12gauge.jpg"),
				),
				ComicVineID: new(2490),
				ResourceURL: parseURL(t, "https://metron.cloud/publisher/12-gauge-comics/"),
				Modified:    parseTime(t, "2025-04-04T23:08:02.258460-04:00"),
			},
		},
	})
}

func TestPublishers(t *testing.T) {
	t.Parallel()
	testList(t, "publisher", (*metron.Client).Publishers, []testCase[*metron.PublisherList]{
		{
			id: 178,
			expected: &metron.PublisherList{
				ID:       178,
				Name:     "10 Ton Press",
				Modified: parseTime(t, "2026-02-11T21:43:43.314351-05:00"),
			},
		},
		{
			id: 29,
			expected: &metron.PublisherList{
				ID:       29,
				Name:     "12-Gauge Comics",
				Modified: parseTime(t, "2025-04-04T23:08:02.258460-04:00"),
			},
		},
		{
			id: 218,
			expected: &metron.PublisherList{
				ID:       218,
				Name:     "3 Finger Prints",
				Modified: parseTime(t, "2026-07-21T12:49:17.396186-04:00"),
			},
		},
		{
			id: 148,
			expected: &metron.PublisherList{
				ID:       148,
				Name:     "Aaaargh! Comics",
				Modified: parseTime(t, "2025-09-24T13:52:41.776649-04:00"),
			},
		},
	})
}
