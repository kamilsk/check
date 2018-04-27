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

func client() string {
	var v *cobra.Command
	for _, cmd := range RootCmd.Commands() {
		if cmd.Use == "version" {
			v = cmd
			break
		}
	}
	if v != nil {
		return fmt.Sprintf("%s/%s", RootCmd.Short, v.Version)
	}
	return RootCmd.Short
}
