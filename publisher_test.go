package metron_test

import (
	"testing"

	"github.com/jyggen/go-metron"
)

func TestPublisherByID(t *testing.T) {
	t.Parallel()
	testByID(t, "publisher", (*metron.Client).PublisherByID, []testCase[metron.Publisher]{
		{
			id: 1,
			expected: metron.Publisher{
				ID:          1,
				Name:        "Marvel",
				Founded:     asReference(1939),
				Description: asReference("Marvel Comics is the brand name and primary imprint of Marvel Worldwide Inc., formerly Marvel Publishing, Inc. and Marvel Comics Group, a publisher of American comic books and related media. In 2009, The Walt Disney Company acquired Marvel Entertainment, Marvel Worldwide's parent company.\r\n\r\nMarvel started in 1939 as Timely Publications, and by the early 1950s, had generally become known as Atlas Comics. The Marvel branding began in 1961, the year that the company launched The Fantastic Four and other superhero titles created by Steve Ditko, Stan Lee, Jack Kirby and many others."),
				ImageURL:    asReference(parseURL(t, "https://static.metron.cloud/media/publisher/2018/11/11/marvel.jpg")),
				ComicVineID: asReference(31),
				ResourceURL: parseURL(t, "https://metron.cloud/publisher/marvel/"),
				Modified:    parseTime(t, "2024-04-07T04:53:45.729670-04:00"),
			},
		},
		{
			id: 29,
			expected: metron.Publisher{
				ID:          29,
				Name:        "12-Gauge Comics",
				Founded:     asReference(2004),
				Description: asReference("Establishing itself from day one as a company dedicated to creating quality material, 12-Gauge has continued to strengthen its reputation with each new comic series. With a focus on crime and action stories, 12-Gauge has carved out a special place in the industry, now in its 15th year of creating stories with a tradition of excellence that can’t be denied."),
				ImageURL:    asReference(parseURL(t, "https://static.metron.cloud/media/publisher/2019/11/03/12gauge.jpg")),
				ComicVineID: nil,
				ResourceURL: parseURL(t, "https://metron.cloud/publisher/12-gauge-comics/"),
				Modified:    parseTime(t, "2019-11-03T17:16:03.835821-05:00"),
			},
		},
	})
}

func TestPublishers(t *testing.T) {
	t.Parallel()
	testList(t, "publisher", (*metron.Client).Publishers, []testCase[metron.PublisherList]{
		{
			id: 29,
			expected: metron.PublisherList{
				ID:       29,
				Name:     "12-Gauge Comics",
				Modified: parseTime(t, "2019-11-03T17:16:03.835821-05:00"),
			},
		},
		{
			id: 64,
			expected: metron.PublisherList{
				ID:       64,
				Name:     "AAA Pop Comics",
				Modified: parseTime(t, "2023-07-24T08:21:06.019465-04:00"),
			},
		},
		{
			id: 36,
			expected: metron.PublisherList{
				ID:       36,
				Name:     "AWA Studios",
				Modified: parseTime(t, "2020-07-06T21:36:46.414152-04:00"),
			},
		},
		{
			id: 56,
			expected: metron.PublisherList{
				ID:       56,
				Name:     "Aardvark-Vanaheim",
				Modified: parseTime(t, "2023-04-29T16:56:46.004672-04:00"),
			},
		},
	})
}
