package availability

import (
	"net/url"
	"sync"

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
	binder Binder
	sites  []*Site
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
			site.Fetch(r.binder)
		}(site)
	}
	close(start)
	wg.Wait()
	return r
}

func (r *Report) Sites() []Site {
	sites := make([]Site, 0, len(r.sites))
	for _, site := range r.sites {
		site := *site
		{
			copied := make([]*Page, len(site.Pages))
			copy(copied, site.Pages)
			for _, page := range copied {
				copied := make([]*Link, len(page.Links))
				copy(copied, page.Links)
				page.Links = copied
			}
			site.Pages = copied
		}
		sites = append(sites, site)
	}
	return sites
}

func NewSite(rawURL string) *Site {
	u, err := url.Parse(rawURL)
	return &Site{
		name:  hostOrRawURL(u, rawURL),
		url:   u,
		error: errors.Wrapf(err, "parse rawURL %q for report", rawURL),

		Pages: make([]*Page, 0, 8),

		mu: &sync.RWMutex{}, journal: make(map[string]*Link),
	}
}

func hostOrRawURL(u *url.URL, raw string) string {
	if u == nil {
		return raw
	}
	return u.Host
}

// ~

func ClientProvider(b Binder) func(*Report) {
	return func(r *Report) {
		r.binder = b
	}
}

type Binder interface {
	Bind(*Site) Client
}

type Client interface {
	Visit(URL string) error
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

func (r *Site) Fetch(binder Binder) error {
	if r.error != nil {
		return r.error
	}
	return binder.Bind(r).Visit(r.url.String())
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

func (l *Link) FullLocation(sep string) string {
	if l.Redirect != "" {
		return l.Location + sep + l.Redirect
	}
	return l.Location
}
