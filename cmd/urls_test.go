package cmd

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestURLs(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	cmd := urlsCmd
	cmd.SetOutput(buf)
	defer cmd.SetOutput(nil)
	site, closer := site()
	defer closer()
	{
		buf.Reset()
		assert.NoError(t, cmd.RunE(cmd, []string{site.URL + "/"}))
		assert.Contains(t, buf.String(), fmt.Sprintf("[200] %s/", site.URL))
	}
	{
		buf.Reset()
		verbose := cmd.Flag("verbose")
		verbose.Value.Set("true")
		assert.NoError(t, cmd.RunE(cmd, []string{site.URL + "/"}))
		verbose.Value.Set(verbose.DefValue)
	}
}
