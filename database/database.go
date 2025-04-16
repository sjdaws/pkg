package database

import (
	"strings"

	"github.com/carlmjohnson/truthy"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/sjdaws/pkg/database/drivers"
	"github.com/sjdaws/pkg/errors"
)

// Connection interface.
type Connection interface {
	Migrate(model ...any) error
	ORM() *gorm.DB
	Transaction() *gorm.DB
}

// Database instance.
type Database struct {
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
) (*Database, error) {
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

	return &Database{orm: orm}, nil
}

// Migrate run database migrations.
func (d *Database) Migrate(model ...any) error {
	// Force InnoDB for MySQL-like DBs
	if d.orm.Name() == "mysql" {
		d.orm.Set("gorm:table_options", "ENGINE=InnoDB")
	}

	err := d.orm.AutoMigrate(model...)
	if err != nil {
		return errors.Wrap(err, "unable to invoke database migrations")
	}

	return nil
}

// ORM return the underlying ORM.
func (d *Database) ORM() *gorm.DB {
	return d.orm
}

// Transaction start a transaction and return transactional connection.
func (d *Database) Transaction() *gorm.DB {
	return d.orm.Begin()
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
