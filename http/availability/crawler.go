package availability

import (
	"io"
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

type CrawlerConfig struct {
	UserAgent string
	Verbose   bool
	Output    io.Writer
}

type CrawlerFunc func(string, EventBus) error

func (fn CrawlerFunc) Visit(url string, bus EventBus) error { return fn(url, bus) }

func CrawlerColly(config CrawlerConfig) Crawler {
	return CrawlerFunc(func(entry string, bus EventBus) error {
		defer close(bus)
		base, err := url.Parse(entry)
		if err != nil {
			return errors.Wrapf(err, "parse entry point URL %q", entry)
		}
		options := make([]func(*colly.Collector), 0, 7)
		options = append(options, colly.UserAgent(config.UserAgent))
		options = append(options, colly.IgnoreRobotsTxt())
		if config.Verbose {
			options = append(options, colly.Debugger(&debug.LogDebugger{Output: config.Output}))
		}
		options = append(options, NoRedirect())
		options = append(options, OnError(bus))
		options = append(options, OnResponse(bus))
		options = append(options, OnHTML(base, bus))
		return colly.NewCollector(options...).Visit(entry)
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
			location := resp.Request.URL.String()

			//issue#30:on investigation
			if location == "" {
				bus <- ProblemEvent{Message: "empty location", Context: struct {
					Response *colly.Response
					Error    error
				}{resp, err}}
				return
			}

			var redirect string
			if resp.Headers != nil {
				redirect = resp.Headers.Get(locationHeader)
			}
			bus <- ErrorEvent{
				StatusCode: resp.StatusCode,
				Location:   location,
				Redirect:   redirect,
				Error:      err,
			}
		})
	}
}

func OnResponse(bus EventBus) func(*colly.Collector) {
	return func(c *colly.Collector) {
		c.OnResponse(func(resp *colly.Response) {
			location := resp.Request.URL.String()

			//issue#30:on investigation
			if location == "" {
				bus <- ProblemEvent{Message: "empty location", Context: resp}
				return
			}

			bus <- ResponseEvent{
				StatusCode: resp.StatusCode,
				Location:   location,
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

				//issue#30:on investigation
				if href == "" {
					bus <- ProblemEvent{Message: "empty location", Context: struct {
						Page string
						Href string
					}{el.Request.URL.String(), el.Attr("href")}}
					return
				}

				bus <- WalkEvent{
					Page: el.Request.URL.String(),
					Href: href,
				}
				el.Request.Visit(href)
			}
		})
	}
}
