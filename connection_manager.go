package progorm

import (
	"errors"
	"log"
	"reflect"
	"sync"

	"github.com/jinzhu/gorm"
)

var (
	// Returned if database connection isn't open
	ErrConnectionClosed = errors.New("db connection closed")

	ErrInvalidConnectionString = errors.New("invalid connection string")
)

type (
	// Manages database connections
	ConnectionManager interface {
		GetConnection() (*gorm.DB, error)
		AutoMigrate(tables ...interface{}) error
		AutoMigrateOrWarn(tables ...interface{})
		Debug() bool
		Dialect() string
		ConnString() string
	}

	// implements ConnectionManager interface
	connectionManager struct {
		dialect        string
		db             *gorm.DB
		once           sync.Once
		debugMode      bool
		connString     string
		migratedTables map[reflect.Type]bool
	}
)

func newConnectionManager(dialect, connString string, debugMode bool) ConnectionManager {
	connMan := &connectionManager{
		dialect:        dialect,
		once:           sync.Once{},
		debugMode:      debugMode,
		connString:     connString,
		migratedTables: make(map[reflect.Type]bool),
	}

	// open database connection
	_, _ = connMan.GetConnection()

	return connMan
}

func (c *connectionManager) GetConnection() (*gorm.DB, error) {
	var err error

	// this func should be once executed and only once,
	// even if GetConnection() is called multiple times
	execOnceOnlyFunc := func() {
		c.db, err = gorm.Open(c.dialect, c.connString)
		if err != nil {
			return
		}

		c.db.LogMode(c.debugMode)

		c.db.DB().SetMaxIdleConns(5)
		c.db.DB().SetMaxOpenConns(-1)
	}

	// ensure execOnceOnlyFunc() is only ever executed once
	c.once.Do(execOnceOnlyFunc)

	return c.db, err
}

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

	return c.db.AutoMigrate(unmigratedTables...).Error
}

func (c *connectionManager) AutoMigrateOrWarn(tables ...interface{}) {
	if err := c.AutoMigrate(tables...); err != nil {
		log.Printf("%v\n", err)
	}
}

func (c *connectionManager) Debug() bool {
	return c.debugMode
}

func (c *connectionManager) Dialect() string {
	return c.dialect
}

func (c *connectionManager) ConnString() string {
	return c.connString
}
