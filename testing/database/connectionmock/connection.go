package connectionmock

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/sjdaws/pkg/database/drivers"
	"github.com/sjdaws/pkg/errors"
)

type DatabaseMock struct {
	Fail bool
	orm  *gorm.DB
}

// New mock database connection.
func New(t *testing.T) *DatabaseMock {
	t.Helper()

	// Create temporary database
	filename := t.TempDir() + "/test.db"

	//nolint:gosec // Filename needs to be a reference to refer to later
	_, err := os.Create(filename)
	require.NoError(t, err)

	// Do an actual connection in case we want to do things like migrations
	options := drivers.SQLite3{
		Filename: filename,
	}

	dialector, err := options.GetDialector()
	require.NoError(t, err)

	config := &gorm.Config{
		AllowGlobalUpdate:                        false,
		ClauseBuilders:                           nil,
		ConnPool:                                 nil,
		CreateBatchSize:                          0,
		Dialector:                                nil,
		DisableAutomaticPing:                     false,
		DisableForeignKeyConstraintWhenMigrating: false,
		DisableNestedTransaction:                 false,
		DryRun:                                   false,
		FullSaveAssociations:                     false,
		IgnoreRelationshipsWhenMigrating:         false,
		Logger:                                   logger.Default.LogMode(logger.Info),
		NamingStrategy:                           nil,
		NowFunc:                                  nil,
		Plugins:                                  nil,
		PrepareStmt:                              false,
		PropagateUnscoped:                        false,
		QueryFields:                              false,
		SkipDefaultTransaction:                   false,
		TranslateError:                           false,
	}

	orm, err := gorm.Open(dialector, config)
	require.NoError(t, err)

	return &DatabaseMock{
		Fail: false,
		orm:  orm,
	}
}

// Migrate perform database migrations.
func (d *DatabaseMock) Migrate(model ...any) error {
	err := d.orm.AutoMigrate(model...)
	if err != nil || d.Fail {
		return errors.Wrap(err, "migration failed")
	}

	return nil
}

// ORM return the underlying ORM.
func (d *DatabaseMock) ORM() *gorm.DB {
	return d.orm
}

// Transaction create a new database transaction.
func (d *DatabaseMock) Transaction() *gorm.DB {
	transaction := d.orm.Begin()

	if d.Fail {
		transaction.Statement.ConnPool = nil
	}

	return transaction
}
