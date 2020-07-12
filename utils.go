package progorm

import (
	"errors"
	"strings"
)

var (
	ErrUnsupportedDatabase = errors.New("unsupported database")
)

func Min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func Max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// Returns a ConnectionManager implementation for the provided connection string.  Will return
// 'ErrUnsupportedDatabase' error if an unsupported connection string is provided.  It will
// also attempt to create a database of the given if one doesn't already exists.
func NewConnectionManager(connStr string, debugMode bool, defaultDB ...string) (ConnectionManager, error) {
	switch {
	case strings.HasPrefix(connStr, "postgres://"):
		if err := PGCreateDbIfNotExists(connStr, defaultDB...); err != nil {
			return nil, err
		}
		return NewPostgresConnectionManager(connStr, debugMode), nil
	case strings.HasPrefix(connStr, "sqlite3"):
		connStr = strings.TrimPrefix(connStr, "sqlite3://")
		return NewSQLiteConnectionManager(connStr, debugMode), nil
	default:
		return nil, ErrUnsupportedDatabase
	}
}
