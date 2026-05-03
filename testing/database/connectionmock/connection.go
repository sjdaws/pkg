package connectionmock

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"sjdaws.com/pkg/database/drivers"
	"sjdaws.com/pkg/errors"
)

// DatabaseMock mock database instance.
type DatabaseMock struct {
	Fail bool
	orm  *gorm.DB
}

// Options various settings which can be toggled when creating a mock connection.
type Options struct {
	AlwaysFail bool
	DebugMode  bool
	FileBased  bool
}

// New mock database connection.
func New(t *testing.T, options ...Options) *DatabaseMock {
	t.Helper()

	// Process options
	var file bool

	if len(options) > 0 {
		file = options[0].FileBased
	}

	filename := ":memory:"

	if file {
		// Create temporary database
		filename = t.TempDir() + "/test.db"

		database, err := os.Create(filename)
		defer func() {
			_ = database.Close()
		}()

		require.NoError(t, err)
	}

	// Do an actual connection in case we want to do things like migrations
	driver := drivers.SQLite3{
		Filename: filename,
	}

	dialector, err := driver.GetDialector()
	require.NoError(t, err)

	return create(t, dialector, options...)
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

// create a DatabaseMock from a gorm.Dialector.
func create(t *testing.T, dialector gorm.Dialector, options ...Options) *DatabaseMock {
	t.Helper()

	// Process options
	var debug, fail bool

	if len(options) > 0 {
		debug = options[0].DebugMode
		fail = options[0].AlwaysFail
	}

	config := &gorm.Config{
		AllowGlobalUpdate:                        false,
		ClauseBuilders:                           nil,
		ConnPool:                                 nil,
		CreateBatchSize:                          0,
		DefaultTransactionTimeout:                0,
		DefaultContextTimeout:                    0,
		Dialector:                                nil,
		DisableAutomaticPing:                     false,
		DisableForeignKeyConstraintWhenMigrating: false,
		DisableNestedTransaction:                 false,
		DryRun:                                   false,
		FullSaveAssociations:                     false,
		IgnoreRelationshipsWhenMigrating:         false,
		Logger:                                   nil,
		NamingStrategy:                           nil,
		NowFunc:                                  nil,
		Plugins:                                  nil,
		PrepareStmt:                              false,
		PrepareStmtMaxSize:                       0,
		PrepareStmtTTL:                           0,
		PropagateUnscoped:                        false,
		QueryFields:                              false,
		SkipDefaultTransaction:                   false,
		TranslateError:                           false,
	}

	if debug {
		config.Logger = logger.Default.LogMode(logger.Info)
	}

	orm, err := gorm.Open(dialector, config)
	require.NoError(t, err)

	return &DatabaseMock{
		Fail: fail,
		orm:  orm,
	}
}
