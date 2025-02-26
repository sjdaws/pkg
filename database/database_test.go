package database_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/sjdaws/pkg/database"
	"github.com/sjdaws/pkg/testing/database/modelmock"
)

func TestConnect(t *testing.T) {
	t.Parallel()

	connection, err := database.Connect(false, "sqlite", "", ":memory:", "", 0, "", "", "")
	require.NoError(t, err)

	assert.IsType(t, &database.Connection{}, connection)
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

	// Get currently running process
	exe, err := os.Executable()
	require.NoError(t, err)

	// Using the current process as a database will fail migrations
	connection, err := database.Connect(false, "sqlite", "", exe, "", 0, "", "", "")
	require.NoError(t, err)

	err = connection.Migrate(modelmock.ModelMock{})
	require.Error(t, err)

	require.EqualError(t, err, "unable to invoke database migrations: file is not a database (26)")
}

func TestRepository(t *testing.T) {
	t.Parallel()

	connection, err := database.Connect(false, "sqlite", "", ":memory:", "", 0, "", "", "")
	require.NoError(t, err)

	instance := database.Repository[modelmock.ModelMock](connection)

	assert.Implements(t, (*database.Persister[modelmock.ModelMock])(nil), instance)
}
