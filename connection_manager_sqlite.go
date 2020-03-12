package progorm

const sqLiteDialect = "sqlite3"

type sqliteConnectionManager struct {
	ConnectionManager
}

// Creates an instance of the SQLite implementation of the ConnectionManager interface.
// Be sure to 'import _ "github.com/jinzhu/gorm/dialects/sqlite"' package in order to run successfully
func NewSQLiteConnectionManager(dbname string, debugMode bool) ConnectionManager {
	m := &sqliteConnectionManager{
		ConnectionManager: newConnectionManager(sqLiteDialect, dbname, debugMode),
	}

	return m
}
