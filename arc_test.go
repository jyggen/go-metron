package metron_test

import (
	"testing"

	"github.com/jyggen/go-metron"
)

func TestArcByID(t *testing.T) {
	t.Parallel()
	testByID(t, "arc", (*metron.Client).ArcByID, []testCase[metron.Arc]{
		{
			id: 659,
			expected: metron.Arc{
				ID:          659,
				Name:        "'Til Death Do Us...",
				Description: asReference(""),
				ImageURL:    nil,
				ComicVineID: nil,
				ResourceURL: parseURL(t, "https://metron.cloud/arc/til-death-do-us/"),
				Modified:    parseTime(t, "2022-01-16T10:12:00.426525-05:00"),
			},
		},
		{
			id: 1493,
			expected: metron.Arc{
				ID:          1493,
				Name:        "Crisis on Infinite Darkwings",
				Description: asReference(""),
				ImageURL:    asReference(parseURL(t, "https://static.metron.cloud/media/arc/2024/03/07/3069ed71448e4cf3bdca8f808fc596ce.jpg")),
				ComicVineID: asReference(56259),
				ResourceURL: parseURL(t, "https://metron.cloud/arc/darkwing-duck-crisis-on-infinite-darkwings/"),
				Modified:    parseTime(t, "2024-03-07T12:36:10.633143-05:00"),
			},
		},
	})
}

func TestArcs(t *testing.T) {
	t.Parallel()
	testList(t, "arc", (*metron.Client).Arcs, []testCase[metron.ArcList]{
		{
			id: 659,
			expected: metron.ArcList{
				ID:       659,
				Name:     "'Til Death Do Us...",
				Modified: parseTime(t, "2022-01-16T10:12:00.426525-05:00"),
			},
		},
		{
			id: 931,
			expected: metron.ArcList{
				ID:       931,
				Name:     "(She) Drunk History",
				Modified: parseTime(t, "2023-02-15T11:47:59.483664-05:00"),
			},
		},
		{
			id: 871,
			expected: metron.ArcList{
				ID:       871,
				Name:     "1+2 = Fantastic Three",
				Modified: parseTime(t, "2022-09-08T09:53:30.626809-04:00"),
			},
		},
		{
			id: 1031,
			expected: metron.ArcList{
				ID:       1031,
				Name:     "1602",
				Modified: parseTime(t, "2023-03-12T15:41:46.220978-04:00"),
			},
		},
	})
}
