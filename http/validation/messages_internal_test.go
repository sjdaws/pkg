package validation

import (
	"reflect"
	"testing"

	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
)

func TestTester_processFailures(t *testing.T) {
	t.Parallel()

	testcases := map[string]struct {
		expected string
		failure  validator.FieldError
	}{
		"email": {
			expected: "Field 'invalid' is not a valid email address",
			failure:  createError(t, "email"),
		},
		"endswith": {
			expected: "Field must end with 'value'",
			failure:  createError(t, "endswith"),
		},
		"oneof": {
			expected: "Field must be one of: 'value'",
			failure:  createError(t, "oneof"),
		},
		"required": {
			expected: "Field is required",
			failure:  createError(t, "required"),
		},
		"uuid": {
			expected: "Field must be a valid uuid",
			failure:  createError(t, "uuid"),
		},
		"default": {
			expected: "Predefined error message",
			failure:  createError(t, "default"),
		},
	}

	for name, testcase := range testcases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			instance := &Tester{}
			actual := instance.processFailures([]validator.FieldError{testcase.failure})

			assert.Equal(t, []string{testcase.expected}, actual)
		})
	}
}

func createError(t *testing.T, tag string) *validationError {
	t.Helper()

	return &validationError{
		message: "invalid",
		tag:     tag,
	}
}

type validationError struct {
	message string
	tag     string
}

func (v *validationError) ActualTag() string {
	return ""
}

func (v *validationError) Error() string {
	return "Predefined error message"
}

func (v *validationError) Field() string {
	return "field"
}

func (v *validationError) Kind() reflect.Kind {
	return reflect.TypeOf(v).Kind()
}

func (v *validationError) Namespace() string {
	return ""
}

func (v *validationError) Param() string {
	return "value"
}

func (v *validationError) StructField() string {
	return ""
}

func (v *validationError) StructNamespace() string {
	return ""
}

func (v *validationError) Tag() string {
	return v.tag
}

func (v *validationError) Translate(_ ut.Translator) string {
	return ""
}

//nolint:ireturn // Return type defined by validator.FieldError interface.
func (v *validationError) Type() reflect.Type {
	return reflect.TypeOf(v)
}

//nolint:ireturn // Return type defined by validator.FieldError interface.
func (v *validationError) Value() interface{} {
	return v.message
}
