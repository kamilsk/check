package availability

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"sort"
	"text/template"

	"github.com/fatih/color"
	"github.com/kamilsk/check/errors"
)

const (
	success = "success"
	warning = "warning"
	danger  = "danger"
)

var colors = map[string]*color.Color{
	success: color.New(color.FgWhite),
	warning: color.New(color.FgYellow),
	danger:  color.New(color.FgRed, color.Bold),
}

var entry = template.Must(template.New("entry").Parse(
	"[{{ .StatusCode }}] {{ .Location }}{{ with .Error }} -> ({{ . }}){{ end }}{{ with .Redirect }} -> {{ . }}{{ end }}",
))

// NewPrinter returns configured printer instance.
func NewPrinter(options ...func(*Printer)) *Printer {
	p := &Printer{}
	for _, f := range options {
		f(p)
	}
	return p
}

// OutputForPrinting sets up printer output.
func OutputForPrinting(output io.Writer) func(*Printer) {
	return func(p *Printer) {
		p.output = output
	}
}

// Reporter defines general behavior of report providers.
type Reporter interface {
	Sites() <-chan Site
}

// Printer represents a printer.
type Printer struct {
	output io.Writer
	report Reporter
}

// For prepares printer for passed report provider.
func (p *Printer) For(report Reporter) *Printer {
	p.report = report
	return p
}

// Print prints a report into the configured output.
// Stdout is used as a fallback if the output is not set up.
func (p *Printer) Print() error {
	var blob = [1024]byte{}
	w, buf := p.outOrStdout(), bytes.NewBuffer(blob[:0])
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
			{
				buf.Reset()
				entry.Execute(buf, page)
			}
			colorize(page.StatusCode).Fprintf(w, "%s\n", buf.String())
			sort.Sort(linksByStatusCode(page.Links))
			for i, link := range page.Links {
				{
					buf.Reset()
					entry.Execute(buf, link)
				}
				if i == last {
					colorize(link.StatusCode).Fprintf(w, "    └───%s\n", buf.String())
					continue
				}
				colorize(link.StatusCode).Fprintf(w, "    ├───%s\n", buf.String())
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
