package progorm

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type sqliteConnectionManager struct {
	ConnectionManager
}

// NewSQLiteConnectionManager creates an instance of the SQLite implementation of the ConnectionManager interface.
// Be sure to 'import _ "github.com/jinzhu/gorm/dialects/sqlite"' package in order to run successfully
func NewSQLiteConnectionManager(dbname string, config *gorm.Config) ConnectionManager {
	dialector := sqlite.Open(dbname)
	m := &sqliteConnectionManager{
		ConnectionManager: newConnectionManager(dialector, config),
	}

	return m
}
