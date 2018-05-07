package availability_test

import (
	"net/http"
	"testing"

	"github.com/kamilsk/check/errors"
	"github.com/kamilsk/check/http/availability"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
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
		{
			"bad site",
			[]string{":bad"},
			func() *availability.Report { return availability.NewReport() },
			func() <-chan availability.Site {
				ch := make(chan availability.Site, 1)
				ch <- *availability.NewSite(":bad")
				close(ch)
				return ch
			},
		},
		{
			"without crawler",
			[]string{"http://test.dev/"},
			func() *availability.Report { return availability.NewReport() },
			func() <-chan availability.Site {
				ch := make(chan availability.Site, 1)
				site := *availability.NewSite("http://test.dev/")
				site.Error = errors.Simple("crawler is not provided")
				ch <- site
				close(ch)
				return ch
			},
		},
		{
			"normal case",
			[]string{"http://test.dev/"},
			func() *availability.Report {
				crawler := &CrawlerMock{shift: func(to availability.EventBus) {
					to <- availability.ResponseEvent{StatusCode: http.StatusOK, Location: "http://test.dev/"}
					to <- availability.WalkEvent{Page: "http://test.dev/", Href: "http://accepted.dev/"}
					to <- availability.WalkEvent{Page: "http://test.dev/", Href: "http://redirect.dev/"}
					to <- availability.WalkEvent{Page: "http://test.dev/", Href: "http://noaccess.dev/"}
					to <- availability.ProblemEvent{Message: "bad url", Context: ":bad"}
					to <- availability.ResponseEvent{StatusCode: http.StatusAccepted, Location: "http://accepted.dev/"}
					to <- availability.ErrorEvent{StatusCode: http.StatusFound,
						Location: "http://redirect.dev/", Redirect: "https://redirect.dev/"}
					to <- availability.ErrorEvent{StatusCode: http.StatusForbidden, Location: "http://noaccess.dev/"}
					close(to)
				}}
				crawler.On("Visit", "http://test.dev/", mock.Anything).Return(nil)
				report := availability.NewReport(availability.CrawlerForSites(crawler))
				return report
			},
			func() <-chan availability.Site {
				ch := make(chan availability.Site, 1)
				ch <- *availability.NewSite("http://test.dev/")
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
			for i := range expected {
				assert.Equal(t, expected[i].Name, obtained[i].Name)
				if expected[i].Error != nil {
					assert.EqualError(t, obtained[i].Error, expected[i].Error.Error())
				} else {
					assert.NoError(t, obtained[i].Error)
				}
			}
		})
	}
}

func TestReporter_handlePanic(t *testing.T) {
	tests := []struct {
		name     string
		rawURLs  []string
		reporter func() *availability.Report
		expected string
	}{
		{
			"unexpected event",
			[]string{"http://test.dev/"},
			func() *availability.Report {
				crawler := &CrawlerMock{shift: func(to availability.EventBus) {
					type unknown struct{ availability.ProblemEvent }
					to <- unknown{availability.ProblemEvent{Message: "bad url", Context: ":bad"}}
					close(to)
				}}
				crawler.On("Visit", "http://test.dev/", mock.Anything).Return(nil)
				report := availability.NewReport(availability.CrawlerForSites(crawler))
				return report
			},
			"panic: unexpected event type availability_test.unknown",
		},
		{
			"not consistent fetch result",
			[]string{"http://test.dev/"},
			func() *availability.Report {
				crawler := &CrawlerMock{shift: func(to availability.EventBus) {
					to <- availability.ResponseEvent{StatusCode: http.StatusOK, Location: "http://test.dev/"}
					to <- availability.WalkEvent{Page: "http://test.dev/without-response/", Href: "http://test.dev/"}
					close(to)
				}}
				crawler.On("Visit", "http://test.dev/", mock.Anything).Return(nil)
				report := availability.NewReport(availability.CrawlerForSites(crawler))
				return report
			},
			"runtime error: invalid memory address or nil pointer dereference",
		},
	}
	for _, test := range tests {
		tc := test
		t.Run(test.name, func(t *testing.T) {
			assert.Panics(t, func() {
				defer func() {
					if r := recover(); r != nil {
						err, is := r.(error)
						assert.True(t, is)
						assert.EqualError(t, err, tc.expected)
						panic(r)
					}
				}()
				tc.reporter().For(tc.rawURLs).Fill()
			})
		})
	}
}
