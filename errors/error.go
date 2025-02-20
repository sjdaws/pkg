package errors

import (
	"errors"
	"fmt"
	"runtime"
	"strings"
)

// InternalError internal error type - can be logged but should not be exposed.
type InternalError struct {
	file     string
	line     int
	message  string
	previous error
}

// As attempts to set an error to target.
func As(err error, target any) bool {
	//goland:noinspection GoErrorsAs
	return errors.As(err, target)
}

// Is unwraps an error to determine if the error type is target.
func Is(err error, target error) bool {
	return errors.Is(err, target)
}

// New creates a new error message.
func New(message any, replacements ...any) error {
	_, file, line, _ := runtime.Caller(1)

	return InternalError{
		file:     file,
		line:     line,
		message:  fmt.Sprintf(fmt.Sprintf("%v", message), replacements...),
		previous: nil,
	}
}

// Error provides the last context message and original error message.
func (e InternalError) Error() string {
	original := e.original()

	if errors.Is(original, e) {
		return e.message
	}

	if e.message == "" {
		return original.Error()
	}

	return fmt.Sprintf("%s: %s", e.message, original.Error())
}

// Trace returns an error message with caller information.
func (e InternalError) Trace() string {
	stack := make([]string, 0)

	err := e
	previous := true

	for previous {
		var line string
		if err.file != "" {
			line = fmt.Sprintf("%s:%d: ", err.file, err.line)
		}

		stack = append(stack, fmt.Sprintf("%s%s", line, err.message))

		previous = errors.As(err.previous, &err)
	}

	return fmt.Sprintf("%s\n- %s", e.Error(), strings.Join(stack, "\n- "))
}

// Code returns the error code passed to Public.
func (p PublicError) Code() int {
	return p.code
}

// Errors returns the messages passed to Public.
func (p PublicError) Errors() []string {
	return p.errors
}
