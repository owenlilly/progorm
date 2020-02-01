package progorm

import (
	"errors"
	"reflect"
	"sync"
	"time"

	"github.com/jinzhu/gorm"
)

var (
	ErrConnectionClosed = errors.New("db connection closed")
)

type (
	ConnectionManager interface {
		GetConnection() (*gorm.DB, error)
		AutoMigrate(values ...interface{}) error
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

func newConnectionManager(dialect, connString string, debugMode bool) connectionManager {
	connMan := connectionManager{
		dialect:        dialect,
		once:           sync.Once{},
		debugMode:      debugMode,
		connString:     connString,
		migratedTables: make(map[reflect.Type]bool),
	}

	return connMan
}

func (c *connectionManager) GetConnection() (*gorm.DB, error) {
	var err error = nil

	// this func should be once executed and only once,
	// even if GetConnection() is called multiple times
	execOnceOnlyFunc := func() {
		// we'll default to UTC time for created at
		gorm.NowFunc = func() time.Time {
			return time.Now().UTC()
		}

		c.db, err = gorm.Open(c.dialect, c.connString)
		if err != nil {
			return
		}

		c.db.LogMode(c.debugMode)

		c.db.DB().SetMaxIdleConns(0)
		c.db.DB().SetMaxOpenConns(0)
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
