package progorm

import (
	"github.com/owenlilly/progorm/connection"
	"gorm.io/gorm"
)

type BaseTypedRepository[T any] struct {
	baseRepo BaseRepository
}

func NewBaseTypedRepository[T any](connMan connection.Manager) BaseTypedRepository[T] {
	return BaseTypedRepository[T]{
		baseRepo: NewBaseRepository(connMan),
	}
}

// InsertRecord model to insert must be a pointer/reference type
func (r BaseTypedRepository[T]) InsertRecord(model *T) error {
	return r.baseRepo.InsertRecord(model)
}

// First model to insert must be a pointer/reference type
func (r BaseTypedRepository[T]) First(model T) (*T, error) {
	var result T

	dbResult := r.DB().First(&result, model)
	if dbResult.Error != nil {
		return nil, dbResult.Error
	}

	return &result, nil
}

// FindRecords finding records with pagination
func (r BaseTypedRepository[T]) FindRecords(page, perPage uint, query *gorm.DB) (PageTyped[T], error) {
	var results []T

	p, err := r.baseRepo.FindRecords(page, perPage, query, &results)
	typedPage := PageTyped[T]{
		Total:   p.Total,
		Page:    p.Page,
		PerPage: p.PerPage,
		Pages:   p.Pages,
	}
	if err != nil {
		return typedPage, err
	}

	typedPage.Results = results

	return typedPage, nil
}

// Count counts total number of records for the given query
func (r BaseTypedRepository[T]) Count(model T, query any, args ...any) (count int64, err error) {
	return r.baseRepo.Count(model, query, args...)
}

// IsRecordNotFoundError returns true if the given error is a gorm.ErrRecordNotFound error, and false otherwise
func (r BaseTypedRepository[T]) IsRecordNotFoundError(err error) bool {
	if err == nil {
		return false
	}

	return err == gorm.ErrRecordNotFound
}

// AutoMigrate create a table for the given model or return an error
func (r BaseTypedRepository[T]) AutoMigrate(model T) error {
	return r.baseRepo.AutoMigrate(&model)
}

// AutoMigrateOrWarn creates a table for the given model or print a warning message if there's an error
func (r BaseTypedRepository[T]) AutoMigrateOrWarn(model T) {
	r.baseRepo.AutoMigrateOrWarn(&model)
}

// ConnectionManager returns the current ConnectionManager
func (r BaseTypedRepository[T]) ConnectionManager() connection.Manager {
	return r.baseRepo.ConnectionManager()
}

// DB return a struct contain gorm database connection information
func (r BaseTypedRepository[T]) DB() *gorm.DB {
	return r.baseRepo.DB()
}

// region: Transaction section

// WithTx start a new database transaction
func (r BaseTypedRepository[T]) WithTx(tx *gorm.DB) BaseTypedRepository[T] {
	return BaseTypedRepository[T]{
		baseRepo: r.baseRepo.WithTx(tx),
	}
}

func (r *BaseTypedRepository[T]) SavePoint(name string) error {
	return r.baseRepo.SavePoint(name)
}

func (r BaseTypedRepository[T]) Commit() error {
	return r.baseRepo.Commit()
}

func (r BaseTypedRepository[T]) Rollback() error {
	return r.baseRepo.Rollback()
}

func (r BaseTypedRepository[T]) RollbackTo(name string) error {
	return r.baseRepo.RollbackTo(name)
}

// endregion: Transaction section
