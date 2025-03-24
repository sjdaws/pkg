package ormmock_test

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	"github.com/sjdaws/pkg/errors"
	"github.com/sjdaws/pkg/testing/database/ormmock"
)

var (
	errExec  = errors.New("exec error")
	errQuery = errors.New("query error")
)

func TestNew(t *testing.T) {
	t.Parallel()

	connection, mock := ormmock.New(t)
	assert.IsType(t, &gorm.DB{}, connection)
	assert.Implements(t, (*sqlmock.Sqlmock)(nil), mock)
}

func TestMockExec(t *testing.T) {
	t.Parallel()

	_, mock := ormmock.New(t)

	ormmock.MockExec(
		mock,
		ormmock.Exec{Query: "query"},
	)

	expected := "there is a remaining expectation which was not matched: ExpectedBegin => expecting database transaction Begin"
	require.EqualError(t, mock.ExpectationsWereMet(), expected)
}

func TestMockExec_Error(t *testing.T) {
	t.Parallel()

	_, mock := ormmock.New(t)

	ormmock.MockExec(
		mock,
		ormmock.Exec{Error: errExec, Query: "query"},
	)

	expected := "there is a remaining expectation which was not matched: ExpectedBegin => expecting database transaction Begin"
	require.EqualError(t, mock.ExpectationsWereMet(), expected)
}

func TestMockSelect(t *testing.T) {
	t.Parallel()

	_, mock := ormmock.New(t)

	ormmock.MockSelect(mock, ormmock.Select{Query: "query"})

	expected := "there is a remaining expectation which was not matched: ExpectedQuery => expecting Query, QueryContext or QueryRow which:\n" +
		"  - matches sql: 'query'\n" +
		"  - is without arguments\n" +
		"  - with empty rows"
	require.EqualError(t, mock.ExpectationsWereMet(), expected)
}

func TestMockSelect_Error(t *testing.T) {
	t.Parallel()

	_, mock := ormmock.New(t)

	ormmock.MockSelect(mock, ormmock.Select{Error: errQuery, Query: "query"})

	expected := "there is a remaining expectation which was not matched: ExpectedQuery => expecting Query, QueryContext or QueryRow which:\n" +
		"  - matches sql: 'query'\n" +
		"  - is without arguments\n" +
		"  - should return error: query error"
	require.EqualError(t, mock.ExpectationsWereMet(), expected)
}
