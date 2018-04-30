package availability_test

import (
	"testing"

	"github.com/kamilsk/check/http/availability"
	"github.com/stretchr/testify/assert"
)

func TestReporter(t *testing.T) {
	tests := []struct {
		name     string
		rawURLs  []string
		reporter func() *availability.Report
		expected func() <-chan availability.Site
	}{
		{
			"empty report",
			nil,
			func() *availability.Report { return availability.NewReport() },
			func() <-chan availability.Site {
				ch := make(chan availability.Site)
				close(ch)
				return ch
			},
		},
	}
	for _, test := range tests {
		tc := test
		t.Run(test.name, func(t *testing.T) {
			obtainedPipe := tc.reporter().For(tc.rawURLs).Fill().Sites()
			expectedPipe := tc.expected()
			assert.Equal(t, len(expectedPipe), len(obtainedPipe))
			obtained := make([]availability.Site, 0, len(obtainedPipe))
			for site := range obtainedPipe {
				obtained = append(obtained, site)
			}
			expected := make([]availability.Site, 0, len(expectedPipe))
			for site := range expectedPipe {
				expected = append(expected, site)
			}
			assert.Equal(t, expected, obtained)
		})
	}
}
