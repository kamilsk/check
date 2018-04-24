package main

import (
	"io"
	"os"
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
	app.Shutdown(success)
	return
}
