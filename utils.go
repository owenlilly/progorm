package progorm

import (
	"errors"
	"log"
	"os"
	"strings"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"
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
		return NewPostgresConnectionManager(connStr, &gorm.Config{
			Logger: logger.New(
				log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
				logger.Config{
					SlowThreshold: time.Second,  // Slow SQL threshold
					LogLevel:      logger.Error, // Log level
					Colorful:      true,         // Disable color
				},
			),
		}), nil
	case strings.HasPrefix(connStr, "sqlite3"):
		connStr = strings.TrimPrefix(connStr, "sqlite3://")
		return NewSQLiteConnectionManager(connStr, &gorm.Config{
			Logger: logger.New(
				log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
				logger.Config{
					SlowThreshold: time.Second,  // Slow SQL threshold
					LogLevel:      logger.Error, // Log level
					Colorful:      true,         // Disable color
				},
			),
		}), nil
	default:
		return nil, ErrUnsupportedDatabase
	}
}
