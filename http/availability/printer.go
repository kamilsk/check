package availability

import (
	"fmt"
	"io"
	"sort"

	"github.com/fatih/color"
)

const (
	success = "success"
	warning = "warning"
	danger  = "danger"
)

func NewPrinter(reports []*Report) *Printer {
	return &Printer{
		colors: map[string]*color.Color{
			success: color.New(color.FgWhite),
			warning: color.New(color.FgYellow),
			danger:  color.New(color.FgRed, color.Bold),
		},
		reports: reports,
	}
}

type Printer struct {
	colors  map[string]*color.Color
	reports []*Report
}

func (p *Printer) Print(w io.Writer) {
	for _, report := range p.reports {
		if err := report.Error(); err != nil {
			p.colorize(999).Fprintf(w, "report %q has error %q\n", report.Name(), err)
			continue
		}
		pages := pagesByLocation(report.Pages())
		sort.Sort(pages)
		for _, page := range pages {
			p.colorize(page.StatusCode).Fprintf(w, "[%d] %s\n", page.StatusCode, page.Location)
			last := len(page.Links) - 1
			sort.Sort(linksByStatusCode(page.Links))
			for i, link := range page.Links {
				if i == last {
					p.colorize(link.StatusCode).Fprintf(w, "    └─── [%d] %s\n", link.StatusCode, link.Location)
					continue
				}
				p.colorize(link.StatusCode).Fprintf(w, "    ├─── [%d] %s\n", link.StatusCode, link.Location)
			}
		}
	}
}

func (p *Printer) colorize(statusCode int) printer {
	var pr printer
	switch {
	case statusCode < 300:
		pr, _ = p.colors[success]
	case statusCode >= 300 && statusCode < 400:
		pr, _ = p.colors[warning]
	case statusCode >= 400:
		pr, _ = p.colors[danger]
	}
	if pr == nil {
		pr = defaultPrinter(fmt.Fprintf)
	}
	return pr
}

type printer interface {
	Fprintf(io.Writer, string, ...interface{}) (int, error)
}

type defaultPrinter func(io.Writer, string, ...interface{}) (int, error)

func (fn defaultPrinter) Fprintf(w io.Writer, format string, a ...interface{}) (n int, err error) {
	return fn(w, format, a...)
}

type pagesByLocation []*Page

func (l pagesByLocation) Len() int { return len(l) }

func (l pagesByLocation) Less(i, j int) bool { return l[i].Location < l[j].Location }

func (l pagesByLocation) Swap(i, j int) { l[i], l[j] = l[j], l[i] }

type linksByStatusCode []*Link

func (l linksByStatusCode) Len() int { return len(l) }

func (l linksByStatusCode) Less(i, j int) bool { return l[i].StatusCode < l[j].StatusCode }

func (l linksByStatusCode) Swap(i, j int) { l[i], l[j] = l[j], l[i] }
