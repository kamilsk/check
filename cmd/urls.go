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
		var reports availability.Reports
		reports.Fill(args)
		printer := availability.NewPrinter(reports)
		printer.Print(cmd.OutOrStdout())
	},
}

func init() {
	urlsCmd.Flags().BoolP("verbose", "v", false, "turn on verbose mode")
}
