package availability

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/url"
	"os"
	"sort"
	"text/template"

	"github.com/fatih/color"
	"go.octolab.org/unsafe"

	"github.com/kamilsk/check/errors"
)

const (
	shaded  = "shaded"
	success = "success"
	warning = "warning"
	danger  = "danger"
)

var base = template.Must(template.New("entry").Parse(`
{{- define "error" }}{{ with .Error }} -> ({{ . }}){{ end }}{{ end -}}
{{- define "redirect" }}{{ with .Redirect }} -> {{ . }}{{ end }}{{ end -}}
[{{ .StatusCode }}] {{ .Location }}{{ template "error" . }}{{ template "redirect" . -}}
`))

// NewPrinter returns configured printer instance.
func NewPrinter(options ...func(*Printer)) *Printer {
	p := &Printer{
		tpl:     template.Must(base.Clone()),
		decoder: func(origin string) string { return origin },
	}
	for _, f := range options {
		f(p)
	}
	return p
}

// ColorizeOutput sets the ink for the printer.
func ColorizeOutput(enabled bool) func(*Printer) {
	return func(p *Printer) {
		if enabled {
			p.ink = map[string]*color.Color{
				shaded:  color.New(color.FgHiBlack),
				success: color.New(color.FgWhite),
				warning: color.New(color.FgYellow),
				danger:  color.New(color.FgRed, color.Bold),
			}
		}
	}
}

// DecodeOutput sets `net/url.PathUnescape` as a decoder.
func DecodeOutput(enabled bool) func(*Printer) {
	return func(p *Printer) {
		if enabled {
			p.decoder = func(origin string) string {
				decoded, _ := url.PathUnescape(origin)
				return decoded
			}
		}
	}
}

// HideError prevents URL's error output.
func HideError(disabled bool) func(*Printer) {
	return func(p *Printer) {
		if disabled {
			unsafe.DoSilent(p.tpl.New("error").Parse("{{ with .Error }}{{/* ignore */}}{{ end }}"))
		}
	}
}

// HideRedirect prevents URL's redirect output.
func HideRedirect(disabled bool) func(*Printer) {
	return func(p *Printer) {
		if disabled {
			unsafe.DoSilent(p.tpl.New("redirect").Parse("{{ with .Redirect }}{{/* ignore */}}{{ end }}"))
		}
	}
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
	tpl     *template.Template
	output  io.Writer
	ink     map[string]*color.Color
	decoder func(string) string
	report  Reporter
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
			p.critical().Fprintf(w, "report %q has error %q\n", site.Name, site.Error)
			if stack := errors.StackTrace(site.Error); stack != nil {
				p.critical().Fprintf(ioutil.Discard, "stack trace: %#+v\n", stack) // for future
			}
			continue
		}
		sort.Sort(pagesByLocation(site.Pages))
		for _, page := range site.Pages {
			last := len(page.Links) - 1
			{
				buf.Reset()
				unsafe.Ignore(p.tpl.Execute(buf, page))
			}
			p.typewriter(page.Link).Fprintf(w, "%s\n", p.decoder(buf.String()))
			sort.Sort(linksByStatusCode(page.Links))
			for i, link := range page.Links {
				link := link
				{
					buf.Reset()
					unsafe.Ignore(p.tpl.Execute(buf, link))
				}
				if i == last {
					p.typewriter(&link).Fprintf(w, "    └───%s\n", p.decoder(buf.String()))
					continue
				}
				p.typewriter(&link).Fprintf(w, "    ├───%s\n", p.decoder(buf.String()))
			}
		}
		if len(site.Problems) > 0 {
			p.critical().Fprintf(w, "found problems on the site %q\n", site.Name)
			for i, problem := range site.Problems {
				p.critical().Fprintf(w, "- [%d] %s `%+v`\n", i, problem.Message, problem.Context)
			}
		}
	}
	return nil
}

func (p *Printer) critical() typewriter {
	return p.typewriter(nil)
}

func (p *Printer) outOrStdout() io.Writer {
	if p.output != nil {
		return p.output
	}
	return os.Stdout
}

func (p *Printer) typewriter(link *Link) typewriter {
	var (
		tw typewriter
		ok bool
	)
	switch {
	case link == nil:
		tw, ok = p.ink[danger]
	case link.StatusCode >= 200 && link.StatusCode < 300:
		if link.Internal {
			tw, ok = p.ink[shaded]
			break
		}
		tw, ok = p.ink[success]
	case link.StatusCode >= 300 && link.StatusCode < 400:
		tw, ok = p.ink[warning]
	case link.StatusCode >= 400:
		tw, ok = p.ink[danger]
	}
	if !ok || tw == nil {
		tw = typewriterFunc(fmt.Fprintf)
	}
	return tw
}

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

type linksByStatusCode []Link

func (l linksByStatusCode) Len() int { return len(l) }

func (l linksByStatusCode) Less(i, j int) bool { return l[i].StatusCode < l[j].StatusCode }

func (l linksByStatusCode) Swap(i, j int) { l[i], l[j] = l[j], l[i] }
