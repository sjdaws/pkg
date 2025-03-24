package database

import (
	"database/sql/driver"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/sjdaws/pkg/testing/database/modelmock"
	"github.com/sjdaws/pkg/testing/database/ormmock"
)

func TestMigrate(t *testing.T) {
	t.Parallel()

	orm, mock := ormmock.New(t)
	connection := &Database{orm: orm}

	query := "SELECT SCHEMA_NAME from Information_schema.SCHEMATA where SCHEMA_NAME LIKE ? ORDER BY SCHEMA_NAME=? DESC,SCHEMA_NAME limit 1"
	queryArgs := []driver.Value{"%", ""}

	model := modelmock.ModelMock{}

	ormmock.MockSelect(mock, ormmock.Select{Query: query, QueryArgs: queryArgs})
	ormmock.MockExec(
		mock,
		ormmock.Exec{Direct: true, Query: "CREATE TABLE `" + model.TableName() + "` .*", Raw: true},
	)
	ormmock.MockSelect(mock, ormmock.Select{Query: query, QueryArgs: queryArgs})

	err := connection.Migrate(model)
	require.NoError(t, err)
}
