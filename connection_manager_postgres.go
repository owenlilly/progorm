package progorm

import (
	"database/sql"
	"fmt"
	"regexp"
	"strings"

	"github.com/lib/pq"
)

type postgresConnectionManager struct {
	ConnectionManager
}

// Creates a new instance of the Postgres implementation of the ConnectionManager interface.
func NewPostgresConnectionManager(connString string, debugMode bool) ConnectionManager {
	connMan := &postgresConnectionManager{
		ConnectionManager: newConnectionManager("postgres", connString, debugMode),
	}

	return connMan
}

// Builds Postgres connection string from individual credential parts
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

// Creates postgres database of the given name if one doesn't already exists. No actions are performed if the database already exists.
func PGCreateDbIfNotExists(connString string, defaultDBs ...string) error {
	var defaultDB string
	if defaultDBs != nil && len(defaultDBs) > 0 {
		if defaultDB == defaultDBs[0] {
			return nil
		}
	} else {
		defaultDB = "postgres"
	}

	re := regexp.MustCompile(`(?m)postgres://.+:?\d?/(\w+)`)
	dbName := re.FindStringSubmatch(connString)[1]
	connStrWithDefaultDB := strings.Replace(connString, dbName, defaultDB, 1)

	db, err := sql.Open("postgres", connStrWithDefaultDB)
	if err != nil {
		return err
	}
	defer func() {
		_ = db.Close()
	}()

	_, err = db.Exec("CREATE DATABASE " + dbName)
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
