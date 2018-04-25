package availability

import (
	"net/url"
	"sync"

	"github.com/gocolly/colly"
	"github.com/pkg/errors"
)

func CreateReport(rawURL string) *Report {
	report := NewReport(rawURL)
	report.error = report.Get()
	return report
}

func NewReport(rawURL string) *Report {
	location, err := url.Parse(rawURL)
	return &Report{
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
	error    error
	location *url.URL
	mu       *sync.RWMutex
	pages    []*Page
	journal  map[string]*Link
}

func (r *Report) Get() error {
	if r.error != nil {
		return r.error
	}
	c := colly.NewCollector(colly.UserAgent("check"), colly.IgnoreRobotsTxt())
	c.OnRequest(func(req *colly.Request) {
		link := r.createLink(req.URL)
		if link.IsPage {
			r.createPage(link)
		}
	})
	c.OnError(func(resp *colly.Response, err error) {
		r.setStatus(resp.Request.URL, resp.StatusCode)
	})
	c.OnResponse(func(resp *colly.Response) {
		r.setStatus(resp.Request.URL, resp.StatusCode)
	})
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		if r.isPage(e.Request.URL) {
			attr := e.Attr("href")
			href := e.Request.AbsoluteURL(attr)
			if href == "" {

				panic("invalid URL " + attr) // TODO set error instead of panic

			}

			// TODO make thread safe
			link := r.createLinkByHref(href)
			page := r.findPage(e.Request.URL)
			page.Links = append(page.Links, link)

			// TODO it can return
			// &errors.errorString{s:""}
			// &errors.errorString{s:"URL already visited"}
			e.Request.Visit(href)
		}
	})
	return c.Visit(r.location.String())
}

func (r *Report) Pages() []*Page {
	r.mu.RLock()

	// TODO return []Page instead []*Page
	pages := make([]*Page, len(r.pages))
	copy(pages, r.pages)

	r.mu.RUnlock()
	return pages
}

func (r *Report) createLink(location *url.URL) *Link {
	href := location.String()

	{
		r.mu.RLock()
		link, ok := r.journal[href]
		if ok {
			r.mu.RUnlock()
			return link
		}
		r.mu.RUnlock()
	}

	{
		r.mu.Lock()
		link, ok := r.journal[href]
		if ok {
			r.mu.Unlock()
			return link
		}
		link = &Link{IsPage: r.isPage(location), Location: href}
		r.journal[href] = link
		r.mu.Unlock()
		return link
	}
}

func (r *Report) createLinkByHref(href string) *Link {
	location, err := url.Parse(href)
	if err != nil {

		panic(err) // TODO set error instead of panic

	}
	return r.createLink(location)
}

func (r *Report) createPage(link *Link) *Page {
	r.mu.Lock()
	page := &Page{Link: link, Links: make([]*Link, 0, 8)}
	r.pages = append(r.pages, page)
	r.mu.Unlock()
	return page
}

func (r *Report) findPage(location *url.URL) *Page {
	href := location.String()
	r.mu.RLock()
	defer r.mu.RUnlock()
	link, ok := r.journal[href]
	if !ok {

		panic("can't find link with URL " + href) // TODO set error instead of panic

	}
	for _, page := range r.pages {
		if page.Link == link {
			return page
		}
	}

	panic("can't find page with URL " + href) // TODO set error instead of panic
}

func (r *Report) isPage(location *url.URL) bool {
	return location.Hostname() == r.location.Hostname()
}

func (r *Report) setStatus(location *url.URL, statusCode int) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	href := location.String()
	link, ok := r.journal[href]
	if !ok {

		panic("unexpected URL " + href) // TODO set error instead of panic

	}
	link.StatusCode = statusCode
}

type Page struct {
	*Link
	Links []*Link
}

type Link struct {
	StatusCode int
	IsPage     bool
	Location   string
}
