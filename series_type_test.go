package metron_test

import (
	"testing"

	"github.com/jyggen/go-metron"
)

func TestSeriesTypes(t *testing.T) {
	t.Parallel()
	testList(t, "series_type", (*metron.Client).SeriesTypes, []testCase[metron.SeriesTypeList]{
		{
			id: 6,
			expected: metron.SeriesTypeList{
				ID:   6,
				Name: "Annual",
			},
		},
		{
			id: 12,
			expected: metron.SeriesTypeList{
				ID:   12,
				Name: "Digital Chapter",
			},
		},
		{
			id: 9,
			expected: metron.SeriesTypeList{
				ID:   9,
				Name: "Graphic Novel",
			},
		},
		{
			id: 8,
			expected: metron.SeriesTypeList{
				ID:   8,
				Name: "Hardcover",
			},
		},
	})
}
