package database

import (
	"strings"

	"github.com/carlmjohnson/truthy"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/sjdaws/pkg/database/drivers"
	"github.com/sjdaws/pkg/errors"
)

// Database interface.
type Database interface {
	GetDialector() (gorm.Dialector, error)
}

// New create a new database connection.
func New(
	debug bool,
	driver string,
	host string,
	name string,
	password string,
	port int,
	socket string,
	sslmode string,
	username string,
) (*gorm.DB, error) {
	var options Database

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

	connection, err := gorm.Open(dialector, createConfiguration(truthy.Cond(debug, logger.Info, logger.Warn)))
	if err != nil {
		return nil, errors.Wrap(err, "unable to open connection to database")
	}

	return connection, nil
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
