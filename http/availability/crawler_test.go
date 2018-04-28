package availability_test

import (
	"html/template"
	"io/ioutil"
	"sync"
	"testing"

	"github.com/kamilsk/check/http/availability"
	"github.com/stretchr/testify/assert"
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

func TestCrawlerColly(t *testing.T) {
	site, closer := site()
	defer closer()

	{
		var errorEvents, responseEvents, walkEvents, problemEvents, unknownEvents int
		wg, bus := &sync.WaitGroup{}, availability.NewReadableEventBus(8)
		wg.Add(1)
		go func() {
			defer wg.Done()
			for event := range bus {
				switch event.(type) {
				case availability.ErrorEvent:
					errorEvents++
				case availability.ResponseEvent:
					responseEvents++
				case availability.WalkEvent:
					walkEvents++
				case availability.ProblemEvent:
					problemEvents++
				default:
					unknownEvents++
				}
			}
		}()
		crawler := availability.CrawlerColly(availability.CrawlerConfig{
			UserAgent: "test/dev",
			Verbose:   true,
			Output:    ioutil.Discard,
		})
		assert.NoError(t, crawler.Visit(site.URL, bus))
		wg.Wait()
		assert.Equal(t, 11, errorEvents)
		assert.Equal(t, 3, responseEvents)
		assert.Equal(t, 42, walkEvents)
		assert.Equal(t, 9, problemEvents)
		assert.Empty(t, unknownEvents)
	}

	{
		crawler := availability.CrawlerColly(availability.CrawlerConfig{})
		assert.Error(t, crawler.Visit(":bad", make(availability.EventBus)))
	}
}
