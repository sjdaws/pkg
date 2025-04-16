package database

import (
	"database/sql/driver"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm/clause"

	"github.com/sjdaws/pkg/errors"
	"github.com/sjdaws/pkg/testing/database/modelmock"
	"github.com/sjdaws/pkg/testing/database/ormmock"
)

const (
	deleteQuery  = "UPDATE `model_mocks` SET `deleted_at`=? WHERE `model_mocks`.`id` = ? AND `model_mocks`.`deleted_at` IS NULL"
	insertQuery  = "INSERT INTO `model_mocks` (`deleted_at`,`test`) VALUES (?,?)"
	restoreQuery = "UPDATE `model_mocks` SET `deleted_at`=? WHERE `id` = ?"
	selectQuery  = "SELECT * FROM `model_mocks` WHERE (`model_mocks`.`id` = ? AND `model_mocks`.`test` = ?) AND `model_mocks`.`deleted_at` IS NULL"
	updateQuery  = "UPDATE `model_mocks` SET `deleted_at`=?,`test`=? WHERE `model_mocks`.`deleted_at` IS NULL AND `id` = ?"
)

func TestRepository_BypassDelete(t *testing.T) {
	t.Parallel()

	connection, _ := ormmock.New(t)
	instance := Repository[modelmock.ModelMock](&Database{orm: connection})

	result := instance.BypassDelete()

	assert.NotEqual(t, instance, result)

	actual, ok := result.(repository[modelmock.ModelMock])

	require.True(t, ok)
	assert.True(t, actual.unscoped)
}

func TestRepository_Create(t *testing.T) {
	t.Parallel()

	connection, mock := ormmock.New(t)
	instance := Repository[modelmock.ModelMock](&Database{orm: connection})

	model := &modelmock.ModelMock{Test: true}

	ormmock.MockExec(
		mock,
		ormmock.Exec{
			Query:     insertQuery,
			QueryArgs: []driver.Value{nil, true},
		},
	)

	err := instance.Create(model)
	require.NoError(t, err)

	assert.Equal(t, 1, model.ID)
}

func TestRepository_Create_ExecError(t *testing.T) {
	t.Parallel()

	connection, mock := ormmock.New(t)
	instance := Repository[modelmock.ModelMock](&Database{orm: connection})

	model := &modelmock.ModelMock{Test: true}

	ormmock.MockExec(
		mock,
		ormmock.Exec{
			Error:     errors.New("test"),
			Query:     insertQuery,
			QueryArgs: []driver.Value{nil, true},
		},
	)

	err := instance.Create(model)
	require.Error(t, err)

	require.EqualError(t, err, "unable to create record: test")
}

func TestRepository_Delete(t *testing.T) {
	t.Parallel()

	connection, mock := ormmock.New(t)
	instance := Repository[modelmock.ModelMock](&Database{orm: connection})

	model := &modelmock.ModelMock{Test: true}

	ormmock.MockExec(
		mock,
		ormmock.Exec{
			Query:     insertQuery,
			QueryArgs: []driver.Value{nil, true},
		},
	)
	ormmock.MockExec(
		mock,
		ormmock.Exec{
			Query:     deleteQuery,
			QueryArgs: []driver.Value{ormmock.TimeArg{}, 1},
		},
	)

	err := instance.Create(model)
	require.NoError(t, err)

	err = instance.Delete(model)
	require.NoError(t, err)

	assert.NotNil(t, model.DeletedAt)
}

func TestRepository_Delete_ExecError(t *testing.T) {
	t.Parallel()

	connection, mock := ormmock.New(t)
	instance := Repository[modelmock.ModelMock](&Database{orm: connection})

	model := &modelmock.ModelMock{Test: true}

	ormmock.MockExec(
		mock,
		ormmock.Exec{
			Query:     insertQuery,
			QueryArgs: []driver.Value{nil, true},
		},
	)
	ormmock.MockExec(
		mock,
		ormmock.Exec{
			Error:     errors.New("test"),
			Query:     deleteQuery,
			QueryArgs: []driver.Value{ormmock.TimeArg{}, 1},
		},
	)

	err := instance.Create(model)
	require.NoError(t, err)

	err = instance.Delete(model)
	require.Error(t, err)

	require.EqualError(t, err, "unable to delete record: test")
}

func TestRepository_Get(t *testing.T) {
	t.Parallel()

	connection, mock := ormmock.New(t)
	instance := Repository[modelmock.ModelMock](&Database{orm: connection})

	ormmock.MockSelect(
		mock,
		ormmock.Select{
			Query:     selectQuery,
			QueryArgs: []driver.Value{1, true},
			Rows:      mock.NewRows([]string{"id"}).AddRow(1),
		},
	)

	result, err := instance.Get(modelmock.ModelMock{ID: 1, Test: true})
	require.NoError(t, err)

	assert.Len(t, result, 1)
}

func TestRepository_Get_NoResultsError(t *testing.T) {
	t.Parallel()

	connection, mock := ormmock.New(t)
	instance := Repository[modelmock.ModelMock](&Database{orm: connection})

	ormmock.MockSelect(
		mock,
		ormmock.Select{
			Query:     selectQuery,
			QueryArgs: []driver.Value{1, true},
		},
	)

	model, err := instance.Get(modelmock.ModelMock{ID: 1, Test: true})
	require.Error(t, err)

	require.EqualError(t, err, "no results returned for query")
	require.ErrorIs(t, err, ErrNoResults)
	assert.Nil(t, model)
}

func TestRepository_Get_QueryError(t *testing.T) {
	t.Parallel()

	connection, mock := ormmock.New(t)
	instance := Repository[modelmock.ModelMock](&Database{orm: connection})

	ormmock.MockSelect(
		mock,
		ormmock.Select{
			Error:     errors.New("test"),
			Query:     selectQuery,
			QueryArgs: []driver.Value{1, true},
		},
	)

	result, err := instance.Get(modelmock.ModelMock{ID: 1, Test: true})
	require.Error(t, err)

	require.EqualError(t, err, "unable to fetch records: test")
	assert.Nil(t, result)
}

func TestRepository_One(t *testing.T) {
	t.Parallel()

	connection, mock := ormmock.New(t)
	instance := Repository[modelmock.ModelMock](&Database{orm: connection})

	ormmock.MockSelect(
		mock,
		ormmock.Select{
			Query:     selectQuery,
			QueryArgs: []driver.Value{1, true, 1},
			Rows:      mock.NewRows([]string{"id"}).AddRow(1),
		},
	)

	model, err := instance.One(modelmock.ModelMock{ID: 1, Test: true})
	require.NoError(t, err)

	assert.IsType(t, &modelmock.ModelMock{}, model)
}

func TestRepository_One_NoResultsError(t *testing.T) {
	t.Parallel()

	connection, mock := ormmock.New(t)
	instance := Repository[modelmock.ModelMock](&Database{orm: connection})

	ormmock.MockSelect(
		mock,
		ormmock.Select{
			Query:     selectQuery,
			QueryArgs: []driver.Value{1, true, 1},
		},
	)

	model, err := instance.One(modelmock.ModelMock{ID: 1, Test: true})
	require.Error(t, err)

	require.EqualError(t, err, "no results returned for query")
	require.ErrorIs(t, err, ErrNoResults)
	assert.Nil(t, model)
}

func TestRepository_One_QueryError(t *testing.T) {
	t.Parallel()

	connection, mock := ormmock.New(t)
	instance := Repository[modelmock.ModelMock](&Database{orm: connection})

	ormmock.MockSelect(
		mock,
		ormmock.Select{
			Error:     errors.New("test"),
			Query:     selectQuery,
			QueryArgs: []driver.Value{1, true, 1},
		},
	)

	model, err := instance.One(modelmock.ModelMock{ID: 1, Test: true})
	require.Error(t, err)

	require.EqualError(t, err, "unable to fetch record: test")
	assert.Nil(t, model)
}

func TestRepository_OrderBy(t *testing.T) {
	t.Parallel()

	connection, _ := ormmock.New(t)
	instance := Repository[modelmock.ModelMock](&Database{orm: connection})

	result := instance.OrderBy(Order{Column: "id", Descending: true})

	assert.NotEqual(t, instance, result)

	actual, ok := result.(repository[modelmock.ModelMock])

	require.True(t, ok)
	assert.Equal(t, []Order{{Column: "id", Descending: true}}, actual.order)
}

func TestRepository_PartOf(t *testing.T) {
	t.Parallel()

	connection, _ := ormmock.New(t)
	instance := Repository[modelmock.ModelMock](&Database{orm: connection})

	transaction, _ := ormmock.New(t)

	result := instance.PartOf(transaction)

	assert.NotEqual(t, instance, result)

	actual, ok := result.(repository[modelmock.ModelMock])

	require.True(t, ok)
	assert.Equal(t, actual.connection, transaction)
}

func TestRepository_Restore(t *testing.T) {
	t.Parallel()

	connection, mock := ormmock.New(t)
	instance := Repository[modelmock.ModelMock](&Database{orm: connection})

	model := &modelmock.ModelMock{Test: true}

	ormmock.MockExec(
		mock,
		ormmock.Exec{
			Query:     insertQuery,
			QueryArgs: []driver.Value{nil, true},
		},
	)
	ormmock.MockExec(
		mock,
		ormmock.Exec{
			Query:     deleteQuery,
			QueryArgs: []driver.Value{ormmock.TimeArg{}, 1},
		},
	)
	ormmock.MockExec(
		mock,
		ormmock.Exec{
			Query:     restoreQuery,
			QueryArgs: []driver.Value{nil, 1},
		},
	)

	err := instance.Create(model)
	require.NoError(t, err)

	err = instance.Delete(model)
	require.NoError(t, err)

	err = instance.Restore(model)
	require.NoError(t, err)

	assert.Nil(t, model.DeletedAt)
}

func TestRepository_Restore_ExecError(t *testing.T) {
	t.Parallel()

	connection, mock := ormmock.New(t)
	instance := Repository[modelmock.ModelMock](&Database{orm: connection})

	model := &modelmock.ModelMock{Test: true}

	ormmock.MockExec(
		mock,
		ormmock.Exec{
			Query:     insertQuery,
			QueryArgs: []driver.Value{nil, true},
		},
	)
	ormmock.MockExec(
		mock,
		ormmock.Exec{
			Query:     deleteQuery,
			QueryArgs: []driver.Value{ormmock.TimeArg{}, 1},
		},
	)
	ormmock.MockExec(
		mock,
		ormmock.Exec{
			Error:     errors.New("test"),
			Query:     restoreQuery,
			QueryArgs: []driver.Value{nil, 1},
		},
	)

	err := instance.Create(model)
	require.NoError(t, err)

	err = instance.Delete(model)
	require.NoError(t, err)

	err = instance.Restore(model)
	require.Error(t, err)

	require.EqualError(t, err, "unable to restore record: test")
}

func TestRepository_Then(t *testing.T) {
	t.Parallel()

	connection, _ := ormmock.New(t)
	instance := Repository[modelmock.ModelMock](&Database{orm: connection})

	result := instance.Then("Relation")

	assert.NotEqual(t, instance, result)

	actual, ok := result.(repository[modelmock.ModelMock])

	require.True(t, ok)
	assert.Equal(t, []relation{{join: false, key: "Relation"}}, actual.relations)
}

func TestRepository_Update(t *testing.T) {
	t.Parallel()

	connection, mock := ormmock.New(t)
	instance := Repository[modelmock.ModelMock](&Database{orm: connection})

	model := &modelmock.ModelMock{Test: true}

	ormmock.MockExec(
		mock,
		ormmock.Exec{
			Query:     insertQuery,
			QueryArgs: []driver.Value{nil, true},
		},
	)

	err := instance.Update(model)
	require.NoError(t, err)

	assert.Equal(t, 1, model.ID)
}

func TestRepository_Update_ExecError(t *testing.T) {
	t.Parallel()

	connection, mock := ormmock.New(t)
	instance := Repository[modelmock.ModelMock](&Database{orm: connection})

	model := &modelmock.ModelMock{Test: true}

	// model must be updated twice: initial update will create/hydrate the model, second update will perform update
	ormmock.MockExec(
		mock,
		ormmock.Exec{
			Query:     insertQuery,
			QueryArgs: []driver.Value{nil, true},
		},
	)

	err := instance.Update(model)
	require.NoError(t, err)

	assert.Equal(t, 1, model.ID)

	model.Test = false

	ormmock.MockExec(
		mock,
		ormmock.Exec{
			Error:     errors.New("test"),
			Query:     updateQuery,
			QueryArgs: []driver.Value{nil, false, model.ID},
		},
	)

	err = instance.Update(model)
	require.Error(t, err)

	require.EqualError(t, err, "unable to update record: test")
}

func TestRepository_With(t *testing.T) {
	t.Parallel()

	connection, _ := ormmock.New(t)
	instance := Repository[modelmock.ModelMock](&Database{orm: connection})

	result := instance.With("Relation")

	assert.NotEqual(t, instance, result)

	actual, ok := result.(repository[modelmock.ModelMock])

	require.True(t, ok)
	assert.Equal(t, []relation{{join: true, key: "Relation"}}, actual.relations)
}

func TestRepository_addMeta(t *testing.T) {
	t.Parallel()

	connection, err := Connect(
		false,
		"sqlite",
		"",
		":memory:",
		"",
		0,
		"",
		"",
		"",
	)
	require.NoError(t, err)

	actual := repository[modelmock.ModelMock]{
		connection: connection.orm,
	}
	transaction := actual.addMeta(connection.orm)

	// test everything is empty
	assert.False(t, transaction.Statement.Unscoped)
	assert.Equal(t, map[string][]any(nil), transaction.Statement.Preloads)

	// add bypass delete
	actual.unscoped = true
	transaction = actual.addMeta(connection.orm)

	assert.True(t, transaction.Statement.Unscoped)

	// add order by
	actual.order = []Order{{Column: "id"}}
	transaction = actual.addMeta(connection.orm)

	assert.Equal(
		t,
		map[string]clause.Clause{
			"ORDER BY": {Name: "ORDER BY", Expression: clause.OrderBy{Columns: []clause.OrderByColumn{{Column: clause.Column{Name: "id asc", Raw: true}}}}},
		},
		transaction.Statement.Clauses,
	)

	// add eager load
	actual.relations = []relation{{join: false, key: "Relation"}}
	transaction = actual.addMeta(connection.orm)

	assert.Len(t, transaction.Statement.Preloads, 1)
	assert.Equal(t, map[string][]interface{}{"Relation": nil}, transaction.Statement.Preloads)

	// add join
	actual.relations = []relation{{join: true, key: "Relation"}}
	transaction = actual.addMeta(connection.orm)

	assert.Len(t, transaction.Statement.Joins, 1)
	assert.Equal(t, "Relation", transaction.Statement.Joins[0].Name)
}

func TestRepository_query(t *testing.T) {
	t.Parallel()

	testcases := map[string]struct {
		expectedParameters []driver.Value
		expectedQuery      string
		where              []any
	}{
		"both": {
			expectedParameters: []driver.Value{3, true, 4},
			expectedQuery:      "SELECT * FROM `model_mocks` WHERE `model_mocks`.`id.` = ? AND (`model_mocks`.`test` = ? OR `model_mocks`.`id` = ?) AND `model_mocks`.`deleted_at` IS NULL",
			where:              []any{&modelmock.ModelMock{ID: 3}, Or{&modelmock.ModelMock{Test: true}, &modelmock.ModelMock{ID: 4}}},
		},
		"or": {
			expectedParameters: []driver.Value{true, 4},
			expectedQuery:      "SELECT * FROM `model_mocks` WHERE (`model_mocks`.`test` = ? OR `model_mocks`.`id` = ?) AND `model_mocks`.`deleted_at` IS NULL",
			where:              []any{Or{&modelmock.ModelMock{Test: true}, &modelmock.ModelMock{ID: 4}}},
		},
		"raw": {
			expectedParameters: []driver.Value{3, true},
			expectedQuery:      "SELECT * FROM `model_mocks` WHERE `model_mocks`.`test` = ? AND `model_mocks`.`deleted_at` IS NULL",
			where:              []any{Raw{Query: "test = ?", Parameters: []any{true}}},
		},
		"where": {
			expectedParameters: []driver.Value{true},
			expectedQuery:      "SELECT * FROM `model_mocks` WHERE `model_mocks`.`test` = ? AND `model_mocks`.`deleted_at` IS NULL",
			where:              []any{&modelmock.ModelMock{Test: true}},
		},
		"nothing": {
			expectedParameters: []driver.Value{},
			expectedQuery:      "SELECT * FROM `model_mocks` WHERE `model_mocks`.`deleted_at` IS NULL",
			where:              []any{nil},
		},
	}

	for name, testcase := range testcases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			connection, mock := ormmock.New(t)
			repo := repository[modelmock.ModelMock]{
				connection: connection,
			}

			ormmock.MockSelect(
				mock,
				ormmock.Select{
					Query:     testcase.expectedQuery,
					QueryArgs: testcase.expectedParameters,
					Rows:      mock.NewRows([]string{"id"}).AddRow("00000000-0000-0000-0000-000000000001"),
				},
			)

			transaction := repo.query(testcase.where...)
			_ = transaction.Find(&modelmock.ModelMock{})
		})
	}
}
