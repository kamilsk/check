module github.com/kamilsk/check/tools

go 1.15

require (
	github.com/golang/mock v1.4.4
	github.com/golangci/golangci-lint v1.32.0
	github.com/kamilsk/egg v0.0.16
	github.com/kyoh86/looppointer v0.1.6
	golang.org/x/tools v0.5.0
)

replace github.com/izumin5210/gex => github.com/kamilsk/gex v0.6.0-e4

replace golang.org/x/tools => github.com/kamilsk/go-tools v0.0.5
