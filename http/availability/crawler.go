package availability

import (
	"net/http"
	"net/url"

	"github.com/gocolly/colly"
	"github.com/gocolly/colly/debug"
	"github.com/pkg/errors"
)

const (
	location = "Location"
)

type Crawler interface {
	Visit(url string, bus EventBus) error
}

type CrawlerFunc func(string, EventBus) error

func (fn CrawlerFunc) Visit(url string, bus EventBus) error { return fn(url, bus) }

func CrawlerColly(userAgent string) Crawler {
	return CrawlerFunc(func(entry string, bus EventBus) error {
		base, err := url.Parse(entry)
		if err != nil {
			return errors.Wrapf(err, "parse entry point URL %q", entry)
		}
		return colly.NewCollector(
			colly.UserAgent(userAgent), colly.IgnoreRobotsTxt(), NoRedirect(),
			OnError(bus), OnResponse(bus), OnHTML(base, bus),
		).Visit(entry)
	})
}

func OnError(bus EventBus) func(*colly.Collector) {
	return func(c *colly.Collector) {
		c.OnError(func(resp *colly.Response, err error) {
			var redirect string
			if err == http.ErrUseLastResponse {
				redirect = resp.Headers.Get(location)
			}
			bus <- ErrorEvent{
				StatusCode: resp.StatusCode,
				Location:   resp.Request.URL.String(),
				Redirect:   redirect,
				Error:      err,
			}
		})
	}
}

func OnResponse(bus EventBus) func(*colly.Collector) {
	return func(c *colly.Collector) {
		c.OnResponse(func(resp *colly.Response) {
			bus <- ResponseEvent{
				StatusCode: resp.StatusCode,
				Location:   resp.Request.URL.String(),
			}
		})
	}
}

func OnHTML(base *url.URL, bus EventBus) func(*colly.Collector) {
	isPage := func(current *url.URL) bool {
		return current.Hostname() == base.Hostname()
	}
	return func(c *colly.Collector) {
		c.OnHTML("a[href]", func(el *colly.HTMLElement) {
			if isPage(el.Request.URL) {
				href := el.Request.AbsoluteURL(el.Attr("href"))
				bus <- WalkEvent{
					Page: el.Request.URL.String(),
					Href: href,
				}
				el.Request.Visit(href)
			}
		})
	}
}

// ~

type Debugger interface {
	debug.Debugger
}

type Option func(*Site)

func WithDebugger() Option {
	return func(*Site) {
		//
	}
}

// ~

func NoRedirect() func(*colly.Collector) {
	return func(c *colly.Collector) {
		c.RedirectHandler = func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		}
	}
}

// ~

func TempOption(r *Site) func(*colly.Collector) {
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

func (r *Site) createLink(location *url.URL) *Link {
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

func (r *Site) createLinkByHref(href string) *Link {
	location, err := url.Parse(href)
	if err != nil {

		panic(err) // TODO set error instead of panic

	}
	return r.createLink(location)
}

func (r *Site) createPage(link *Link) *Page {
	r.mu.Lock()
	page := &Page{Link: link, Links: make([]*Link, 0, 8)}
	r.Pages = append(r.Pages, page)
	r.mu.Unlock()
	return page
}

func (r *Site) findPage(location *url.URL) *Page {
	href := location.String()
	r.mu.RLock()
	defer r.mu.RUnlock()
	link, ok := r.journal[href]
	if !ok {

		panic("can't find link with URL " + href) // TODO set error instead of panic

	}
	for _, page := range r.Pages {
		if page.Link == link {
			return page
		}
	}

	panic("can't find page with URL " + href) // TODO set error instead of panic
}

func (r *Site) isPage(location *url.URL) bool {
	return location.Hostname() == r.url.Hostname()
}

func (r *Site) setStatus(resp *colly.Response) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	href := resp.Request.URL.String()
	link, ok := r.journal[href]
	if !ok {

		panic("unexpected URL " + href) // TODO set error instead of panic

	}
	link.StatusCode = resp.StatusCode
}
