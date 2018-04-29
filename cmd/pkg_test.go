package cmd

import (
	"bytes"
	"testing"

	"fmt"

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
	assert.NoError(t, urlsCmd.RunE(urlsCmd, []string{site.URL + "/"}))
	assert.Contains(t, buf.String(), fmt.Sprintf("[200] %s/", site.URL))
}
