package progorm

import (
	"log"
	"math"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"
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
	if perPage > 1000 {
		// cap at 10000 records per call
		perPage = 1000
	}

	if page == 0 {
		// 1 based page index
		page = 1
	}

	results := Page{
		PerPage: perPage,
		Page:    page,
	}

	session := query.Session(&gorm.Session{
		DryRun:         true,
		WithConditions: true,
		Logger: logger.New(nil, logger.Config{
			LogLevel: logger.Silent,
		}),
	})

	var total int64
	countStmt := session.Count(&total).Statement
	countSqlStr := countStmt.SQL.String()
	countArgs := countStmt.Vars

	err := r.db.Raw(countSqlStr, countArgs...).Count(&total).Error
	if err != nil {
		return results, err
	}

	queryStmt := session.
		Offset(int(page - 1*perPage)).
		Limit(int(perPage)).
		Find(nil).Statement

	sqlStr := queryStmt.SQL.String()
	vars := queryStmt.Vars

	err = r.db.Raw(sqlStr, vars...).Scan(out).Error
	if err != nil {
		return results, err
	}

	results.Total = uint(total)
	results.Pages = r.calcPageCount(uint64(results.PerPage), uint64(results.Total))

	return results, nil
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

func (r BaseRepository) calcPageCount(perPage, total uint64) uint {
	if perPage == 0 || total == 0 {
		return 0
	}
	return uint(math.Ceil(float64(total) / float64(perPage)))
}
