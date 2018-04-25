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

	head = "[%d] %s\n"
	body = "    ├─── [%d] %s\n"
	foot = "    └─── [%d] %s\n"
)

var colors = map[string]*color.Color{
	success: color.New(color.FgWhite),
	warning: color.New(color.FgYellow),
	danger:  color.New(color.FgRed, color.Bold),
}

func NewPrinter(reports []*Report) *Printer {
	return &Printer{reports: reports}
}

type Printer struct {
	reports []*Report
}

func (p *Printer) Print(w io.Writer) {
	for _, report := range p.reports {
		if err := report.Error(); err != nil {
			important().Fprintf(w, "report %q has error %q\n", report.Name(), err)
			continue
		}
		pages := pagesByLocation(report.Pages())
		sort.Sort(pages)
		for _, page := range pages {
			colorize(page.StatusCode).Fprintf(w, head, page.StatusCode, page.Location)
			last := len(page.Links) - 1
			sort.Sort(linksByStatusCode(page.Links))
			for i, link := range page.Links {
				if i == last {
					colorize(link.StatusCode).Fprintf(w, foot, link.StatusCode, join(link.Location, link.Redirect))
					continue
				}
				colorize(link.StatusCode).Fprintf(w, body, link.StatusCode, join(link.Location, link.Redirect))
			}
		}
	}
}

func colorize(statusCode int) printer {
	var p printer
	switch {
	case statusCode < 300:
		p, _ = colors[success]
	case statusCode >= 300 && statusCode < 400:
		p, _ = colors[warning]
	case statusCode >= 400:
		p, _ = colors[danger]
	}
	if p == nil {
		p = defaultPrinter(fmt.Fprintf)
	}
	return p
}

func important() printer { return colorize(999) }

func join(location, redirect string) string {
	if redirect == "" {
		return location
	}
	return location + " -> " + redirect
}

type printer interface {
	Fprintf(io.Writer, string, ...interface{}) (int, error)
}

type defaultPrinter func(io.Writer, string, ...interface{}) (int, error)

func (fn defaultPrinter) Fprintf(w io.Writer, format string, a ...interface{}) (int, error) {
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
