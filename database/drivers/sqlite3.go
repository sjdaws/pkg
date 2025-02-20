package drivers

import (
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"

	"github.com/sjdaws/pkg/errors"
)

// SQLite3 options available when connecting to this driver.
type SQLite3 struct {
	Filename string
}

// GetDialector gorm dialector for this driver.
func (d SQLite3) GetDialector() (gorm.Dialector, error) {
	// refer https://github.com/glebarez/sqlite#usage for details
	if d.Filename == "" {
		return nil, errors.New(ErrInvalidDatabase)
	}

	return sqlite.Dialector{
		DriverName: "",
		DSN:        d.Filename + "?_pragma=foreign_keys(1)",
		Conn:       nil,
	}, nil
}
