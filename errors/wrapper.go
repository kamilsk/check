package errors

import (
	"fmt"

	core "errors"

	"github.com/pkg/errors"
)

func Errorf(format string, args ...interface{}) error {
	return errors.Errorf(format, args...)
}

func Recover(err *error) {
	if r := recover(); r != nil {
		if e, is := r.(error); is {
			*err = e
			return
		}
		*err = fmt.Errorf("%v", r)
	}
}

func Simple(message string) error {
	return core.New(message)
}

func StackTrace(err error) errors.StackTrace {
	type stackTracer interface {
		StackTrace() errors.StackTrace
	}

	if stack, is := err.(stackTracer); is {
		return stack.StackTrace()
	}
	return nil
}

func WithMessage(err error, message string) error {
	return errors.WithMessage(err, message)
}

func Wrapf(err error, format string, args ...interface{}) error {
	return errors.Wrapf(err, format, args...)
}
