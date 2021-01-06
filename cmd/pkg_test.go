package cmd

import (
	"html/template"
	"net/http"
	"net/http/httptest"

	"go.octolab.org/unsafe"
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
<a href="{{ .Href }}">{{ .Text }}</a>
</body>
`
	tpl = template.Must(template.New("links").Parse(html))
)

func site() (server *httptest.Server, closer func()) {
	main := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		unsafe.Ignore(tpl.Execute(rw, struct {
			Href string
			Text string
		}{"/", "Home"}))
	}))
	return main, main.Close
}
