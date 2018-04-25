package availability

import (
	"fmt"
	"io"
	"sort"
)

func NewPrinter(reports []*Report) *Printer {
	return &Printer{reports: reports}
}

type Printer struct {
	reports []*Report
}

func (p *Printer) Print(w io.Writer) {
	for _, report := range p.reports {
		pages := pagesByLocation(report.Pages())
		sort.Sort(pages)
		for _, page := range pages {
			fmt.Fprintf(w, "[%d] %s \n", page.StatusCode, page.Location)
			last := len(page.Links) - 1
			sort.Sort(linksByStatusCode(page.Links))
			for i, link := range page.Links {
				if i == last {
					fmt.Fprintf(w, "└─── [%d] %s \n", link.StatusCode, link.Location)
					continue
				}
				fmt.Fprintf(w, "├─── [%d] %s \n", link.StatusCode, link.Location)
			}
		}
	}
}

type pagesByLocation []*Page

func (l pagesByLocation) Len() int { return len(l) }

func (l pagesByLocation) Less(i, j int) bool { return l[i].Link.Location < l[j].Link.Location }

func (l pagesByLocation) Swap(i, j int) { l[i], l[j] = l[j], l[i] }

type linksByStatusCode []*Link

func (l linksByStatusCode) Len() int { return len(l) }

func (l linksByStatusCode) Less(i, j int) bool { return l[i].StatusCode < l[j].StatusCode }

func (l linksByStatusCode) Swap(i, j int) { l[i], l[j] = l[j], l[i] }
