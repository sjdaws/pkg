package drivers

import (
	"fmt"
	"strings"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/sjdaws/pkg/errors"
)

// PostgreSQL options available when connecting to this driver.
type PostgreSQL struct {
	Host            string
	Name            string
	Password        string
	Port            int
	SSLModeDisabled bool
	Username        string
}

// GetDialector gorm dialector for this driver.
func (d PostgreSQL) GetDialector() (gorm.Dialector, error) {
	// refer https://gorm.io/docs/connecting_to_the_database.html#PostgreSQL for details
	if d.Name == "" {
		return nil, errors.New(ErrInvalidDatabase)
	}

	dsn := []string{"TimeZone=UTC dbname=" + d.Name}

	if d.Host != "" {
		dsn = append(dsn, "host="+d.Host)
	}

	if d.Password != "" {
		dsn = append(dsn, "password="+d.Password)
	}

	if d.Port > 0 {
		dsn = append(dsn, fmt.Sprintf("port=%d", d.Port))
	}

	if d.SSLModeDisabled {
		dsn = append(dsn, "sslmode=disabled")
	}

	if d.Username != "" {
		dsn = append(dsn, "user="+d.Username)
	}

	return postgres.Dialector{
		Config: &postgres.Config{
			DriverName:           "",
			DSN:                  strings.Join(dsn, " "),
			WithoutQuotingCheck:  false,
			PreferSimpleProtocol: false,
			WithoutReturning:     false,
			Conn:                 nil,
		},
	}, nil
}
