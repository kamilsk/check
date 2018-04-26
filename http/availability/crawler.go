package availability

import (
	"net/http"
	"net/url"

	"github.com/gocolly/colly"
	"github.com/gocolly/colly/debug"
	"github.com/pkg/errors"
)

const (
	locationHeader = "Location"
)

type Crawler interface {
	Visit(url string, bus EventBus) error
}

type CrawlerFunc func(string, EventBus) error

func (fn CrawlerFunc) Visit(url string, bus EventBus) error { return fn(url, bus) }

func CrawlerColly(userAgent string) Crawler {
	return CrawlerFunc(func(entry string, bus EventBus) error {
		defer close(bus)
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

func NoRedirect() func(*colly.Collector) {
	return func(c *colly.Collector) {
		c.RedirectHandler = func(*http.Request, []*http.Request) error {
			return http.ErrUseLastResponse
		}
	}
}

func OnError(bus EventBus) func(*colly.Collector) {
	return func(c *colly.Collector) {
		c.OnError(func(resp *colly.Response, err error) {
			var redirect string
			if err == http.ErrUseLastResponse {
				redirect = resp.Headers.Get(locationHeader)
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

// TODO use

type Debugger interface {
	debug.Debugger
}

type Option func(*Site)

func WithDebugger() Option {
	return func(*Site) {
		//
	}
}
