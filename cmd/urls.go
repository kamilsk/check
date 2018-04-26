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
			NewPrinter(availability.Output(cmd.OutOrStdout())).
			For(availability.NewReport().
				For(args).
				Fill()).
			Print()
	},
}

func init() {
	urlsCmd.Flags().BoolP("verbose", "v", false, "turn on verbose mode")
}
