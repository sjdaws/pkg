package connectionmock

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	"github.com/sjdaws/pkg/database"
	"github.com/sjdaws/pkg/database/drivers"
	"github.com/sjdaws/pkg/errors"
)

type DatabaseMock struct {
	Fail bool
	orm  *gorm.DB
}

// New mock database connection.
func New(t *testing.T) database.Connection {
	t.Helper()

	// Create temporary database
	filename := t.TempDir() + "/test.db"

	_, err := os.Create(filename)
	require.NoError(t, err)

	// Do an actual connection in case we want to do things like migrations
	options := drivers.SQLite3{
		Filename: filename,
	}

	dialector, err := options.GetDialector()
	require.NoError(t, err)

	orm, err := gorm.Open(dialector, nil)
	require.NoError(t, err)

	return &DatabaseMock{
		Fail: false,
		orm:  orm,
	}
}

// Migrate perform database migrations.
func (d DatabaseMock) Migrate(_ ...any) error {
	if d.Fail {
		return errors.New("migration failed")
	}

	return nil
}

// Transaction create a new database transaction.
func (d DatabaseMock) Transaction() *gorm.DB {
	transaction := d.orm.Begin()

	if d.Fail {
		transaction.Statement.ConnPool = nil
	}

	return transaction
}
