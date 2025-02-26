package database

import (
	"database/sql/driver"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/sjdaws/pkg/testing/database/connectionmock"
	"github.com/sjdaws/pkg/testing/database/modelmock"
)

func TestMigrate(t *testing.T) {
	t.Parallel()

	orm, mock := connectionmock.New(t)
	connection := &Connection{orm: orm}

	query := "SELECT SCHEMA_NAME from Information_schema.SCHEMATA where SCHEMA_NAME LIKE ? ORDER BY SCHEMA_NAME=? DESC,SCHEMA_NAME limit 1"
	queryArgs := []driver.Value{"%", ""}

	model := modelmock.ModelMock{}

	connectionmock.MockSelect(mock, connectionmock.Select{Query: query, QueryArgs: queryArgs})
	connectionmock.MockExec(
		mock,
		connectionmock.Exec{Direct: true, Query: "CREATE TABLE `" + model.TableName() + "` .*", Raw: true},
	)
	connectionmock.MockSelect(mock, connectionmock.Select{Query: query, QueryArgs: queryArgs})

	err := connection.Migrate(model)
	require.NoError(t, err)
}
