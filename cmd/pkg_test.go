package cmd

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestCompletion(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	completionCmd.SetOutput(buf)
	defer completionCmd.SetOutput(nil)
	{
		assert.Error(t, completionCmd.Args(completionCmd, nil))
		assert.Error(t, completionCmd.Args(completionCmd, []string{"shell"}))
		assert.NoError(t, completionCmd.Args(completionCmd, []string{bashFormat}))
		assert.NoError(t, completionCmd.Args(completionCmd, []string{zshFormat}))
	}
	{
		buf.Reset()
		assert.NoError(t, completionCmd.RunE(completionCmd, []string{bashFormat}))
		assert.Contains(t, buf.String(), "# bash completion for check")
	}
	{
		buf.Reset()
		assert.NoError(t, completionCmd.RunE(completionCmd, []string{zshFormat}))
		assert.Contains(t, buf.String(), "#compdef check")
	}
}

func TestURLs(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	site, closer := site()
	urlsCmd.SetOutput(buf)
	defer closer()
	defer urlsCmd.SetOutput(nil)
	{
		buf.Reset()
		assert.NoError(t, urlsCmd.RunE(urlsCmd, []string{site.URL + "/"}))
		assert.Contains(t, buf.String(), fmt.Sprintf("[200] %s/", site.URL))
	}
	{
		buf.Reset()
		verbose := urlsCmd.Flag("verbose")
		verbose.Value.Set("true")
		assert.NoError(t, urlsCmd.RunE(urlsCmd, []string{site.URL + "/"}))
		verbose.Value.Set(verbose.DefValue)
	}
}

func Test_client(t *testing.T) {
	tests := []struct {
		name     string
		cmd      func() *cobra.Command
		expected string
	}{
		{"without parent", func() *cobra.Command {
			return &cobra.Command{Use: "child"}
		}, "child"},
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
