package cmd

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func Test_client(t *testing.T) {
	tests := []struct {
		name     string
		cmd      func() *cobra.Command
		expected string
	}{
		{"without parent", func() *cobra.Command { return &cobra.Command{Use: "child"} }, "child"},
		{"without version", func() *cobra.Command {
			parent := &cobra.Command{Use: "parent"}
			child := &cobra.Command{Use: "child"}
			parent.AddCommand(child)
			return child
		}, "parent"},
		{"with version", func() *cobra.Command {
			parent := &cobra.Command{Use: "parent"}
			child := &cobra.Command{Use: "child"}
			version := &cobra.Command{Use: "version", Version: "ver."}
			parent.AddCommand(child, version)
			return child
		}, "parent/ver."},
	}
	for _, test := range tests {
		tc := test
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, tc.expected, client(tc.cmd()))
		})
	}
}
