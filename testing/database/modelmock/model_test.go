package modelmock_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/sjdaws/pkg/testing/database/modelmock"
)

func TestModelMock_TableName(t *testing.T) {
	t.Parallel()

	model := modelmock.ModelMock{}

	assert.Equal(t, "model_mocks", model.TableName())
}
