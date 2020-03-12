package progorm

import (
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

const sqLiteDialect = "sqlite3"

type sqliteConnectionManager struct {
	ConnectionManager
}

// Creates an instance of the SQLite implementation of the ConnectionManager interface.
func NewSQLiteConnectionManager(dbname string, debugMode bool) ConnectionManager {
	m := &sqliteConnectionManager{
		ConnectionManager: newConnectionManager(sqLiteDialect, dbname, debugMode),
	}

	return m
}
