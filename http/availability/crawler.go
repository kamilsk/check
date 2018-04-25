package availability

import (
	"net/http"

	"github.com/gocolly/colly"
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
