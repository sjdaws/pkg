package drivers_test

import (
	"testing"

	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/sjdaws/pkg/database/drivers"
)

func TestSQLite3_GetDialector(t *testing.T) {
	t.Parallel()

	options := &drivers.SQLite3{
		Filename: ":memory:",
	}

	dialector, err := options.GetDialector()
	require.NoError(t, err)

	actual, ok := dialector.(sqlite.Dialector)
	require.True(t, ok)

	assert.Equal(t, ":memory:?_pragma=foreign_keys(1)", actual.DSN)
}

func TestSQLite3_GetDialector_ErrInvalidDatabase(t *testing.T) {
	t.Parallel()

	options := &drivers.SQLite3{}

	dialector, err := options.GetDialector()
	require.Error(t, err)

	require.EqualError(t, err, "database name must be provided")
	assert.Nil(t, dialector)
}
