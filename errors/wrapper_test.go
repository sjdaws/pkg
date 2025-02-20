package errors_test

import (
	errs "errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/sjdaws/pkg/errors"
)

var (
	errPackage = errors.New("original error")
	errStdlib  = errs.New("original error")
)

func TestPublic(t *testing.T) {
	t.Parallel()

	var public errors.PublicError

	assert.Equal(t, "original error", errPackage.Error())

	err := errors.Public(errPackage, 1, "public error")

	require.ErrorAs(t, err, &public)
	assert.Equal(t, "original error", err.Error())
	assert.Equal(t, "original error", public.Error())
	assert.Equal(t, 1, public.Code())
	assert.Equal(t, []string{"public error"}, public.Errors())
}

func TestWrap(t *testing.T) {
	t.Parallel()

	assert.Equal(t, "original error", errPackage.Error())

	firstWrap := errors.Wrap(errPackage, "some context")

	require.ErrorAs(t, firstWrap, &errors.InternalError{})
	assert.Equal(t, "some context: original error", firstWrap.Error())

	secondWrap := errors.Wrap(firstWrap, "more context")

	require.ErrorAs(t, secondWrap, &errors.InternalError{})
	assert.Equal(t, "more context: original error", secondWrap.Error())
}

func TestError_Unwrap(t *testing.T) {
	t.Parallel()

	var unwrappable errors.InternalError

	err := errors.Wrap(errStdlib, "an error")

	assert.True(t, errors.As(err, &unwrappable))
	assert.Equal(t, errStdlib, unwrappable.Unwrap())

	err = errors.New("an error")

	assert.True(t, errors.As(err, &unwrappable))
	assert.NoError(t, unwrappable.Unwrap())
}

func TestPublic_Code(t *testing.T) {
	t.Parallel()

	var public errors.PublicError

	err := errors.Public(nil, 1, "public error")

	require.ErrorAs(t, err, &public)
	assert.Equal(t, 1, public.Code())
}

func TestPublic_Errors(t *testing.T) {
	t.Parallel()

	var public errors.PublicError

	err := errors.Public(nil, 1, "public error")

	require.ErrorAs(t, err, &public)
	assert.Equal(t, []string{"public error"}, public.Errors())
}
