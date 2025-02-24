package connectionmock

import (
	"database/sql/driver"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// Exec options for MockExec.
type Exec struct {
	Direct    bool
	Error     error
	Query     string
	QueryArgs []driver.Value
	Raw       bool
}

// Select options for MockSelect.
type Select struct {
	Error     error
	Query     string
	QueryArgs []driver.Value
	Raw       bool
	Rows      *sqlmock.Rows
}

// New mock database and mock handler.
func New(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
	t.Helper()

	database, mock, err := sqlmock.New()
	require.NoError(t, err)

	mock.ExpectQuery("SELECT VERSION()").
		WillReturnRows(mock.NewRows([]string{"mock"}).AddRow("mock"))

	configuration := mysql.Config{
		DriverName:                    "",
		ServerVersion:                 "",
		DSN:                           "",
		DSNConfig:                     nil,
		Conn:                          database,
		SkipInitializeWithVersion:     false,
		DefaultStringSize:             0,
		DefaultDatetimePrecision:      nil,
		DisableWithReturning:          false,
		DisableDatetimePrecision:      false,
		DontSupportRenameIndex:        false,
		DontSupportRenameColumn:       false,
		DontSupportForShareClause:     false,
		DontSupportNullAsDefaultValue: false,
		DontSupportRenameColumnUnique: false,
		DontSupportDropConstraint:     false,
	}
	connection, err := gorm.Open(mysql.New(configuration))

	require.NoError(t, err)

	return connection, mock
}

// MockExec mock a delete, insert or update.
func MockExec(mock sqlmock.Sqlmock, options Exec) {
	if !options.Raw {
		options.Query = regexp.QuoteMeta(options.Query)
	}

	if !options.Direct {
		mock.ExpectBegin()
	}

	if options.QueryArgs == nil {
		options.QueryArgs = []driver.Value{}
	}

	expect := mock.ExpectExec(options.Query).
		WithArgs(options.QueryArgs...)

	if options.Error != nil {
		expect.WillReturnError(options.Error)

		if !options.Direct {
			mock.ExpectRollback()
		}

		return
	}

	expect.WillReturnResult(sqlmock.NewResult(1, 1))

	if !options.Direct {
		mock.ExpectCommit()
	}
}

// MockSelect mock a select.
func MockSelect(mock sqlmock.Sqlmock, options Select) {
	if !options.Raw {
		options.Query = regexp.QuoteMeta(options.Query)
	}

	if options.QueryArgs == nil {
		options.QueryArgs = []driver.Value{}
	}

	expect := mock.ExpectQuery(options.Query).
		WithArgs(options.QueryArgs...)

	if options.Error != nil {
		expect.WillReturnError(options.Error)

		return
	}

	if options.Rows == nil {
		options.Rows = mock.NewRows(nil)
	}

	expect.WillReturnRows(options.Rows)
}
