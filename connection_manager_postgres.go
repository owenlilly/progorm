package progorm

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/lib/pq"
)

type postgresConnectionManager struct {
	connectionManager
}

func NewPostgresConnectionManager(connString string, debugMode bool) ConnectionManager {
	connMan := &postgresConnectionManager{
		connectionManager: newConnectionManager("postgres", connString, debugMode),
	}

	return connMan
}

func MakePostgresConnString(user, pass, host, dbName, sslMode string, defaultsDBs ...string) string {
	var connStr = "postgres://"

	var defaultDB = "postgres"
	if defaultsDBs != nil && len(defaultsDBs) > 0 {
		defaultDB = defaultsDBs[0]
	}

	if user != "" {
		connStr += fmt.Sprintf("%s:%s@", user, pass)
	}

	if host != "" {
		connStr += host
	} else {
		connStr += "localhost"
	}

	if dbName != "" {
		connStr += "/" + dbName + "?sslmode=" + sslMode
	} else {
		connStr += fmt.Sprintf("/%s?sslmode=%s", defaultDB, sslMode)
	}

	return connStr
}

// Creates postgres database of the given name if one doesn't already exists.
// No actions are performed if the database already exists.
func PGCreateDbIfNotExists(user, pass, host, name, sslMode string, defaultDB ...string) error {
	if defaultDB != nil && len(defaultDB) > 0 {
		if name == defaultDB[0] {
			return nil
		}
	} else if name == "" {
		return errors.New("database name is required")
	}

	var connStr = MakePostgresConnString(user, pass, host, "", sslMode, defaultDB...)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return err
	}
	defer func() {
		_ = db.Close()
	}()

	_, err = db.Exec("CREATE DATABASE " + name)
	if err != nil {
		switch e := err.(type) {
		case *pq.Error:
			if strings.Contains(e.Message, "already exists") {
				return nil
			}
		}
		return err
	}

	return db.Close()
}
