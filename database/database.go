package database

import (
	"strings"

	"github.com/carlmjohnson/truthy"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/sjdaws/pkg/database/drivers"
	"github.com/sjdaws/pkg/errors"
)

// Connection instance.
type Connection struct {
	orm *gorm.DB
}

// Driver interface.
type Driver interface {
	GetDialector() (gorm.Dialector, error)
}

// Connect create a new database connection.
func Connect(
	debug bool,
	driver string,
	host string,
	name string,
	password string,
	port int,
	socket string,
	sslmode string,
	username string,
) (*Connection, error) {
	var options Driver

	switch strings.ToLower(driver) {
	case "mariadb", "mysql":
		options = drivers.MySQL{
			Host:     host,
			Name:     name,
			Password: password,
			Port:     port,
			Socket:   socket,
			Username: username,
		}
	case "postgres", "postgresql":
		options = drivers.PostgreSQL{
			Host:            host,
			Name:            name,
			Password:        password,
			Port:            port,
			SSLModeDisabled: strings.EqualFold(sslmode, "disabled"),
			Username:        username,
		}
	case "sqlite", "sqlite3":
		options = drivers.SQLite3{
			Filename: name,
		}
	case "sqlserver":
		options = drivers.SQLServer{
			Host:     host,
			Name:     name,
			Password: password,
			Port:     port,
			Username: username,
		}
	default:
		return nil, errors.New("unsupported database type requested: %s", driver)
	}

	dialector, err := options.GetDialector()
	if err != nil {
		return nil, errors.Wrap(err, "unable to create dialector")
	}

	orm, err := gorm.Open(dialector, createConfiguration(truthy.Cond(debug, logger.Info, logger.Warn)))
	if err != nil {
		return nil, errors.Wrap(err, "unable to open connection to database")
	}

	return &Connection{orm: orm}, nil
}

// Migrate run database migrations.
func (c *Connection) Migrate(model ...any) error {
	// Force InnoDB for MySQL-like DBs
	if c.orm.Dialector.Name() == "mysql" {
		c.orm.Set("gorm:table_options", "ENGINE=InnoDB")
	}

	err := c.orm.AutoMigrate(model...)
	if err != nil {
		return errors.Wrap(err, "unable to invoke database migrations")
	}

	return nil
}

// Transaction start a transaction and return transactional connection.
func (c *Connection) Transaction() *gorm.DB {
	return c.orm.Begin()
}

// createConfiguration creates a configuration for a dialector.
func createConfiguration(logMode logger.LogLevel) *gorm.Config {
	return &gorm.Config{
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
		Logger:                                   logger.Default.LogMode(logMode),
		NamingStrategy:                           nil,
		NowFunc:                                  nil,
		Plugins:                                  nil,
		PrepareStmt:                              false,
		PropagateUnscoped:                        false,
		QueryFields:                              false,
		SkipDefaultTransaction:                   false,
		TranslateError:                           false,
	}
}
