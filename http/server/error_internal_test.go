package server

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"sjdaws.com/pkg/errors"
	"sjdaws.com/pkg/testing/logging/logmock"
)

func TestHandler(t *testing.T) {
	t.Parallel()

	logger := logmock.New()
	handler := errorHandler{
		logger: logger,
	}

	handler.log(errors.New("test error"), nil)

	assert.Equal(t, "error", logger.GetLastLevel())
	assert.Equal(t, "test error", logger.GetLastMessage())
}
