package cmd

import (
	"github.com/kamilsk/check/http/availability"
	"github.com/spf13/cobra"
)

var urlsCmd = &cobra.Command{
	Use:   "urls",
	Short: "Check all internal URLs on availability",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		availability.
			NewPrinter(
				availability.OutputForPrinting(cmd.OutOrStdout()),
			).
			For(
				availability.NewReport(
					availability.CrawlerForSites(availability.CrawlerColly(client())),
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
