package drivers_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlserver"

	"github.com/sjdaws/pkg/database/drivers"
)

func TestSQLServer_GetDialector(t *testing.T) {
	t.Parallel()

	testcases := map[string]struct {
		options  *drivers.SQLServer
		expected string
	}{
		"name only": {
			options: &drivers.SQLServer{
				Name: "test",
			},
			expected: "sqlserver://@?database=test",
		},
		"hostname": {
			options: &drivers.SQLServer{
				Host: "sqlserver.host",
				Name: "test",
			},
			expected: "sqlserver://@sqlserver.host?database=test",
		},
		"hostname and port": {
			options: &drivers.SQLServer{
				Host: "sqlserver.host",
				Name: "test",
				Port: 9930,
			},
			expected: "sqlserver://@sqlserver.host:9930?database=test",
		},
		"invalid port": {
			options: &drivers.SQLServer{
				Host: "sqlserver.host",
				Name: "test",
				Port: -1,
			},
			expected: "sqlserver://@sqlserver.host?database=test",
		},
		"username": {
			options: &drivers.SQLServer{
				Name:     "test",
				Username: "root",
			},
			expected: "sqlserver://root@?database=test",
		},
		"username and password": {
			options: &drivers.SQLServer{
				Name:     "test",
				Password: "password",
				Username: "root",
			},
			expected: "sqlserver://root:password@?database=test",
		},
		"full": {
			options: &drivers.SQLServer{
				Host:     "sqlserver.host",
				Name:     "test",
				Password: "password",
				Port:     9930,
				Username: "root",
			},
			expected: "sqlserver://root:password@sqlserver.host:9930?database=test",
		},
	}

	for name, testcase := range testcases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			dialector, err := testcase.options.GetDialector()
			require.NoError(t, err)

			actual, ok := dialector.(sqlserver.Dialector)
			require.True(t, ok)
			assert.Equal(t, testcase.expected, actual.DSN)
		})
	}
}

func TestSQLServer_GetDialector_ErrInvalidDatabase(t *testing.T) {
	t.Parallel()

	options := &drivers.SQLServer{}

	dialector, err := options.GetDialector()
	require.Error(t, err)

	require.EqualError(t, err, "database name must be provided")
	assert.Nil(t, dialector)
}
