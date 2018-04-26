package availability

import (
	"net/http"
	"net/url"

	"github.com/gocolly/colly"
	"github.com/gocolly/colly/debug"
)

func UserAgent() func(*colly.Collector) {
	return colly.UserAgent("check")
}

func NoRedirect() func(*colly.Collector) {
	return func(c *colly.Collector) {
		c.RedirectHandler = func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		}
	}
}

// ~

func TempOption(r *Report) func(*colly.Collector) {
	return func(c *colly.Collector) {
		c.OnRequest(func(req *colly.Request) {
			link := r.createLink(req.URL)
			if link.IsPage {
				r.createPage(link)
			}
		})

		c.OnError(func(resp *colly.Response, err error) { r.setStatus(resp) })
		c.OnResponse(func(resp *colly.Response) { r.setStatus(resp) })

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

				// TODO it can return error
				// &errors.errorString{s:""}
				// &errors.errorString{s:"URL already visited"}
				e.Request.Visit(href)
			}
		})
	}
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

func (r *Report) setStatus(resp *colly.Response) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	href := resp.Request.URL.String()
	link, ok := r.journal[href]
	if !ok {

		panic("unexpected URL " + href) // TODO set error instead of panic

	}
	link.StatusCode = resp.StatusCode
	if _, is := redirects[link.StatusCode]; is {
		link.Redirect = resp.Headers.Get(location)
	}
}

// ~

type Debugger interface {
	debug.Debugger
}

type Option func(*Report)

func WithDebugger() Option {
	return func(*Report) {
		//
	}
}
