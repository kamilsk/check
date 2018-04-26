package availability

import (
	"net/url"
	"sync"

	"github.com/gocolly/colly"
	"github.com/pkg/errors"
)

func NewReport(options ...func(*Report)) *Report {
	r := &Report{}
	for _, f := range options {
		f(r)
	}
	return r
}

type Report struct {
	sites []*Site
}

func (r *Report) For(rawURLs []string) *Report {
	r.sites = make([]*Site, 0, len(rawURLs))
	for _, rawURL := range rawURLs {
		r.sites = append(r.sites, NewSite(rawURL))
	}
	return r
}

func (r *Report) Fill() *Report {
	wg, start := &sync.WaitGroup{}, make(chan struct{})
	for _, site := range r.sites {
		wg.Add(1)
		go func(site *Site) {
			defer wg.Done()
			<-start
			site.Fetch()
		}(site)
	}
	close(start)
	wg.Wait()
	return r
}

func (r *Report) Sites() []Site {
	sites := make([]Site, 0, len(r.sites))
	for _, site := range r.sites {
		sites = append(sites, *site)
	}
	return sites
}

func NewSite(rawURL string) *Site {
	u, err := url.Parse(rawURL)
	return &Site{
		name:  hostOrRawURL(u, rawURL),
		url:   u,
		error: errors.Wrapf(err, "parse rawURL %q for report", rawURL),

		// deprecated
		mu: &sync.RWMutex{}, Pages: make([]*Page, 0, 8), journal: make(map[string]*Link),
	}
}

func hostOrRawURL(u *url.URL, raw string) string {
	if u == nil {
		return raw
	}
	return u.Host
}

// ~

type Site struct {
	name  string
	url   *url.URL
	error error

	Pages []*Page

	// deprecated
	mu      *sync.RWMutex
	journal map[string]*Link
}

func (r *Site) Name() string { return r.name }

func (r *Site) Error() error { return r.error }

func (r *Site) Fetch() error {
	if r.error != nil {
		return r.error
	}
	c := colly.NewCollector(
		UserAgent(), NoRedirect(), colly.IgnoreRobotsTxt(),

		TempOption(r),
	)
	return c.Visit(r.url.String())
}

type Page struct {
	*Link
	Links []*Link
}

type Link struct {
	// deprecated
	IsPage bool

	Page       *Page
	StatusCode int
	Location   string
	Redirect   string
	Error      error
}
