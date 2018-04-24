package cmd

import "github.com/spf13/cobra"

var urlsCmd = &cobra.Command{
	Use:   "urls",
	Short: "Check all attainable URLs",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Println("urls... urls")
	},
}
