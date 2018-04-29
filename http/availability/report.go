package availability

import (
	"fmt"
	"net/url"
	"sync"

	"github.com/kamilsk/check/errors"
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
	ready   chan Site
}

func (r *Report) For(rawURLs []string) *Report {
	r.sites = make([]*Site, 0, len(rawURLs))
	r.ready = make(chan Site, len(rawURLs))
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
			var copied Site
			copied.name = site.name
			defer wg.Done()
			defer func() { r.ready <- copied }()
			defer errors.Recover(&copied.error)
			site.error = site.Fetch(r.crawler)
			{
				copied = *site
				pages := make([]*Page, 0, len(site.Pages))
				for _, page := range site.Pages {
					page := *page
					pages = append(pages, &page)
					links := make([]*Link, 0, len(page.Links))
					for _, link := range page.Links {
						link := *link
						links = append(links, &link)
					}
					page.Links = links
				}
				copied.Pages = pages
			}
		}(site)
	}
	wg.Wait()
	close(r.ready)
	return r
}

func (r *Report) Sites() <-chan Site {
	return r.ready
}

func NewSite(rawURL string) *Site {
	u, err := url.Parse(rawURL)
	return &Site{
		name:  hostOrRawURL(u, rawURL),
		url:   u,
		error: errors.WithMessage(err, fmt.Sprintf("parse rawURL %q for report", rawURL)),
	}
}

type Site struct {
	name  string
	url   *url.URL
	error error

	Pages    []*Page
	Problems []ProblemEvent
}

func (s *Site) Name() string { return s.name }

func (s *Site) Error() error { return s.error }

func (s *Site) Fetch(crawler Crawler) error {
	if s.error != nil {
		return s.error
	}
	if crawler == nil {
		s.error = errors.Simple("crawler is not provided")
		return s.error
	}
	var unexpected error
	wg, events := &sync.WaitGroup{}, make(chan event, 512)
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer errors.Recover(&unexpected)
		s.listen(events)
	}()
	s.error = crawler.Visit(s.url.String(), events)
	wg.Wait()
	if unexpected != nil {
		panic(unexpected)
	}
	return s.error
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
		case ProblemEvent:
			s.Problems = append(s.Problems, e)
		default:
			panic(errors.Errorf("panic: unexpected event type %T", e))
		}
	}
	barrier := make(map[*Page]map[*Link]struct{})
	s.Pages = make([]*Page, 0, len(pages))
	for location, page := range pages {
		link, found := links[location]
		if !found {
			panic(errors.Errorf("panic: not consistent fetch result. link %q not found", location))
		}
		page.Link, link.Page = link, page
		s.Pages = append(s.Pages, page)
		barrier[page] = make(map[*Link]struct{})
	}
	for _, linkAndPage := range linkToPage {
		linkLocation, pageLocation := linkAndPage[0], linkAndPage[1]
		link, found := links[linkLocation]
		if !found {
			panic(errors.Errorf("panic: not consistent fetch result. link %q not found", linkLocation))
		}
		if _, found := pages[linkLocation]; found {
			continue // exclude internal links
		}
		page, found := pages[pageLocation]
		if !found {
			panic(errors.Errorf("panic: not consistent fetch result. page %q not found", pageLocation))
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

func hostOrRawURL(u *url.URL, raw string) string {
	if u == nil {
		return raw
	}
	return u.Host
}

type event interface {
	family()
}

// NewReadableEventBus returns rw-channel of events.
func NewReadableEventBus(size int) chan event {
	return make(chan event, size)
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

type ProblemEvent struct {
	event

	Message string
	Context interface{}
}
