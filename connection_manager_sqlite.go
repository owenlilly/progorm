package progorm

import (
	"errors"
	"sync"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	_ "github.com/mattn/go-sqlite3"
)

const sqLiteDialect = "sqlite3"

type sqliteConnectionManager struct {
	db         *gorm.DB
	once       sync.Once
	debugMode  bool
	connString string
}

func (m sqliteConnectionManager) Debug() bool {
	return m.debugMode
}

func (m *sqliteConnectionManager) Dialect() string {
	return sqLiteDialect
}

func (m *sqliteConnectionManager) ConnString() string {
	return m.connString
}

// Creates an instance of the SQLite implementation of the ConnectionManager interface.
func NewSQLiteConnectionManager(dbname string, debugMode bool) ConnectionManager {
	m := &sqliteConnectionManager{
		connString: dbname,
		debugMode:  debugMode,
	}

	return m
}

func (m *sqliteConnectionManager) GetConnection() (*gorm.DB, error) {
	var err error
	m.once.Do(func() {
		m.db, err = gorm.Open(sqLiteDialect, m.connString)

		m.db.LogMode(m.debugMode)
	})

	return m.db, err
}

func (m sqliteConnectionManager) AutoMigrate(models ...interface{}) error {
	if m.db != nil {
		return m.db.AutoMigrate(models...).Error
	}

	return errors.New("database is nil")
}
