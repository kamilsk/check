package availability

import (
	"net/url"
	"sync"

	"github.com/kamilsk/check/errors"
)

// NewReport returns configured report builder.
func NewReport(options ...func(*Report)) *Report {
	r := &Report{}
	for _, f := range options {
		f(r)
	}
	return r
}

// CrawlerForSites sets the website crawler to a report builder.
func CrawlerForSites(crawler Crawler) func(*Report) {
	return func(r *Report) {
		r.crawler = crawler
	}
}

// Report represents a report builder.
type Report struct {
	crawler Crawler
	sites   []*Site
	ready   chan Site
}

// For prepares report builder for passed websites' URLs.
func (r *Report) For(rawURLs []string) *Report {
	r.sites = make([]*Site, 0, len(rawURLs))
	for _, rawURL := range rawURLs {
		r.sites = append(r.sites, NewSite(rawURL))
	}
	return r
}

// Fill starts to fetch sites and prepared them for reading.
func (r *Report) Fill() *Report {
	r.ready = make(chan Site, len(r.sites))
	for _, site := range r.sites {
		site.Error = site.Fetch(r.crawler)
		{
			copied := *site
			pages := make([]*Page, 0, len(site.Pages))
			for _, page := range site.Pages {
				page := *page
				pages = append(pages, &page)
				links := make([]Link, len(page.Links))
				copy(links, page.Links)
				page.Links = links
			}
			copied.Pages = pages
			r.ready <- copied
		}
	}
	close(r.ready)
	return r
}

// Sites returns a channel what will be closed when the report is available.
func (r *Report) Sites() <-chan Site {
	return r.ready
}

// NewSite returns new instance of the website.
func NewSite(rawURL string) *Site {
	u, err := url.Parse(rawURL)
	return &Site{
		url:   u,
		Name:  hostOrRawURL(u, rawURL),
		Error: errors.Wrapf(err, "parse rawURL %q for report", rawURL),
	}
}

// Site contains a meta information about a website.
type Site struct {
	url *url.URL

	Name     string
	Error    error
	Pages    []*Page
	Problems []ProblemEvent
}

// Fetch runs the website crawler and starts listen its events to build a website tree.
func (s *Site) Fetch(crawler Crawler) error {
	if s.Error != nil {
		return s.Error
	}
	if crawler == nil {
		s.Error = errors.Simple("crawler is not provided")
		return s.Error
	}
	var unexpected error
	wg, events := &sync.WaitGroup{}, make(chan event, 512)
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer errors.Recover(&unexpected)
		s.listen(events)
	}()
	s.Error = crawler.Visit(s.url.String(), events)
	wg.Wait()
	if unexpected != nil {
		panic(unexpected)
	}
	return s.Error
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
				pages[e.Page] = &Page{Links: make([]Link, 0, 8)}
			}
			linkToPage = append(linkToPage, [2]string{e.Href, e.Page})
		case ProblemEvent:
			s.Problems = append(s.Problems, e)
		default:
			panic(errors.Errorf("panic: unexpected event type %T", e))
		}
	}
	barrier := make(map[*Page]map[*Link]struct{})
	s.Pages = make([]*Page, 0, len(pages))
	for location, page := range pages {
		page.Link = links[location]
		s.Pages = append(s.Pages, page)
		barrier[page] = make(map[*Link]struct{})
	}
	for _, linkAndPage := range linkToPage {
		linkLocation, pageLocation := linkAndPage[0], linkAndPage[1]
		link := links[linkLocation]
		page := pages[pageLocation]
		if _, exists := barrier[page][link]; !exists {
			barrier[page][link] = struct{}{}
			{
				link := *link
				link.Page = page
				link.Internal = hasSameHost(page.Link.Location, link.Location)
				page.Links = append(page.Links, link)
			}
		}
	}
}

// Page contains meta information about a website page.
type Page struct {
	*Link
	Links []Link
}

// Link contains meta information about a web link.
type Link struct {
	Page       *Page
	Internal   bool
	StatusCode int
	Location   string
	Redirect   string
	Error      error
}

func hostOrRawURL(u *url.URL, raw string) string {
	if u == nil {
		return raw
	}
	return u.Host
}

func hasSameHost(link1, link2 string) bool {
	u1, _ := url.Parse(link1)
	u2, _ := url.Parse(link2)
	return u1 != nil && u2 != nil && u1.Host == u2.Host
}

type event interface {
	family()
}

// NewReadableEventBus returns read/write-channel of events.
func NewReadableEventBus(size int) chan event {
	return make(chan event, size)
}

// EventBus is a write-only channel to communicate between a website crawler and a report builder.
type EventBus chan<- event

// ErrorEvent contains a response' status code, its URL and an encountered error.
type ErrorEvent struct {
	event

	StatusCode int
	Location   string
	Redirect   string
	Error      error
}

// ResponseEvent contains a response' status code and its URL.
type ResponseEvent struct {
	event

	StatusCode int
	Location   string
}

// WalkEvent contains information about a page and a link located on it.
type WalkEvent struct {
	event

	Page string
	Href string
}

// ProblemEvent contains information about unexpected error.
type ProblemEvent struct {
	event

	Message string
	Context interface{}
}
