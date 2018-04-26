package availability

import (
	"net/http"
	"net/url"
	"sync"

	"github.com/gocolly/colly"
	"github.com/pkg/errors"
)

const (
	location = "Location"
)

var redirects = map[int]struct{}{
	http.StatusMovedPermanently:  {},
	http.StatusFound:             {},
	http.StatusTemporaryRedirect: {},
	http.StatusPermanentRedirect: {},
}

func CreateReport(rawURL string) *Report {
	report := NewReport(rawURL)
	report.error = report.Get()
	return report
}

func NewReport(rawURL string) *Report {
	location, err := url.Parse(rawURL)
	return &Report{
		name:     rawURL,
		error:    errors.Wrapf(err, "parse rawURL %q for report", rawURL),
		location: location,
		mu:       &sync.RWMutex{}, pages: make([]*Page, 0, 8), journal: make(map[string]*Link),
	}
}

type Reports []*Report

func (l *Reports) Fill(rawURLs []string) {
	*l = make([]*Report, len(rawURLs))

	wg, start := &sync.WaitGroup{}, make(chan struct{})
	for i, rawURL := range rawURLs {
		wg.Add(1)
		go func(idx int, rawURL string) {
			defer wg.Done()
			<-start
			(*l)[i] = CreateReport(rawURL)
		}(i, rawURL)
	}
	close(start)
	wg.Wait()
}

type Report struct {
	name     string
	error    error
	location *url.URL
	mu       *sync.RWMutex
	pages    []*Page
	journal  map[string]*Link
}

func (r *Report) Name() string { return r.name }

func (r *Report) Error() error { return r.error }

func (r *Report) Get() error {
	if r.error != nil {
		return r.error
	}
	c := colly.NewCollector(
		UserAgent(), NoRedirect(), colly.IgnoreRobotsTxt(),

		TempOption(r),
	)
	return c.Visit(r.location.String())
}

func (r *Report) Pages() []*Page {
	return r.pages
}

type Page struct {
	*Link
	Links []*Link
}

type Link struct {
	IsPage     bool
	StatusCode int
	Location   string
	Redirect   string
}
