package cmd

import (
	"github.com/kamilsk/check/http/availability"
	"github.com/spf13/cobra"
)

var urlsCmd = &cobra.Command{
	Use:   "urls",
	Short: "Check all internal URLs on availability",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return availability.
			NewPrinter(
				availability.OutputForPrinting(cmd.OutOrStdout()),
			).
			For(
				availability.NewReport(
					availability.CrawlerForSites(availability.CrawlerColly(
						availability.CrawlerConfig{
							UserAgent: client(cmd),
							Verbose:   asBool(cmd.Flag("verbose").Value),
							Output:    cmd.OutOrStderr(),
						},
					)),
				).
					For(args).
					Fill(),
			).
			Print()
	},
}

func init() {
	urlsCmd.Flags().BoolP("verbose", "v", false, "turn on verbose mode")
}
