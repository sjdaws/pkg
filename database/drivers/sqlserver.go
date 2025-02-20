package drivers

import (
	"fmt"

	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"

	"github.com/sjdaws/pkg/errors"
)

// SQLServer options available when connecting to this driver.
type SQLServer struct {
	Host     string
	Name     string
	Password string
	Port     int
	Username string
}

// GetDialector gorm dialector for this driver.
func (d SQLServer) GetDialector() (gorm.Dialector, error) {
	// refer https://gorm.io/docs/connecting_to_the_database.html#SQL-Server for details
	if d.Name == "" {
		return nil, errors.New(ErrInvalidDatabase)
	}

	host := d.Host
	if d.Port > 0 {
		host = fmt.Sprintf("%s:%d", d.Host, d.Port)
	}

	credentials := d.Username
	if d.Password != "" {
		credentials = fmt.Sprintf("%s:%s", d.Username, d.Password)
	}

	dsn := fmt.Sprintf(
		"sqlserver://%s@%s?database=%s",
		credentials,
		host,
		d.Name,
	)

	return sqlserver.Dialector{
		Config: &sqlserver.Config{
			DriverName:        "",
			DSN:               dsn,
			DefaultStringSize: 0,
			Conn:              nil,
		},
	}, nil
}
