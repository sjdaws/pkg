package repositorymock_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	"github.com/sjdaws/pkg/database"
	"github.com/sjdaws/pkg/errors"
	"github.com/sjdaws/pkg/testing/database/modelmock"
	"github.com/sjdaws/pkg/testing/database/repositorymock"
)

func TestRepositoryMock(t *testing.T) {
	t.Parallel()

	repository := repositorymock.RepositoryMock[modelmock.ModelMock]{}

	assert.Implements(t, (*database.Persister[modelmock.ModelMock])(nil), &repository)
}

func TestRepositoryMock_BypassDelete(t *testing.T) {
	t.Parallel()

	repository := repositorymock.RepositoryMock[modelmock.ModelMock]{}

	result := repository.BypassDelete()

	assert.Equal(t, repository, result)
}

func TestRepositoryMock_Create(t *testing.T) {
	t.Parallel()

	model := &modelmock.ModelMock{ID: 1}
	repository := repositorymock.RepositoryMock[modelmock.ModelMock]{
		CreateMock: func(_ *modelmock.ModelMock) error {
			return errors.New("create")
		},
	}

	err := repository.Create(model)
	require.Error(t, err)

	require.EqualError(t, err, "create")
}

func TestRepositoryMock_Delete(t *testing.T) {
	t.Parallel()

	repository := repositorymock.RepositoryMock[modelmock.ModelMock]{
		DeleteMock: func(_ *modelmock.ModelMock, _ ...any) error {
			return errors.New("delete")
		},
	}

	err := repository.Delete(nil)
	require.Error(t, err)

	require.EqualError(t, err, "delete")
}

func TestRepositoryMock_Get(t *testing.T) {
	t.Parallel()

	model := modelmock.ModelMock{ID: 1}
	repository := repositorymock.RepositoryMock[modelmock.ModelMock]{
		GetMock: func(_ ...any) ([]modelmock.ModelMock, error) {
			return []modelmock.ModelMock{model}, errors.New("get")
		},
	}

	get, err := repository.Get(nil)
	require.Error(t, err)

	require.EqualError(t, err, "get")
	assert.Equal(t, []modelmock.ModelMock{model}, get)
}

func TestRepositoryMock_One(t *testing.T) {
	t.Parallel()

	model := &modelmock.ModelMock{ID: 1}
	repository := repositorymock.RepositoryMock[modelmock.ModelMock]{
		OneMock: func(_ ...any) (*modelmock.ModelMock, error) {
			return model, errors.New("one")
		},
	}

	one, err := repository.One(nil)
	require.Error(t, err)

	require.EqualError(t, err, "one")
	assert.Equal(t, model, one)
}

func TestRepositoryMock_Order(t *testing.T) {
	t.Parallel()

	repository := repositorymock.RepositoryMock[modelmock.ModelMock]{}

	result := repository.OrderBy(database.Order{Column: "test", Descending: true})

	assert.Equal(t, repository, result)
}

func TestRepositoryMock_PartOf(t *testing.T) {
	t.Parallel()

	repository := repositorymock.RepositoryMock[modelmock.ModelMock]{}

	result := repository.PartOf(&gorm.DB{})

	assert.Equal(t, repository, result)
}

func TestRepositoryMock_Restore(t *testing.T) {
	t.Parallel()

	repository := repositorymock.RepositoryMock[modelmock.ModelMock]{
		RestoreMock: func(_ *modelmock.ModelMock) error {
			return errors.New("restore")
		},
	}

	err := repository.Restore(nil)
	require.Error(t, err)

	require.EqualError(t, err, "restore")
}

func TestRepositoryMock_Then(t *testing.T) {
	t.Parallel()

	repository := repositorymock.RepositoryMock[modelmock.ModelMock]{}

	result := repository.Then("")

	assert.Equal(t, repository, result)
}

func TestRepositoryMock_Update(t *testing.T) {
	t.Parallel()

	model := &modelmock.ModelMock{ID: 1}
	repository := repositorymock.RepositoryMock[modelmock.ModelMock]{
		UpdateMock: func(_ *modelmock.ModelMock) error {
			return errors.New("update")
		},
	}

	err := repository.Update(model)
	require.Error(t, err)

	require.EqualError(t, err, "update")
}

func TestRepositoryMock_With(t *testing.T) {
	t.Parallel()

	repository := repositorymock.RepositoryMock[modelmock.ModelMock]{}

	result := repository.With("")

	assert.Equal(t, repository, result)
}
