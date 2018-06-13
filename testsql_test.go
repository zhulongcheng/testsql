package testsql

import (
	"strings"
	"testing"

	"database/sql"

	"github.com/go-sql-driver/mysql"
)

func newTS() *TestSQL {
	dsn := "testsql:password@tcp(localhost:3306)/test_testsql"
	tableSchemaPath := "./db/schema.sql"
	dirPath := "./db/fixtures"

	ts := New(dsn, tableSchemaPath, dirPath)
	return ts
}

func TestTS_New(t *testing.T) {
	dsn := "testsql:password@tcp(localhost:3306)/test_testsql"
	tableSchemaPath := "./db/schema.sql"
	dirPath := "./db/fixtures"

	ts := New(dsn, tableSchemaPath, dirPath)

	// check ts db name
	cfg, err := mysql.ParseDSN(ts.Config.DSN)
	if err != nil {
		t.Errorf(err.Error())
	}

	if !strings.HasPrefix(cfg.DBName, "test_testsql") {
		t.Errorf("db name error, %s not start with %s", cfg.DBName, "test_testsql")

	}

	// check fixture path
	if ts.Config.FixtureDirPath != dirPath {
		t.Errorf("fixture path error, %s != %s", ts.Config.FixtureDirPath, dirPath)
	}
}

func TestTS_Exec(t *testing.T) {
	ts := newTS()

	sqlString := `
		INSERT INTO users (id, name)
		VALUES
    		(1, 'foo')
	`
	ts.Exec(sqlString)

	var name string
	err := ts.DB.QueryRow("select name from users limit 1").Scan(&name)
	if err != nil {
		t.Errorf(err.Error())
	}
	if name != "foo" {
		t.Errorf("")
	}
}

func TestTS_Use(t *testing.T) {
	ts := newTS()
	ts.Use("users.sql")

	var name string
	err := ts.DB.QueryRow("select name from users limit 1").Scan(&name)
	if err != nil {
		t.Errorf(err.Error())
	}
	if name != "foo" {
		t.Errorf("")
	}
}

func TestTS_Clear(t *testing.T) {
	ts := newTS()
	ts.Use("users.sql")
	ts.Clear()

	var name string
	err := ts.DB.QueryRow("select name from users limit 1").Scan(&name)
	if err == nil {
		t.Errorf("TestTS Clear failed, get db result: %s", name)
	}

	if err.Error() != "sql: no rows in result set" {
		t.Errorf("TestTS Clear failed, err: %s", err.Error())
	}

}

func TestTS_DropTestDB(t *testing.T) {
	ts := newTS()

	// new sql.DB
	conf, err := mysql.ParseDSN(ts.Config.DSN)
	if err != nil {
		t.Errorf(err.Error())
	}

	TestDBName := conf.DBName

	conf.DBName = ""
	dsn := conf.FormatDSN()

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		t.Errorf(err.Error())
	}

	// check if database exists
	var name string
	err = db.QueryRow("SELECT SCHEMA_NAME FROM INFORMATION_SCHEMA.SCHEMATA WHERE SCHEMA_NAME = ?", TestDBName).Scan(&name)
	if err != nil {
		t.Errorf(err.Error())
	}
	if name != TestDBName {
		t.Errorf("TestTS DropTestDB failed, get db name: %s", name)
	}

	// call TS.DropTestDB and check again
	ts.DropTestDB()
	err = db.QueryRow("SELECT SCHEMA_NAME FROM INFORMATION_SCHEMA.SCHEMATA WHERE SCHEMA_NAME = ?", TestDBName).Scan(&name)
	if err == nil {
		t.Errorf("TestTS DropTestDB failed")
	}
}
