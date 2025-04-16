package connectionmock_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	"github.com/sjdaws/pkg/database"
	"github.com/sjdaws/pkg/testing/database/connectionmock"
)

func TestNew(t *testing.T) {
	t.Parallel()

	connection := connectionmock.New(t)

	assert.IsType(t, &connectionmock.DatabaseMock{}, connection)
	assert.Implements(t, (*database.Connection)(nil), connection)
}

func TestConnection_Migrate(t *testing.T) {
	t.Parallel()

	connection := connectionmock.New(t)

	err := connection.Migrate()
	require.NoError(t, err)
}

func TestConnection_Migrate_Error(t *testing.T) {
	t.Parallel()

	connection := connectionmock.New(t)
	connection.Fail = true

	err := connection.Migrate()
	require.Error(t, err)

	require.EqualError(t, err, "migration failed")
}

func TestConnection_ORM(t *testing.T) {
	t.Parallel()

	connection := connectionmock.New(t)

	assert.IsType(t, &gorm.DB{}, connection.ORM())
}

func TestConnection_Transaction(t *testing.T) {
	t.Parallel()

	connection := connectionmock.New(t)

	transaction := connection.Transaction()

	assert.IsType(t, &gorm.DB{}, transaction)
}

func TestConnection_Transaction_Error(t *testing.T) {
	t.Parallel()

	connection := connectionmock.New(t)
	connection.Fail = true

	transaction := connection.Transaction()

	assert.IsType(t, &gorm.DB{}, transaction)

	err := transaction.Commit().Error
	require.Error(t, err)

	require.EqualError(t, err, "invalid transaction")
}
