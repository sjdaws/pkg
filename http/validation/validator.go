package validation

import (
	"github.com/go-playground/validator/v10"

	"github.com/sjdaws/pkg/errors"
)

// Validator interface.
type Validator interface {
	AddCustomValidation(key string, validation validator.Func, message MessageFunc) error
	AddFailureMessage(key string, message MessageFunc)
	Validate(target any) ([]string, error)
}

// Tester implementation of Validator.
type Tester struct {
	messages  map[string]func(field string, failure validator.FieldError) string
	validator *validator.Validate
}

// MessageFunc function for converting failures into better error messages.
type MessageFunc func(field string, failure validator.FieldError) string

// New create a new Validator.
func New() *Tester {
	messages := make(map[string]func(field string, failure validator.FieldError) string)

	return &Tester{
		messages:  messages,
		validator: validator.New(validator.WithRequiredStructEnabled()),
	}
}

// AddCustomValidation add a custom validator.
func (t *Tester) AddCustomValidation(key string, validation validator.Func, message MessageFunc) error {
	err := t.validator.RegisterValidation(key, validation)
	if err != nil {
		return errors.Wrap(err, "unable to add custom validator")
	}

	t.AddFailureMessage(key, message)

	return nil
}

// AddFailureMessage add a custom failure message handler.
func (t *Tester) AddFailureMessage(key string, message MessageFunc) {
	if message != nil {
		t.messages[key] = message
	}
}

// Validate a request against the rules on a struct.
func (t *Tester) Validate(target any) ([]string, error) {
	err := t.validator.Struct(target)
	if err != nil {
		// If error is due to failures, process them into messages the frontend can handle
		var failures validator.ValidationErrors
		if errors.As(err, &failures) {
			return t.processFailures(failures), nil
		}

		return nil, errors.Wrap(err, "unable to perform validation")
	}

	return nil, nil
}
