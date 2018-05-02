package availability

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/gocolly/colly"
	"github.com/gocolly/colly/debug"
	"github.com/kamilsk/check/errors"
)

const (
	locationHeader = "Location"
	clickOptHeader = "X-Click-Options"
)

var clickOptions = []string{"anonym", "nolog"}

// Crawler defines general behavior of website crawlers.
type Crawler interface {
	// Visit starts to crawl a website starting with the passed URL.
	// All implementations must communicate with report builder
	// through the passed EventBus and close it when they're done.
	Visit(url string, bus EventBus) error
}

// CrawlerConfig holds a website crawler's configuration.
type CrawlerConfig struct {
	UserAgent string
	Verbose   bool
	Output    io.Writer
}

// CrawlerFunc adds possibility to use functions as a website crawler.
type CrawlerFunc func(string, EventBus) error

// Visit calls itself.
func (fn CrawlerFunc) Visit(url string, bus EventBus) error { return fn(url, bus) }

// CrawlerColly returns configured website crawler.
func CrawlerColly(config CrawlerConfig) Crawler {
	return CrawlerFunc(func(entry string, bus EventBus) error {
		defer close(bus)
		base, err := url.Parse(entry)
		if err != nil {
			return errors.WithMessage(err, fmt.Sprintf("parse entry point URL %q", entry))
		}
		options := make([]func(*colly.Collector), 0, 9)
		if config.UserAgent != "" {
			options = append(options, colly.UserAgent(config.UserAgent))
		}
		if config.Verbose {
			options = append(options, colly.Debugger(&debug.LogDebugger{Output: config.Output}))
		}
		options = append(options,
			colly.IgnoreRobotsTxt(),
			NoCookie(),
			NoRedirect(),
			OnRequest(),
			OnError(bus),
			OnResponse(bus),
			OnHTML(base, bus),
		)
		return colly.NewCollector(options...).Visit(entry)
	})
}

// NoRedirect disables redirects for `github.com/gocolly/colly.Collector`.
func NoRedirect() func(*colly.Collector) {
	return func(c *colly.Collector) {
		c.RedirectHandler = func(*http.Request, []*http.Request) error {
			return http.ErrUseLastResponse
		}
	}
}

// NoCookie disables cookie for `github.com/gocolly/colly.Collector`.
func NoCookie() func(*colly.Collector) {
	return func(c *colly.Collector) {
		c.DisableCookies()
	}
}

// OnRequest registers a callback by `github.com/gocolly/colly.Collector.OnRequest()`.
func OnRequest() func(*colly.Collector) {
	options := strings.Join(clickOptions, ";")
	return func(c *colly.Collector) {
		c.OnRequest(func(req *colly.Request) {
			req.Headers.Set(clickOptHeader, options)
		})
	}
}

// OnError registers a callback by `github.com/gocolly/colly.Collector.OnError()`.
func OnError(bus EventBus) func(*colly.Collector) {
	return func(c *colly.Collector) {
		c.OnError(func(resp *colly.Response, err error) {
			location, redirect := resp.Request.URL.String(), ""
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

// OnResponse registers a callback by `github.com/gocolly/colly.Collector.OnResponse()`.
func OnResponse(bus EventBus) func(*colly.Collector) {
	return func(c *colly.Collector) {
		c.OnResponse(func(resp *colly.Response) {
			location := resp.Request.URL.String()
			bus <- ResponseEvent{
				StatusCode: resp.StatusCode,
				Location:   location,
			}
		})
	}
}

// OnHTML registers a callback by `github.com/gocolly/colly.Collector.OnHTML()`.
func OnHTML(base *url.URL, bus EventBus) func(*colly.Collector) {
	isPage := func(current *url.URL) bool {
		return current.Host == base.Host
	}
	return func(c *colly.Collector) {
		c.OnHTML("a[href]", func(el *colly.HTMLElement) {
			if isPage(el.Request.URL) {
				attr := el.Attr("href")
				if strings.HasPrefix(attr, "#") {
					return
				}
				href := el.Request.AbsoluteURL(attr)
				if href == "" {
					bus <- ProblemEvent{Message: "bad url", Context: struct {
						Page string
						Href string
					}{el.Request.URL.String(), attr}}
					return
				}
				if !strings.HasPrefix(href, "http") {
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
