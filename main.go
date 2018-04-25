package main

import (
	"fmt"
	"io"
	"os"
	"runtime"

	"github.com/kamilsk/check/cmd"
	"github.com/spf13/cobra"

	_ "github.com/spf13/viper"
	_ "github.com/stretchr/testify"
)

const (
	success = 0
	failed  = 1
)

func main() { application{Stderr: os.Stderr, Stdout: os.Stdout, Shutdown: os.Exit}.Run() }

type application struct {
	Stderr, Stdout io.Writer
	Shutdown       func(code int)
}

// Run executes the application logic.
func (app application) Run() {
	cmd.RootCmd.AddCommand(&cobra.Command{
		Use:   "version",
		Short: "Show application version",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Fprintf(app.Stderr,
				"Version %s (commit: %s, build date: %s, go version: %s, compiler: %s, platform: %s)\n",
				version, commit, date, runtime.Version(), runtime.Compiler, runtime.GOOS+"/"+runtime.GOARCH)
		},
	})
	cmd.RootCmd.SetOutput(app.Stderr)
	if err := cmd.RootCmd.Execute(); err != nil {
		fmt.Fprintln(app.Stderr, err)
		app.Shutdown(failed)
	}
	app.Shutdown(success)
}
