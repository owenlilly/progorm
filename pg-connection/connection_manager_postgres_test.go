package pg_connection

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMakePostgresConnectionString(t *testing.T) {
	connStr1 := MakePostgresConnString("user", "pass", "localhost", "testdb", "disable")
	assert.Equal(t, "postgres://user:pass@localhost/testdb?sslmode=disable", connStr1)

	connStr2 := MakePostgresConnString("", "pass", "localhost", "testdb", "disable")
	assert.Equal(t, "postgres://localhost/testdb?sslmode=disable", connStr2)

	connStr3 := MakePostgresConnString("", "", "localhost", "testdb", "disable")
	assert.Equal(t, "postgres://localhost/testdb?sslmode=disable", connStr3)

	connStr4 := MakePostgresConnString("", "", "localhost", "", "disable")
	assert.Equal(t, "postgres://localhost/postgres?sslmode=disable", connStr4)

	connStr5 := MakePostgresConnString("", "", "localhost", "", "disable", "defaultdb")
	assert.Equal(t, "postgres://localhost/defaultdb?sslmode=disable", connStr5)

	connStr6 := MakePostgresConnString("user", "pass", "localhost", "testdb", "require")
	assert.Equal(t, "postgres://user:pass@localhost/testdb?sslmode=require", connStr6)

	connStr7 := MakePostgresConnString("user", "pass", "", "testdb", "require")
	assert.Equal(t, "postgres://user:pass@localhost/testdb?sslmode=require", connStr7)
}
