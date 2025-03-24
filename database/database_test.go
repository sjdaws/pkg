package database_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	"github.com/sjdaws/pkg/database"
	"github.com/sjdaws/pkg/testing/database/modelmock"
)

func TestConnect(t *testing.T) {
	t.Parallel()

	connection, err := database.Connect(false, "sqlite", "", ":memory:", "", 0, "", "", "")
	require.NoError(t, err)

	assert.IsType(t, &database.Database{}, connection)
	assert.Implements(t, (*database.Connection)(nil), connection)
}

func TestConnect_ErrInvalidDialector(t *testing.T) {
	t.Parallel()

	testcases := []string{
		"mariadb",
		"mysql",
		"postgres",
		"postgresql",
		"sqlite",
		"sqlite3",
		"sqlserver",
	}

	for _, name := range testcases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			connection, err := database.Connect(false, name, "", "", "", 0, "", "", "")
			require.Error(t, err)

			require.EqualError(t, err, "unable to create dialector: database name must be provided")
			assert.Nil(t, connection)
		})
	}
}

func TestConnect_ErrInvalidDriver(t *testing.T) {
	t.Parallel()

	connection, err := database.Connect(false, "invalid", "", "", "", 0, "", "", "")
	require.Error(t, err)

	require.EqualError(t, err, "unsupported database type requested: invalid")
	assert.Nil(t, connection)
}

func TestConnect_ErrOpenDatabaseFailure(t *testing.T) {
	t.Parallel()

	connection, err := database.Connect(false, "sqlite", "", "\\/:", "", 0, "", "", "")
	require.Error(t, err)

	require.EqualError(t, err, "unable to open connection to database: unable to open database file: out of memory (14)")
	assert.Nil(t, connection)
}

func TestNew_ErrMigrationFailure(t *testing.T) {
	t.Parallel()

	// Enforce read only database
	connection, err := database.Connect(false, "sqlite", "", ":memory:?_pragma=query_only(true)", "", 0, "", "", "")
	require.NoError(t, err)

	err = connection.Migrate(modelmock.ModelMock{})
	require.Error(t, err)

	require.EqualError(t, err, "unable to invoke database migrations: attempt to write a readonly database (8)")
}

func TestRepository(t *testing.T) {
	t.Parallel()

	connection, err := database.Connect(false, "sqlite", "", ":memory:", "", 0, "", "", "")
	require.NoError(t, err)

	instance := database.Repository[modelmock.ModelMock](connection)

	assert.Implements(t, (*database.Persister[modelmock.ModelMock])(nil), instance)
}

func TestTransaction(t *testing.T) {
	t.Parallel()

	connection, err := database.Connect(false, "sqlite", "", ":memory:", "", 0, "", "", "")
	require.NoError(t, err)

	transaction := connection.Transaction()

	assert.NotEqual(t, connection, transaction)
	assert.IsType(t, &gorm.DB{}, transaction)
}
