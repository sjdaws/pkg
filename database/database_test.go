package database_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/sjdaws/pkg/database"
)

func TestNew(t *testing.T) {
	t.Parallel()

	connection, err := database.New(false, "sqlite", "", ":memory:", "", 0, "", "", "")
	require.NoError(t, err)

	assert.NotNil(t, connection)
}

func TestNew_ErrInvalidDialector(t *testing.T) {
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

	for _, driver := range testcases {
		t.Run(driver, func(t *testing.T) {
			t.Parallel()

			connection, err := database.New(false, driver, "", "", "", 0, "", "", "")
			require.Error(t, err)

			require.EqualError(t, err, "unable to create dialector: database name must be provided")
			assert.Nil(t, connection)
		})
	}
}

func TestNew_ErrInvalidDriver(t *testing.T) {
	t.Parallel()

	connection, err := database.New(false, "invalid", "", "", "", 0, "", "", "")
	require.Error(t, err)

	require.EqualError(t, err, "unsupported database type requested: invalid")
	assert.Nil(t, connection)
}

func TestNew_ErrOpenDatabaseFailure(t *testing.T) {
	t.Parallel()

	connection, err := database.New(false, "sqlite", "", "\\/:", "", 0, "", "", "")
	require.Error(t, err)

	require.EqualError(t, err, "unable to open connection to database: unable to open database file: out of memory (14)")
	assert.Nil(t, connection)
}
