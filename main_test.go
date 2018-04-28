package main

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestApplication_Run(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	app := application{Output: buf}

	tests := []struct {
		name string
		cmd  func() interface {
			AddCommand(...*cobra.Command)
			Execute() error
		}
		expected int
	}{
		{
			name: "success run",
			cmd: func() interface {
				AddCommand(...*cobra.Command)
				Execute() error
			} {
				cmd := &CmdMock{}
				cmd.On("AddCommand", mock.Anything)
				cmd.On("Execute").Return(nil)
				return cmd
			},
			expected: success,
		},
		{
			name: "failed run",
			cmd: func() interface {
				AddCommand(...*cobra.Command)
				Execute() error
			} {
				cmd := &CmdMock{}
				cmd.On("AddCommand", mock.Anything)
				cmd.On("Execute").Return(fmt.Errorf("mocking"))
				return cmd
			},
			expected: failed,
		},
	}
	for _, test := range tests {
		tc := test
		t.Run(test.name, func(t *testing.T) {
			buf.Reset()
			app.Cmd = tc.cmd()
			app.Shutdown = func(code int) { panic(assert.Equal(t, tc.expected, code)) }
			assert.Panics(t, func() { app.Run() })
			assert.Contains(t, buf.String(), "Version dev")
		})
	}
}
