package validation_test

import (
	"fmt"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/sjdaws/pkg/http/validation"
)

func TestNew(t *testing.T) {
	t.Parallel()

	instance := validation.New()

	assert.IsType(t, &validation.Tester{}, instance)
	assert.Implements(t, (*validation.Validator)(nil), instance)
}

func TestAddCustomValidation(t *testing.T) {
	t.Parallel()

	handler := func(_ validator.FieldLevel) bool {
		return true
	}

	instance := validation.New()
	err := instance.AddCustomValidation("test", handler, nil)
	require.NoError(t, err)
}

func TestAddCustomValidation_InvalidValidation(t *testing.T) {
	t.Parallel()

	instance := validation.New()
	err := instance.AddCustomValidation("test", nil, nil)
	require.Error(t, err)

	require.EqualError(t, err, "unable to add custom validator: function cannot be empty")
}

func TestAddFailureMessage(t *testing.T) {
	t.Parallel()

	handler := func(_ string, _ validator.FieldError) string {
		return ""
	}

	instance := validation.New()
	instance.AddFailureMessage("test", handler)
}

func TestValidate_Pass(t *testing.T) {
	t.Parallel()

	type Data struct {
		Required string `validate:"required"`
	}

	target := Data{
		Required: "test",
	}

	instance := validation.New()

	failures, err := instance.Validate(&target)
	require.NoError(t, err)

	assert.Empty(t, failures)
}

func TestValidate_Failure(t *testing.T) {
	t.Parallel()

	type Data struct {
		CustomField string `validate:"custom"`
		Required    string `validate:"required"`
	}

	target := Data{
		CustomField: "invalid",
	}

	handler := func(field validator.FieldLevel) bool {
		return field.Field().String() != "invalid"
	}

	message := func(field string, failure validator.FieldError) string {
		return fmt.Sprintf("%s '%s' must not equal invalid", field, failure.Value())
	}

	instance := validation.New()

	err := instance.AddCustomValidation("custom", handler, message)
	require.NoError(t, err)

	failures, err := instance.Validate(&target)
	require.NoError(t, err)

	require.Len(t, failures, 2)
	assert.Equal(t, "Custom field 'invalid' must not equal invalid", failures[0])
	assert.Equal(t, "Required is required", failures[1])
}

func TestValidate_InvalidTarget(t *testing.T) {
	t.Parallel()

	target := "test"

	instance := validation.New()

	failures, err := instance.Validate(&target)
	require.Error(t, err)

	require.EqualError(t, err, "unable to perform validation: validator: (nil *string)")
	require.Nil(t, failures)
}
