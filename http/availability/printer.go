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
		pages := byLocation(report.Pages())
		sort.Sort(pages)
		for _, page := range pages {
			fmt.Fprintf(w, "[%d] %s \n", page.StatusCode, page.Location)
			max := len(page.Links)
			for i, link := range page.Links {
				if max == i+1 {
					fmt.Fprintf(w, "└─── [%d] %s \n", link.StatusCode, link.Location)
					continue
				}
				fmt.Fprintf(w, "├─── [%d] %s \n", link.StatusCode, link.Location)
			}
		}
	}
}

type byLocation []*Page

func (l byLocation) Len() int {
	return len(l)
}

func (l byLocation) Less(i, j int) bool {
	return l[i].Link.Location < l[j].Link.Location
}

func (l byLocation) Swap(i, j int) {
	l[i], l[j] = l[j], l[i]
}
