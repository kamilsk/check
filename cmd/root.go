package cmd

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
)

// RootCmd is the entry point.
var RootCmd = &cobra.Command{Use: "check"}

func init() {
	RootCmd.AddCommand(completionCmd, urlsCmd)
}

func asBool(value fmt.Stringer) bool {
	is, _ := strconv.ParseBool(value.String())
	return is
}

func client(cmd *cobra.Command) string {
	var version *cobra.Command
	if cmd.Parent() != nil {
		cmd = cmd.Parent()
	}
	for _, cmd := range cmd.Commands() {
		if cmd.Use == "version" {
			version = cmd
			break
		}
	}
	if version != nil {
		return fmt.Sprintf("%s/%s", cmd.Use, version.Version)
	}
	return cmd.Use
}
