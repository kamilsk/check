package main

import (
	"fmt"
	"io"
	"os"
	"runtime"

	"github.com/kamilsk/check/cmd"
	"github.com/spf13/cobra"

	_ "github.com/oxequa/grace"
	_ "github.com/spf13/viper"
)

const (
	success = 0
	failed  = 1
)

func main() { application{Cmd: cmd.RootCmd, Output: os.Stderr, Shutdown: os.Exit}.Run() }

type application struct {
	Cmd interface {
		AddCommand(...*cobra.Command)
		Execute() error
	}
	Output   io.Writer
	Shutdown func(code int)
}

// Run executes the application logic.
func (app application) Run() {
	app.Cmd.AddCommand(&cobra.Command{
		Use:   "version",
		Short: "Show application version",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Fprintf(app.Output,
				"Version %s (commit: %s, build date: %s, go version: %s, compiler: %s, platform: %s)\n",
				version, commit, date, runtime.Version(), runtime.Compiler, runtime.GOOS+"/"+runtime.GOARCH)
		},
		Version: version,
	})
	if err := app.Cmd.Execute(); err != nil {
		// so, when `issue` project will be ready
		// I have to integrate it to open GitHub issues
		// with stack trace from terminal
		app.Shutdown(failed)
	}
	app.Shutdown(success)
}
