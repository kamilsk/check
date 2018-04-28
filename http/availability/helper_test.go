package availability_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
)

var echoCode = func(rw http.ResponseWriter, req *http.Request) error {
	code, _ := strconv.Atoi(strings.TrimLeft(req.URL.Path, "/"))

	switch code {

	// redirects
	case http.StatusMovedPermanently:
		fallthrough
	case http.StatusFound:
		fallthrough
	case http.StatusSeeOther:
		fallthrough
	case http.StatusTemporaryRedirect:
		fallthrough
	case http.StatusPermanentRedirect:
		rw.Header().Set("Location", "https://test.dev/")
		http.Error(rw, http.StatusText(code), code)
		return nil

	default:
		return fmt.Errorf("can't handle request %q", req.URL.Path)
	}
}

var echoLinks = func(site1, site2 string) func(http.ResponseWriter, *http.Request) error {
	return func(rw http.ResponseWriter, req *http.Request) error {
		if req.URL.Path == "/" {
			links := []struct {
				Href string
				Text string
			}{
				{Href: site1 + "/", Text: "site1"},
				{Href: site1 + "/301", Text: "site1 - 301"},
				{Href: site1 + "/302", Text: "site1 - 302"},
				{Href: site1 + "/303", Text: "site1 - 303"},
				{Href: site1 + "/307", Text: "site1 - 307"},
				{Href: site1 + "/308", Text: "site1 - 308"},

				{Href: site2 + "/", Text: "site2"},
				{Href: site2 + "/301", Text: "site2 - 301"},
				{Href: site2 + "/302", Text: "site2 - 302"},
				{Href: site2 + "/303", Text: "site2 - 303"},
				{Href: site2 + "/307", Text: "site2 - 307"},
				{Href: site2 + "/308", Text: "site2 - 308"},

				//issue#34, fixes issue#30
				{Href: "#anchor", Text: "anchor"},
				{Href: "/#/tab", Text: "some tab"},
				//issue#35, fixes issue#30
				{Href: "mailto:test@my.email", Text: "test@my.email"},
				{Href: "tel:+01234567", Text: "+01234567"},
				// other bad urls
				{Href: ":bad", Text: "something bad"},
			}
			tpl.Execute(rw, links)
			return nil
		}
		return fmt.Errorf("can't handle request %q", req.URL.Path)
	}
}

func site() (server *httptest.Server, closer func()) {
	var (
		chain1, chain2, chain3 = &chain{}, &chain{}, &chain{}
		site1                  = httptest.NewServer(chain1)
		site2                  = httptest.NewServer(chain2)
		main                   = httptest.NewServer(chain3)
	)
	chain1.handlers = append(chain1.handlers, echoCode, echoLinks(site2.URL, main.URL))
	chain2.handlers = append(chain1.handlers, echoCode, echoLinks(site1.URL, main.URL))
	chain3.handlers = append(chain1.handlers, echoCode, echoLinks(site1.URL, site2.URL))
	return main, func() {
		site1.Close()
		site2.Close()
		main.Close()
	}
}

type chain struct {
	http.Handler
	handlers []func(rw http.ResponseWriter, req *http.Request) error
}

func (c *chain) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	for _, handler := range c.handlers {
		if err := handler(rw, req); err == nil {
			return
		}
	}
	if c.Handler == nil {
		c.Handler = http.NotFoundHandler()
	}
	c.Handler.ServeHTTP(rw, req)
}
