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

	mock, ok := connection.(*connectionmock.DatabaseMock)
	require.True(t, ok)

	mock.Fail = true

	err := mock.Migrate()
	require.Error(t, err)

	require.EqualError(t, err, "migration failed")
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

	mock, ok := connection.(*connectionmock.DatabaseMock)
	require.True(t, ok)

	mock.Fail = true

	transaction := mock.Transaction()

	assert.IsType(t, &gorm.DB{}, transaction)

	err := transaction.Commit().Error
	require.Error(t, err)

	require.EqualError(t, err, "invalid transaction")
}
