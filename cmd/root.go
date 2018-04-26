package cmd

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
)

// RootCmd is the entry point.
var RootCmd = &cobra.Command{Short: "check"}

func init() {
	RootCmd.AddCommand(urlsCmd)
}

func asBool(value fmt.Stringer) bool {
	is, _ := strconv.ParseBool(value.String())
	return is
}
