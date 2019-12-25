package availability_test

import (
	"io/ioutil"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/kamilsk/check/http/availability"
)

func TestCrawlerColly(t *testing.T) {
	site, closer := site()
	defer closer()
	{
		var errorEvents, redirectEvents, responseEvents, walkEvents, problemEvents, unknownEvents int
		wg, bus := &sync.WaitGroup{}, availability.NewReadableEventBus(8)
		wg.Add(1)
		go func() {
			defer wg.Done()
			for event := range bus {
				switch e := event.(type) {
				case availability.ErrorEvent:
					errorEvents++
					if e.Redirect != "" {
						redirectEvents++
					}
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
		assert.NoError(t, crawler.Visit(site.URL+"/", bus))
		wg.Wait()
		assert.Equal(t, 28, errorEvents)
		assert.Equal(t, 10, redirectEvents)
		assert.Equal(t, 8, responseEvents)
		assert.Equal(t, 37, walkEvents)
		assert.Equal(t, 1, problemEvents)
		assert.Empty(t, unknownEvents)
	}
	{
		crawler := availability.CrawlerColly(availability.CrawlerConfig{})
		assert.Error(t, crawler.Visit(":bad", make(availability.EventBus)))
	}
}
