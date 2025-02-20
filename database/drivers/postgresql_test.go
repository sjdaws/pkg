package drivers_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"

	"github.com/sjdaws/pkg/database/drivers"
)

func TestPostgreSQL_GetDialector(t *testing.T) {
	t.Parallel()

	testcases := map[string]struct {
		options  *drivers.PostgreSQL
		expected string
	}{
		"name only": {
			options: &drivers.PostgreSQL{
				Name:            "test",
				SSLModeDisabled: false,
			},
			expected: "dbname=test",
		},
		"hostname": {
			options: &drivers.PostgreSQL{
				Host: "postgres.host",
				Name: "test",
			},
			expected: "dbname=test host=postgres.host",
		},
		"hostname and port": {
			options: &drivers.PostgreSQL{
				Host: "postgres.host",
				Name: "test",
				Port: 5432,
			},
			expected: "dbname=test host=postgres.host port=5432",
		},
		"invalid port": {
			options: &drivers.PostgreSQL{
				Host: "postgres.host",
				Name: "test",
				Port: -1,
			},
			expected: "dbname=test host=postgres.host",
		},
		"username": {
			options: &drivers.PostgreSQL{
				Name:     "test",
				Username: "root",
			},
			expected: "dbname=test user=root",
		},
		"username and password": {
			options: &drivers.PostgreSQL{
				Name:     "test",
				Password: "password",
				Username: "root",
			},
			expected: "dbname=test password=password user=root",
		},
		"sslmode disabled": {
			options: &drivers.PostgreSQL{
				Name:            "test",
				SSLModeDisabled: true,
			},
			expected: "dbname=test sslmode=disabled",
		},
		"full": {
			options: &drivers.PostgreSQL{
				Host:            "postgres.host",
				Name:            "test",
				Password:        "password",
				Port:            5432,
				SSLModeDisabled: true,
				Username:        "root",
			},
			expected: "dbname=test host=postgres.host password=password port=5432 sslmode=disabled user=root",
		},
	}

	for name, testcase := range testcases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			dialector, err := testcase.options.GetDialector()
			require.NoError(t, err)

			actual, ok := dialector.(postgres.Dialector)
			require.True(t, ok)
			assert.Equal(t, "TimeZone=UTC "+testcase.expected, actual.DSN)
		})
	}
}

func TestPostgreSQL_GetDialector_ErrInvalidDatabase(t *testing.T) {
	t.Parallel()

	options := &drivers.PostgreSQL{}

	dialector, err := options.GetDialector()
	require.Error(t, err)

	require.EqualError(t, err, "database name must be provided")
	assert.Nil(t, dialector)
}
