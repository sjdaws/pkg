package repositorymock

import (
	"gorm.io/gorm"

	"github.com/sjdaws/pkg/database"
)

// RepositoryMock fakes a repository.
type RepositoryMock[m database.Model] struct {
	CreateMock  func(model *m) error
	DeleteMock  func(model *m) error
	GetMock     func(where ...any) ([]m, error)
	OneMock     func(where ...any) (*m, error)
	RestoreMock func(model *m) error
	UpdateMock  func(model *m) error
}

// Also do nothing.
func (r RepositoryMock[m]) Also(_ string, _ ...any) database.Persister[m] {
	return r
}

// And do nothing.
func (r RepositoryMock[m]) And(_ string, _ ...any) database.Persister[m] {
	return r
}

// BypassDelete do nothing.
func (r RepositoryMock[m]) BypassDelete() database.Persister[m] {
	return r
}

// Create run CreateMock() function.
func (r RepositoryMock[m]) Create(model *m) error {
	return r.CreateMock(model)
}

// Delete run DeleteMock() function.
func (r RepositoryMock[m]) Delete(model *m) error {
	return r.DeleteMock(model)
}

// Get run GetMock() function.
func (r RepositoryMock[m]) Get(where ...any) ([]m, error) {
	return r.GetMock(where...)
}

// One run OneMock() function.
func (r RepositoryMock[m]) One(where ...any) (*m, error) {
	return r.OneMock(where...)
}

// OrderBy run OrderByMock() function.
func (r RepositoryMock[m]) OrderBy(_ ...database.Order) database.Persister[m] {
	return r
}

// PartOf do nothing.
func (r RepositoryMock[m]) PartOf(_ *gorm.DB) database.Persister[m] {
	return r
}

// Restore run RestoreMock() function.
func (r RepositoryMock[m]) Restore(model *m) error {
	return r.RestoreMock(model)
}

// Update run UpdateMock() function.
func (r RepositoryMock[m]) Update(model *m) error {
	return r.UpdateMock(model)
}
