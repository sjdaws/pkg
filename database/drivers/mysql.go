package drivers

import (
	"fmt"

	base "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/sjdaws/pkg/errors"
)

// MySQL options available when connecting to this driver.
type MySQL struct {
	Host     string
	Name     string
	Password string
	Port     int
	Socket   string
	Username string
}

// GetDialector gorm dialector for this driver.
func (d MySQL) GetDialector() (gorm.Dialector, error) {
	// refer https://github.com/go-sql-driver/mysql#dsn-data-source-name for details
	if d.Name == "" {
		return nil, errors.New(ErrInvalidDatabase)
	}

	host := ""
	if d.Socket != "" {
		host = fmt.Sprintf("unix(%s)", d.Socket)
	} else if d.Host != "" {
		host = fmt.Sprintf("tcp(%s)", d.Host)
		if d.Port > 0 {
			host = fmt.Sprintf("tcp(%s:%d)", d.Host, d.Port)
		}
	}

	credentials := d.Username
	if d.Password != "" {
		credentials = fmt.Sprintf("%s:%s", d.Username, d.Password)
	}

	dsn := fmt.Sprintf(
		"%s@%s/%s?charset=utf8mb4&parseTime=True&loc=UTC",
		credentials,
		host,
		d.Name,
	)

	dsnConfig, _ := base.ParseDSN(dsn)

	return mysql.Dialector{
		Config: &mysql.Config{
			DriverName:                    "",
			ServerVersion:                 "",
			DSN:                           dsn,
			DSNConfig:                     dsnConfig,
			Conn:                          nil,
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
		},
	}, nil
}
