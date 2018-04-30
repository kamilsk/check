package availability_test

import (
	"bytes"
	"net/http"
	"testing"

	"github.com/kamilsk/check/http/availability"
	"github.com/stretchr/testify/assert"
)

func TestPrinter(t *testing.T) {
	buf := bytes.NewBuffer(nil)

	tests := []struct {
		name     string
		printer  func() *availability.Printer
		report   func() availability.Reporter
		checker  func(assert.TestingT, error, ...interface{}) bool
		expected string
	}{
		{
			"empty report",
			func() *availability.Printer { return availability.NewPrinter() },
			func() availability.Reporter { return nil },
			assert.Error,
			"",
		},
		{
			"site with error",
			func() *availability.Printer { return availability.NewPrinter(availability.OutputForPrinting(buf)) },
			func() availability.Reporter {
				m := &PrinterMock{}
				data := make(chan availability.Site, 1)
				data <- *availability.NewSite(":bad")
				close(data)
				var pipe <-chan availability.Site = data
				m.On("Sites").Return(pipe)
				return m
			},
			assert.NoError,
			`report ":bad" has error`,
		},
		{
			"site with problems",
			func() *availability.Printer { return availability.NewPrinter(availability.OutputForPrinting(buf)) },
			func() availability.Reporter {
				m := &PrinterMock{}
				data := make(chan availability.Site, 1)
				data <- availability.Site{Problems: []availability.ProblemEvent{{Message: "something happened"}}}
				close(data)
				var pipe <-chan availability.Site = data
				m.On("Sites").Return(pipe)
				return m
			},
			assert.NoError,
			"- [0] something happened `<nil>`",
		},
		{
			"normal case",
			func() *availability.Printer { return availability.NewPrinter(availability.OutputForPrinting(buf)) },
			func() availability.Reporter {
				m := &PrinterMock{}
				data := make(chan availability.Site, 1)
				data <- availability.Site{Pages: []*availability.Page{
					{
						&availability.Link{StatusCode: http.StatusOK, Location: "https://kamil.samigullin.info/en/"},
						[]*availability.Link{
							{StatusCode: http.StatusServiceUnavailable, Location: "https://github.com/kamilsk"},
							{StatusCode: http.StatusForbidden, Location: "https://www.linkedin.com/in/kamilsk"},
							{StatusCode: http.StatusFound,
								Location: "http://howilive.ru/en/", Redirect: "https://howilive.ru/en/"},
							{StatusCode: http.StatusProcessing, Location: "https://twitter.com/ikamilsk"},
						},
					},
					{
						Link: &availability.Link{StatusCode: http.StatusOK, Location: "https://kamil.samigullin.info/"},
						Links: []*availability.Link{
							{StatusCode: http.StatusServiceUnavailable, Location: "https://github.com/kamilsk"},
							{StatusCode: http.StatusForbidden, Location: "https://www.linkedin.com/in/kamilsk"},
							{StatusCode: http.StatusFound,
								Location: "http://howilive.ru/en/", Redirect: "https://howilive.ru/en/"},
							{StatusCode: http.StatusProcessing, Location: "https://twitter.com/ikamilsk"},
						},
					},
				}}
				close(data)
				var pipe <-chan availability.Site = data
				m.On("Sites").Return(pipe)
				return m
			},
			assert.NoError,
			"[200] https://kamil.samigullin.info/",
		},
	}
	for _, test := range tests {
		tc := test
		t.Run(test.name, func(t *testing.T) {
			buf.Reset()
			tc.checker(t, tc.printer().For(tc.report()).Print())
			assert.Contains(t, buf.String(), tc.expected)
		})
	}
}
