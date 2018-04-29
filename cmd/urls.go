package cmd

import (
	"time"

	"github.com/briandowns/spinner"
	"github.com/kamilsk/check/http/availability"
	"github.com/spf13/cobra"
)

var urlsCmd = &cobra.Command{
	Use:   "urls",
	Short: "Check all internal URLs on availability",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		var spin = func() func() { return func() {} }
		verbose := asBool(cmd.Flag("verbose").Value)
		if !verbose {
			spin = func() func() {
				s := spinner.New(spinner.CharSets[34], 100*time.Millisecond)
				s.Writer = cmd.OutOrStderr()
				s.Start()
				return s.Stop
			}
		}
		stop := spin()
		report := availability.NewReport(
			availability.CrawlerForSites(availability.CrawlerColly(
				availability.CrawlerConfig{
					UserAgent: client(cmd),
					Verbose:   verbose,
					Output:    cmd.OutOrStderr(),
				},
			)),
		).
			For(args).
			Fill()
		stop()
		return availability.
			NewPrinter(availability.OutputForPrinting(cmd.OutOrStdout())).
			For(report).
			Print()
	},
}

func init() {
	urlsCmd.Flags().BoolP("verbose", "v", false, "turn on verbose mode")
}
