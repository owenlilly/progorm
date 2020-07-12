package progorm

import (
	"log"

	"github.com/jinzhu/gorm"
	"github.com/vcraescu/go-paginator"
	"github.com/vcraescu/go-paginator/adapter"
)

// Contains information about database connection and methods to access data, should be extended by more specific repository types.
type BaseRepository struct {
	connMan ConnectionManager
	db      *gorm.DB
}

// Instantiate new instance of BaseRepository
func NewBaseRepository(connMan ConnectionManager) BaseRepository {
	r := BaseRepository{
		connMan: connMan,
	}

	var err error
	r.db, err = connMan.GetConnection()
	if err != nil {
		log.Fatalf("%+v\n", err)
	}

	return r
}

// Model to insert must be a pointer/reference type
func (r BaseRepository) InsertRecord(model interface{}) error {
	return r.db.Create(model).Error
}

func (r BaseRepository) FindRecords(page, perPage uint, query *gorm.DB, out interface{}) (Page, error) {
	resultPage := Page{PerPage: perPage}

	p := paginator.New(adapter.NewGORMAdapter(query), int(perPage))
	p.SetPage(int(page))
	if err := p.Results(out); err != nil {
		return resultPage, err
	}

	resultPage.Pages = uint(p.PageNums())
	resultPage.Total = uint(p.Nums())
	resultPage.Page = uint(p.Page())

	return resultPage, nil
}

func (r BaseRepository) Count(model, query interface{}, args ...interface{}) (count int64, err error) {
	err = r.db.Model(model).Where(query, args...).Count(&count).Error

	return
}

// Create tables for the given models or return an error
func (r BaseRepository) AutoMigrate(models ...interface{}) error {
	return r.connMan.AutoMigrate(models...)
}

// Create tables for the given models or print a warning message if there's an error
func (r BaseRepository) AutoMigrateOrWarn(models ...interface{}) {
	if err := r.connMan.AutoMigrate(models...); err != nil {
		log.Println("warning:", err.Error())
	}
}

// Return the ConnectionManager
func (r BaseRepository) ConnectionManager() ConnectionManager {
	return r.connMan
}

// Returns a struct contain gorm database connection information
func (r BaseRepository) DB() *gorm.DB {
	return r.db
}
