package availability_test

import (
	"fmt"
	"html/template"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
)

var (
	html = `
<!doctype html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>Links</title>
</head>
<body>
<ul>{{ range . }}
    <li><a href="{{ .Href }}">{{ .Text }}</a></li>
{{ end }}</ul>
</body>
`
	tpl = template.Must(template.New("links").Parse(html))
)

var echoCode = func(rw http.ResponseWriter, req *http.Request) error {
	code, _ := strconv.Atoi(strings.TrimLeft(req.URL.Path, "/"))

	switch code {

	// ok
	case http.StatusOK:
		fallthrough
	case http.StatusCreated:
		fallthrough
	case http.StatusAccepted:
		rw.WriteHeader(code)
		return nil

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

	// client errors
	case http.StatusBadRequest:
		fallthrough
	case http.StatusUnauthorized:
		fallthrough
	case http.StatusForbidden:
		fallthrough
	case http.StatusNotFound:
		fallthrough
	// server errors
	case http.StatusInternalServerError:
		fallthrough
	case http.StatusNotImplemented:
		fallthrough
	case http.StatusBadGateway:
		fallthrough
	case http.StatusServiceUnavailable:
		fallthrough
	// and others
	case http.StatusMultipleChoices:
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
				Href template.URL
				Text string
			}{
				{Href: template.URL(site1 + "/"), Text: "site1"},
				{Href: template.URL(site1 + "/200"), Text: "site1 - 200"},
				{Href: template.URL(site1 + "/201"), Text: "site1 - 201"},
				{Href: template.URL(site1 + "/202"), Text: "site1 - 202"},
				{Href: template.URL(site1 + "/301"), Text: "site1 - 301"},
				{Href: template.URL(site1 + "/302"), Text: "site1 - 302"},
				{Href: template.URL(site1 + "/303"), Text: "site1 - 303"},
				{Href: template.URL(site1 + "/307"), Text: "site1 - 307"},
				{Href: template.URL(site1 + "/308"), Text: "site1 - 308"},
				{Href: template.URL(site1 + "/400"), Text: "site1 - 400"},
				{Href: template.URL(site1 + "/401"), Text: "site1 - 401"},
				{Href: template.URL(site1 + "/403"), Text: "site1 - 403"},
				{Href: template.URL(site1 + "/404"), Text: "site1 - 404"},
				{Href: template.URL(site1 + "/500"), Text: "site1 - 500"},
				{Href: template.URL(site1 + "/501"), Text: "site1 - 501"},
				{Href: template.URL(site1 + "/502"), Text: "site1 - 502"},
				{Href: template.URL(site1 + "/503"), Text: "site1 - 503"},
				{Href: template.URL(site1 + "/300"), Text: "site1 - 300"},

				{Href: template.URL(site2 + "/"), Text: "site2"},
				{Href: template.URL(site2 + "/200"), Text: "site2 - 200"},
				{Href: template.URL(site2 + "/201"), Text: "site2 - 201"},
				{Href: template.URL(site2 + "/202"), Text: "site2 - 202"},
				{Href: template.URL(site2 + "/301"), Text: "site2 - 301"},
				{Href: template.URL(site2 + "/302"), Text: "site2 - 302"},
				{Href: template.URL(site2 + "/303"), Text: "site2 - 303"},
				{Href: template.URL(site2 + "/307"), Text: "site2 - 307"},
				{Href: template.URL(site2 + "/308"), Text: "site2 - 308"},
				{Href: template.URL(site2 + "/400"), Text: "site2 - 400"},
				{Href: template.URL(site2 + "/401"), Text: "site2 - 401"},
				{Href: template.URL(site2 + "/403"), Text: "site2 - 403"},
				{Href: template.URL(site2 + "/404"), Text: "site2 - 404"},
				{Href: template.URL(site2 + "/500"), Text: "site2 - 500"},
				{Href: template.URL(site2 + "/501"), Text: "site2 - 501"},
				{Href: template.URL(site2 + "/502"), Text: "site2 - 502"},
				{Href: template.URL(site2 + "/503"), Text: "site2 - 503"},
				{Href: template.URL(site2 + "/300"), Text: "site2 - 300"},

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
