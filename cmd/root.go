package cmd

import "github.com/spf13/cobra"

// RootCmd is the entry point.
var RootCmd = &cobra.Command{Short: "check"}

func init() {
	RootCmd.AddCommand(urlsCmd)
}
