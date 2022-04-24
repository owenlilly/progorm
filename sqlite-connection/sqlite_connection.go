package sqlite_connection

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/owenlilly/progorm/connection"
)

type sqliteConnectionManager struct {
	connection.Manager
}

// NewConnectionManager creates an instance of the SQLite implementation of the Manager interface.
func NewConnectionManager(dbname string, config *gorm.Config) connection.Manager {
	dialector := sqlite.Open(dbname)
	m := &sqliteConnectionManager{
		Manager: connection.NewBaseConnectionManager(dbname, dialector, config),
	}

	return m
}
