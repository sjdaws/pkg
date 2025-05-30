package database

import (
	"github.com/carlmjohnson/truthy"
	"gorm.io/gorm"

	"github.com/sjdaws/pkg/errors"
)

// Persister interface.
type Persister[m Model] interface {
	BypassDelete() Persister[m]
	Create(model *m) error
	Delete(model *m, where ...any) error
	Get(where ...any) ([]m, error)
	One(where ...any) (*m, error)
	OrderBy(order ...Order) Persister[m]
	PartOf(connection *gorm.DB) Persister[m]
	Restore(model *m) error
	Then(relationship string, where ...any) Persister[m]
	Update(model *m) error
	With(relationship string, where ...any) Persister[m]
}

// Or type for holding where queries which should be OR.
type Or []any

// Order parameter for Persister.OrderBy.
type Order struct {
	Column     string
	Descending bool
}

// Raw type for holding a raw query with optional parameters.
type Raw struct {
	Parameters []any
	Query      string
}

// repository base repository which all repositories extend.
type repository[m Model] struct {
	connection *gorm.DB
	model      m
	models     []m
	order      []Order
	relations  []relation
	unscoped   bool
}

// relation to fetch with the initial request.
type relation struct {
	join  bool
	key   string
	where []any
}

// ErrNoResults error to return when there are no results returned from a query.
var ErrNoResults = errors.New("no results returned for query")

// Repository create a repository for a model.
func Repository[m Model](connection Connection) Persister[m] {
	var model m

	instance := repository[m]{
		connection: connection.ORM(),
		model:      model,
		models:     make([]m, 0),
		order:      make([]Order, 0),
		relations:  make([]relation, 0),
		unscoped:   false,
	}

	return instance
}

// BypassDelete return deleted records.
func (r repository[m]) BypassDelete() Persister[m] {
	transaction := r
	transaction.unscoped = true

	return transaction
}

// Create a new record from a model.
func (r repository[m]) Create(model *m) error {
	result := r.connection.Create(model)
	if result.Error != nil {
		return errors.Wrap(result.Error, "unable to create record")
	}

	return nil
}

// Delete a record.
func (r repository[m]) Delete(model *m, where ...any) error {
	result := r.connection.Delete(model, where...)
	if result.Error != nil {
		return errors.Wrap(result.Error, "unable to delete record")
	}

	return nil
}

// Get record(s) from a query.
func (r repository[m]) Get(where ...any) ([]m, error) {
	result := r.query(where...).Find(&r.models)
	if result.Error != nil {
		return nil, errors.Wrap(result.Error, "unable to fetch records")
	}

	if len(r.models) == 0 {
		return nil, ErrNoResults
	}

	return r.models, nil
}

// One fetches a single record from a query.
func (r repository[m]) One(where ...any) (*m, error) {
	result := r.query(where...).First(&r.model)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ErrNoResults
		}

		return nil, errors.Wrap(result.Error, "unable to fetch record")
	}

	return &r.model, nil
}

// OrderBy order results from a query.
func (r repository[m]) OrderBy(order ...Order) Persister[m] {
	transaction := r
	transaction.order = r.order
	transaction.order = append(transaction.order, order...)

	return transaction
}

// PartOf define a different connection for transactions.
func (r repository[m]) PartOf(connection *gorm.DB) Persister[m] {
	transaction := r
	transaction.connection = connection

	return transaction
}

// Restore a deleted record.
func (r repository[m]) Restore(model *m) error {
	result := r.connection.Unscoped().Model(model).Update("deleted_at", nil)
	if result.Error != nil {
		return errors.Wrap(result.Error, "unable to restore record")
	}

	return nil
}

// Then eager load relationship after initial query is complete.
func (r repository[m]) Then(relationship string, where ...any) Persister[m] {
	transaction := r
	transaction.relations = r.relations
	transaction.relations = append(transaction.relations, relation{join: false, key: relationship, where: where})

	return transaction
}

// Update a record from a model.
func (r repository[m]) Update(model *m) error {
	result := r.connection.Save(model)
	if result.Error != nil {
		return errors.Wrap(result.Error, "unable to update record")
	}

	return nil
}

// With get a relationship with query, otherwise return nothing.
func (r repository[m]) With(relationship string, where ...any) Persister[m] {
	transaction := r
	transaction.relations = r.relations
	transaction.relations = append(transaction.relations, relation{join: true, key: relationship, where: where})

	return transaction
}

// addMeta eager load requested relationships, process order.
func (r repository[m]) addMeta(transaction *gorm.DB) *gorm.DB {
	if r.unscoped {
		transaction = transaction.Unscoped()
	}

	for _, by := range r.order {
		transaction = transaction.Order(truthy.Cond(by.Descending, by.Column+" desc", by.Column+" asc"))
	}

	for _, relationship := range r.relations {
		// Use inner join for hasone relationships, this will cause no records to be returned if join is empty
		if relationship.join {
			transaction = transaction.InnerJoins(relationship.key, relationship.where...)

			continue
		}

		// Preload hasmany relationships, this will do a second select for the relationship
		transaction = transaction.Preload(relationship.key, relationship.where...)
	}

	return transaction
}

// query starts a query using map expectedParameters.
func (r repository[m]) query(where ...any) *gorm.DB {
	query := r.connection

	for _, condition := range where {
		switch state := condition.(type) {
		case Or:
			subQuery := r.connection
			for _, orCondition := range state {
				subQuery = subQuery.Or(orCondition)
			}

			query = query.Where(subQuery)

		case Raw:
			query = query.Where(state.Query, state.Parameters...)

		default:
			query = query.Where(condition)
		}
	}

	return r.addMeta(query)
}
