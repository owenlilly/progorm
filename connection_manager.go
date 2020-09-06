package progorm

import (
	"database/sql"
	"errors"
	"log"
	"os"
	"reflect"
	"sync"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	// ErrConnectionClosed returned if database connection isn't open
	ErrConnectionClosed = errors.New("db connection closed")

	// ErrInvalidConnectionString returned if database is connection string is malformed/invalid
	ErrInvalidConnectionString = errors.New("invalid connection string")
)

type (
	// ConnectionManager manages database connections
	ConnectionManager interface {
		GetConnection() (*gorm.DB, error)
		AutoMigrate(tables ...interface{}) error
		AutoMigrateOrWarn(tables ...interface{})
		Dialect() string
		ConnString() string
	}

	// connectionManager implements ConnectionManager interface
	connectionManager struct {
		dialector      gorm.Dialector
		config         *gorm.Config
		connStr        string
		db             *gorm.DB
		once           sync.Once
		migratedTables map[reflect.Type]bool
	}
)

func newConnectionManager(dialector gorm.Dialector, config *gorm.Config) ConnectionManager {
	connMan := &connectionManager{
		dialector:      dialector,
		config:         config,
		once:           sync.Once{},
		migratedTables: make(map[reflect.Type]bool),
	}

	if connMan.config == nil {
		defaultLogger := logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
			logger.Config{
				SlowThreshold: time.Second,  // Slow SQL threshold
				LogLevel:      logger.Error, // Log level
				Colorful:      true,         // Disable color
			},
		)
		connMan.config = &gorm.Config{Logger: defaultLogger}
	}

	// open database connection
	_, err := connMan.GetConnection()
	if err != nil {
		log.Fatalf("failed to connect to database: %s", err.Error())
	}

	return connMan
}

// GetConnection get current *gorm.DB connection
func (c *connectionManager) GetConnection() (*gorm.DB, error) {
	var err error

	// this func should be once executed and only once,
	// even if GetConnection() is called multiple times
	execOnceOnlyFunc := func() {
		c.db, err = gorm.Open(c.dialector, c.config)
		if err != nil {
			return
		}

		var sqlDB *sql.DB
		sqlDB, err = c.db.DB()
		if err != nil {
			return
		}
		sqlDB.SetMaxIdleConns(10)
		sqlDB.SetMaxOpenConns(-1)
	}

	// ensure execOnceOnlyFunc() is only ever executed once
	c.once.Do(execOnceOnlyFunc)

	return c.db, err
}

// AutoMigrate create/change database table definition based on the given models
func (c *connectionManager) AutoMigrate(tables ...interface{}) error {
	if c.db == nil {
		return ErrConnectionClosed
	}

	var unmigratedTables []interface{}
	for _, table := range tables {
		t := reflect.ValueOf(table).Type()
		if !c.migratedTables[t] {
			// add current table to list of tables to be migrated
			unmigratedTables = append(unmigratedTables, table)
			// mark current table as migrated
			c.migratedTables[t] = true
		}
	}

	return c.db.AutoMigrate(unmigratedTables...)
}

// AutoMigrateOrWarn same as AutoMigrate but prints a log instead of returning an error
func (c *connectionManager) AutoMigrateOrWarn(tables ...interface{}) {
	if err := c.AutoMigrate(tables...); err != nil {
		log.Printf("%v\n", err)
	}
}

// Dialect return the current database dialect
func (c *connectionManager) Dialect() string {
	return c.config.Name()
}

// ConnString return the connection string for the current database
func (c *connectionManager) ConnString() string {
	return c.connStr
}
