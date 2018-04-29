package cmd

import (
	"bytes"
	"testing"

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
