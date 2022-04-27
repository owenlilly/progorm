package progorm

import (
	"log"
	"math"
	"regexp"

	"github.com/owenlilly/progorm-connection/connection"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// BaseRepository base repository type for accessing tables
type BaseRepository struct {
	connMan connection.Manager
	db      *gorm.DB
	inTx    bool
}

// NewBaseRepository instantiate new instance of BaseRepository
func NewBaseRepository(connMan connection.Manager) BaseRepository {
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

// InsertRecord model to insert must be a pointer/reference type
func (r BaseRepository) InsertRecord(model any) error {
	return r.db.Create(model).Error
}

// FindRecords page finding records
func (r BaseRepository) FindRecords(page, perPage uint, query *gorm.DB, out any) (Page, error) {
	if perPage > 1000 {
		// cap at 1000 records per call
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
		DryRun: true,
		Logger: logger.New(nil, logger.Config{
			LogLevel: logger.Silent,
		}),
	})

	var total int64
	countSqlStr, countArgs := r.buildCountSql(session)
	err := r.db.Raw(countSqlStr, countArgs...).Count(&total).Error
	if err != nil {
		return results, err
	}

	queryStmt := session.
		Offset(int(page - 1*perPage)).
		Limit(int(perPage)).
		Find(nil).Statement

	// for Postgresql, db.Raw() expects '?' as placeholders so replace '$1' placeholders with '?'
	sqlStr := r.replaceNumericPlaceholders(queryStmt.SQL.String())
	vars := queryStmt.Vars

	err = r.db.Raw(sqlStr, vars...).Scan(out).Error
	if err != nil {
		return results, err
	}

	results.Total = uint(total)
	results.Pages = r.calcPageCount(uint64(results.PerPage), uint64(results.Total))

	return results, nil
}

// Count counts total number of records for the given query
func (r BaseRepository) Count(model, query any, args ...any) (count int64, err error) {
	err = r.db.Model(model).Where(query, args...).Count(&count).Error

	return
}

// AutoMigrate create tables for the given models or return an error
func (r BaseRepository) AutoMigrate(models ...any) error {
	return r.connMan.AutoMigrate(models...)
}

// AutoMigrateOrWarn create tables for the given models or print a warning message if there's an error
func (r BaseRepository) AutoMigrateOrWarn(models ...any) {
	if err := r.connMan.AutoMigrate(models...); err != nil {
		log.Println("warning:", err.Error())
	}
}

// ConnectionManager return the current ConnectionManager
func (r BaseRepository) ConnectionManager() connection.Manager {
	return r.connMan
}

// DB return a struct contain gorm database connection information
func (r BaseRepository) DB() *gorm.DB {
	return r.db
}

// region: Transaction section

// WithTx start a new database transaction
func (r BaseRepository) WithTx(tx *gorm.DB) BaseRepository {
	return BaseRepository{
		connMan: r.connMan,
		db:      tx,
		inTx:    true,
	}
}

func (r *BaseRepository) SavePoint(name string) error {
	return r.db.SavePoint(name).Error
}

func (r BaseRepository) Commit() error {
	return r.db.Commit().Error
}

func (r BaseRepository) Rollback() error {
	return r.db.Rollback().Error
}

func (r BaseRepository) RollbackTo(name string) error {
	return r.db.RollbackTo(name).Error
}

// endregion: Transaction section

func (r BaseRepository) calcPageCount(perPage, total uint64) uint {
	if perPage == 0 || total == 0 {
		return 0
	}
	return uint(math.Ceil(float64(total) / float64(perPage)))
}

func (r BaseRepository) buildCountSql(db *gorm.DB) (countSql string, vars []any) {
	if orderByClause, ok := db.Statement.Clauses["ORDER BY"]; ok {
		if _, ok := db.Statement.Clauses["GROUP BY"]; !ok {
			delete(db.Statement.Clauses, "ORDER BY")
			defer func() {
				db.Statement.Clauses["ORDER BY"] = orderByClause
			}()
		}
	}
	var count int64
	countStmt := db.Count(&count).Statement

	// for Postgresql, db.Raw() expects '?' as placeholders so replace '$1' placeholders with '?'
	countSql = r.replaceNumericPlaceholders(countStmt.SQL.String())
	vars = countStmt.Vars

	return
}

func (r BaseRepository) replaceNumericPlaceholders(sqlStr string) string {
	var numericPlaceholder = regexp.MustCompile("\\$(\\d+)")

	return numericPlaceholder.ReplaceAllString(sqlStr, "?")
}
