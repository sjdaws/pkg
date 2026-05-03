package database_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"sjdaws.com/pkg/testing/database/modelmock"
)

func TestModels_TableName(t *testing.T) {
	t.Parallel()

	model := modelmock.ModelMock{}

	assert.Equal(t, "model_mocks", model.TableName())
}
