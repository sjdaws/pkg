package drivers_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/mysql"

	"github.com/sjdaws/pkg/database/drivers"
)

func TestMySQL_GetDialector(t *testing.T) {
	t.Parallel()

	testcases := map[string]struct {
		options  *drivers.MySQL
		expected string
	}{
		"name only": {
			options: &drivers.MySQL{
				Name: "test",
			},
			expected: "@/test",
		},
		"hostname": {
			options: &drivers.MySQL{
				Host: "mysql.host",
				Name: "test",
			},
			expected: "@tcp(mysql.host)/test",
		},
		"hostname and port": {
			options: &drivers.MySQL{
				Host: "mysql.host",
				Name: "test",
				Port: 3306,
			},
			expected: "@tcp(mysql.host:3306)/test",
		},
		"invalid port": {
			options: &drivers.MySQL{
				Host: "mysql.host",
				Name: "test",
				Port: -1,
			},
			expected: "@tcp(mysql.host)/test",
		},
		"socket": {
			options: &drivers.MySQL{
				Name:   "test",
				Socket: "/var/lib/mysql.sock",
			},
			expected: "@unix(/var/lib/mysql.sock)/test",
		},
		"hostname and socket": {
			options: &drivers.MySQL{
				Host:   "mysql.host",
				Name:   "test",
				Socket: "/var/lib/mysql.sock",
			},
			expected: "@unix(/var/lib/mysql.sock)/test",
		},
		"username": {
			options: &drivers.MySQL{
				Name:     "test",
				Username: "root",
			},
			expected: "root@/test",
		},
		"username and password": {
			options: &drivers.MySQL{
				Name:     "test",
				Password: "password",
				Username: "root",
			},
			expected: "root:password@/test",
		},
		"full": {
			options: &drivers.MySQL{
				Host:     "mysql.host",
				Name:     "test",
				Password: "password",
				Port:     3306,
				Socket:   "/var/lib/mysql.sock",
				Username: "root",
			},
			expected: "root:password@unix(/var/lib/mysql.sock)/test",
		},
	}

	for name, testcase := range testcases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			dialector, err := testcase.options.GetDialector()
			require.NoError(t, err)

			actual, ok := dialector.(mysql.Dialector)
			require.True(t, ok)
			assert.Equal(t, testcase.expected+"?charset=utf8mb4&parseTime=True&loc=UTC", actual.DSN)
		})
	}
}

func TestMySQL_GetDialector_ErrInvalidDatabase(t *testing.T) {
	t.Parallel()

	options := &drivers.MySQL{}

	dialector, err := options.GetDialector()
	require.Error(t, err)

	require.EqualError(t, err, "database name must be provided")
	assert.Nil(t, dialector)
}
