package connectionmock_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/sjdaws/pkg/database"
	"github.com/sjdaws/pkg/testing/database/connectionmock"
)

func TestNew(t *testing.T) {
	t.Parallel()

	connection := connectionmock.New(t)

	assert.IsType(t, &database.Connection{}, connection)
}
