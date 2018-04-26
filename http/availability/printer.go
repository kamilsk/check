package availability

import (
	"fmt"
	"io"
	"os"
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

func Output(output io.Writer) func(*Printer) {
	return func(p *Printer) {
		p.output = output
	}
}

func NewPrinter(options ...func(*Printer)) *Printer {
	p := &Printer{}
	for _, f := range options {
		f(p)
	}
	return p
}

type Printer struct {
	output io.Writer
	report *Report
}

func (p *Printer) For(report *Report) *Printer {
	p.report = report
	return p
}

func (p *Printer) Print() {
	w := p.outOrStdout()

	for _, site := range p.report.Sites() {
		if err := site.Error(); err != nil {
			important().Fprintf(w, "report %q has error %q\n", site.Name(), err)
			continue
		}
		sort.Sort(pagesByLocation(site.Pages))
		for _, page := range site.Pages {
			last := len(page.Links) - 1
			colorize(page.StatusCode).Fprintf(w, head, page.StatusCode, page.Location)
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

func (p *Printer) outOrStdout() io.Writer {
	if p.output != nil {
		return p.output
	}
	return os.Stdout
}

func colorize(statusCode int) typewriter {
	var tw typewriter
	switch {
	case statusCode < 300:
		tw, _ = colors[success]
	case statusCode >= 300 && statusCode < 400:
		tw, _ = colors[warning]
	case statusCode >= 400:
		tw, _ = colors[danger]
	}
	if tw == nil {
		tw = defaultTypewriter(fmt.Fprintf)
	}
	return tw
}

func important() typewriter { return colorize(999) }

type typewriter interface {
	Fprintf(io.Writer, string, ...interface{}) (int, error)
}

type defaultTypewriter func(io.Writer, string, ...interface{}) (int, error)

func (fn defaultTypewriter) Fprintf(w io.Writer, format string, a ...interface{}) (int, error) {
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

func join(location, redirect string) string {
	if redirect == "" {
		return location
	}
	return location + " -> " + redirect
}
