package errors

import (
	core "errors"
	"fmt"

	"github.com/pkg/errors"
)

// Errorf is a proxy for `github.com/pkg/errors.Errorf`.
func Errorf(format string, args ...interface{}) error {
	return errors.Errorf(format, args...)
}

// Recover recovers execution flow and sets error to the passed error pointer.
func Recover(err *error) {
	if r := recover(); r != nil {
		if e, is := r.(error); is {
			*err = e
			return
		}
		*err = fmt.Errorf("%v", r)
	}
}

// Simple returns text-based error without stack trace.
func Simple(message string) error {
	return core.New(message)
}

// StackTrace tries to extract stack trace from the error
// or returns nil if it can't.
func StackTrace(err error) errors.StackTrace {
	type stackTracer interface {
		StackTrace() errors.StackTrace
	}

	if stack, is := err.(stackTracer); is {
		return stack.StackTrace()
	}
	return nil
}

// WithMessage is a proxy for `github.com/pkg/errors.WithMessage`.
func WithMessage(err error, message string) error {
	return errors.WithMessage(err, message)
}

// Wrapf is a proxy for `github.com/pkg/errors.Wrapf`.
func Wrapf(err error, format string, args ...interface{}) error {
	return errors.Wrapf(err, format, args...)
}
