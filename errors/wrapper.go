package errors

import (
	"errors"
	"fmt"
	"runtime"
)

// PublicError errors which can be exposed to end users.
type PublicError struct {
	InternalError
	code   int
	errors []string
}

// Public wraps an error in public error message.
func Public(err error, code int, message ...string) error {
	_, file, line, _ := runtime.Caller(1)

	internal := InternalError{
		file:     file,
		line:     line,
		message:  "",
		previous: err,
	}

	// If err is an InternalError, use directly
	var previous InternalError

	if errors.As(err, &previous) {
		internal = previous
	}

	return PublicError{
		InternalError: internal,
		code:          code,
		errors:        message,
	}
}

// Wrap wraps an error with additional context.
func Wrap(err error, context any, replacements ...any) error {
	_, file, line, _ := runtime.Caller(1)

	return InternalError{
		file:     file,
		line:     line,
		message:  fmt.Sprintf(fmt.Sprintf("%v", context), replacements...),
		previous: err,
	}
}

// Unwrap returns the next error in the error stack.
func (e InternalError) Unwrap() error {
	return e.previous
}

// original continuously unwraps an error until the original error is found.
func (e InternalError) original() error {
	next := e.Unwrap()

	if next == nil {
		return e
	}

	var err InternalError
	if errors.As(next, &err) {
		return err.original()
	}

	return next
}
