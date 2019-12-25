package cmd

import (
	"time"

	"github.com/briandowns/spinner"
	"github.com/spf13/cobra"

	"github.com/kamilsk/check/http/availability"
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
			NewPrinter(
				availability.ColorizeOutput(!asBool(cmd.Flag("no-color").Value)),
				availability.DecodeOutput(asBool(cmd.Flag("decode").Value)),
				availability.HideError(asBool(cmd.Flag("no-error").Value)),
				availability.HideRedirect(asBool(cmd.Flag("no-redirect").Value)),
				availability.OutputForPrinting(cmd.OutOrStdout()),
			).
			For(report).
			Print()
	},
}

func init() {
	urlsCmd.Flags().BoolP("decode", "d", false, "decode URLs")
	urlsCmd.Flags().Bool("no-color", false, "disable colorized output")
	urlsCmd.Flags().Bool("no-error", false, "do not show URL's error")
	urlsCmd.Flags().Bool("no-redirect", false, "do not show URL's redirect")
	urlsCmd.Flags().BoolP("verbose", "v", false, "turn on verbose mode")
}
