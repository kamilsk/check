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

func CrawlerForSites(crawler Crawler) func(*Report) {
	return func(r *Report) {
		r.crawler = crawler
	}
}

type Report struct {
	crawler Crawler
	sites   []*Site
}

func (r *Report) For(rawURLs []string) *Report {
	r.sites = make([]*Site, 0, len(rawURLs))
	for _, rawURL := range rawURLs {
		r.sites = append(r.sites, NewSite(rawURL))
	}
	return r
}

func (r *Report) Fill() *Report {
	wg := &sync.WaitGroup{}
	for _, site := range r.sites {
		wg.Add(1)
		go func(site *Site) {
			defer wg.Done()
			_ = site.Fetch(r.crawler)
		}(site)
	}
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

type event interface {
	family()
}

type EventBus chan<- event

type ErrorEvent struct {
	event

	StatusCode int
	Location   string
	Redirect   string
	Error      error
}

type ResponseEvent struct {
	event

	StatusCode int
	Location   string
}

type WalkEvent struct {
	event

	Page string
	Href string
}

// ~

type Site struct {
	name  string
	url   *url.URL
	error error

	// not thread-safe
	Pages []*Page

	// deprecated
	mu      *sync.RWMutex
	journal map[string]*Link
}

func (s *Site) Name() string { return s.name }

func (s *Site) Error() error { return s.error }

func (s *Site) Fetch(crawler Crawler) error {
	if s.error != nil {
		return s.error
	}
	wg, events := &sync.WaitGroup{}, make(chan event, 512)
	wg.Add(1)
	go func() {
		defer wg.Done()
		s.listen(events)
	}()
	if err := crawler.Visit(s.url.String(), events); err != nil {
		return err
	}
	wg.Wait()
	return nil
}

func (s *Site) listen(events <-chan event) {
	links := make(map[string]*Link)
	pages := make(map[string]*Page)
	linkToPage := make([][2]string, 0, 512)
	for event := range events {
		switch e := event.(type) {
		case ErrorEvent:
			if _, exists := links[e.Location]; !exists {
				links[e.Location] = &Link{
					StatusCode: e.StatusCode,
					Location:   e.Location,
					Redirect:   e.Redirect,
					Error:      e.Error,
				}
			}
		case ResponseEvent:
			if _, exists := links[e.Location]; !exists {
				links[e.Location] = &Link{
					StatusCode: e.StatusCode,
					Location:   e.Location,
				}
			}
		case WalkEvent:
			if _, exists := pages[e.Page]; !exists {
				pages[e.Page] = &Page{Links: make([]*Link, 0, 8)}
			}
			linkToPage = append(linkToPage, [2]string{e.Href, e.Page})
		default:
			panic(errors.Errorf("unexpected event type %T", e))
		}
	}
	barrier := make(map[*Page]map[*Link]struct{})
	s.Pages = make([]*Page, 0, len(pages))
	for location, page := range pages {
		link, found := links[location]
		if !found {
			panic(errors.Errorf("not consistent fetch result. link %q not found", location))
		}
		page.Link, link.Page = link, page
		s.Pages = append(s.Pages, page)
		barrier[page] = make(map[*Link]struct{})
	}
	for _, linkAndPage := range linkToPage {
		linkLocation, pageLocation := linkAndPage[0], linkAndPage[1]
		link, found := links[linkLocation]
		if !found {
			panic(errors.Errorf("not consistent fetch result. link %q not found", linkLocation))
		}
		if _, found := pages[linkLocation]; found {
			continue // exclude internal links
		}
		page, found := pages[pageLocation]
		if !found {
			panic(errors.Errorf("not consistent fetch result. page %q not found", pageLocation))
		}
		if _, exists := barrier[page][link]; !exists {
			page.Links = append(page.Links, link)
			barrier[page][link] = struct{}{}
		}
	}
}

type Page struct {
	*Link
	Links []*Link
}

type Link struct {
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
