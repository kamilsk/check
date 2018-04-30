package availability

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"sort"

	"github.com/fatih/color"
	"github.com/kamilsk/check/errors"
)

const (
	success = "success"
	warning = "warning"
	danger  = "danger"

	head = "[%d] %s\n"
	body = "    ├─── [%d] %s\n"
	foot = "    └─── [%d] %s\n"
	sep  = " -> "
)

var colors = map[string]*color.Color{
	success: color.New(color.FgWhite),
	warning: color.New(color.FgYellow),
	danger:  color.New(color.FgRed, color.Bold),
}

func NewPrinter(options ...func(*Printer)) *Printer {
	p := &Printer{}
	for _, f := range options {
		f(p)
	}
	return p
}

func OutputForPrinting(output io.Writer) func(*Printer) {
	return func(p *Printer) {
		p.output = output
	}
}

type Reporter interface {
	Sites() <-chan Site
}

type Printer struct {
	output io.Writer
	report Reporter
}

func (p *Printer) For(report Reporter) *Printer {
	p.report = report
	return p
}

func (p *Printer) Print() error {
	w := p.outOrStdout()
	if p.report == nil {
		return errors.Simple("nothing to print")
	}
	for site := range p.report.Sites() {
		if site.Error != nil {
			critical().Fprintf(w, "report %q has error %q\n", site.Name, site.Error)
			if stack := errors.StackTrace(site.Error); stack != nil {
				critical().Fprintf(ioutil.Discard, "stack trace: %#+v\n", stack) // for future
			}
			continue
		}
		sort.Sort(pagesByLocation(site.Pages))
		for _, page := range site.Pages {
			last := len(page.Links) - 1
			colorize(page.StatusCode).Fprintf(w, head, page.StatusCode, page.Location)
			sort.Sort(linksByStatusCode(page.Links))
			for i, link := range page.Links {
				if i == last {
					colorize(link.StatusCode).Fprintf(w, foot, link.StatusCode, link.FullLocation(sep))
					continue
				}
				colorize(link.StatusCode).Fprintf(w, body, link.StatusCode, link.FullLocation(sep))
			}
		}
		if len(site.Problems) > 0 {
			critical().Fprintf(w, "found problems on the site %q\n", site.Name)
			for i, problem := range site.Problems {
				critical().Fprintf(w, "- [%d] %s `%+v`\n", i, problem.Message, problem.Context)
			}
		}
	}
	return nil
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
	case statusCode == 0:
		tw, _ = colors[danger]
	case statusCode >= 200 && statusCode < 300:
		tw, _ = colors[success]
	case statusCode >= 300 && statusCode < 400:
		tw, _ = colors[warning]
	case statusCode >= 400:
		tw, _ = colors[danger]
	}
	if tw == nil {
		tw = typewriterFunc(fmt.Fprintf)
	}
	return tw
}

func critical() typewriter { return colorize(0) }

type typewriter interface {
	Fprintf(io.Writer, string, ...interface{}) (int, error)
}

type typewriterFunc func(io.Writer, string, ...interface{}) (int, error)

func (fn typewriterFunc) Fprintf(w io.Writer, format string, a ...interface{}) (int, error) {
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
