// Package testsql generates test data from SQL files before testing and clears it after finished.
package testsql

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"regexp"
	"strings"

	"math/rand"
	"path/filepath"
	"time"

	"github.com/go-sql-driver/mysql"
)

const (
	driverName     string = "mysql"
	maxDBRandValue int    = 10000
	dbNameTemplate string = "%s_%d_%d"
)

// Config represents a configuration for TestSQL.
type Config struct {
	DSN             string
	TableSchemaPath string
	FixtureDirPath  string
}

// TestSQL is the main object to generate/clear test data.
type TestSQL struct {
	Config *Config
	DB     *sql.DB
	tables *Set
}

// New initializes sql.DB, creates test database, and creates all tables.
func New(dsn string, tableSchemaPath string, fixtureDirPath string) *TestSQL {
	ts := &TestSQL{}
	dsn = replaceDBName(dsn)
	config := &Config{
		DSN:             dsn,
		TableSchemaPath: tableSchemaPath,
		FixtureDirPath:  fixtureDirPath,
	}

	ts.Config = config
	ts.createTestDB()

	db, err := sql.Open(driverName, dsn)
	panicErr(err)

	ts.DB = db
	ts.tables = NewSet()
	ts.createTables()
	return ts
}

// Exec generates data from sql string.
func (ts *TestSQL) Exec(sqlString string) {
	_, err := ts.DB.Exec(sqlString)
	panicErr(err)
}

// Use generates data from sql files.
func (ts *TestSQL) Use(sqlFileNames ...string) {
	for _, sqlFileName := range sqlFileNames {
		sqlFilePath := filepath.Join(ts.Config.FixtureDirPath, sqlFileName)
		ts.sqlExec(sqlFilePath)
	}
}

// Clear deletes data of all tables.
func (ts *TestSQL) Clear() {
	for _, tableName := range ts.tables.Values() {
		ts.clearTable(tableName)
	}
}

// DropTestDB drops the database created by TestSQL.
func (ts *TestSQL) DropTestDB() {
	conf, err := mysql.ParseDSN(ts.Config.DSN)
	panicErr(err)

	dbName := conf.DBName

	conf.DBName = ""
	dsn := conf.FormatDSN()

	db, err := sql.Open(driverName, dsn)
	panicErr(err)
	defer db.Close()

	_, err = db.Exec("DROP DATABASE " + dbName)
	panicErr(err)

	log.Printf("DROP DATABASE %s", dbName)
}

func (ts *TestSQL) createTestDB() {
	conf, err := mysql.ParseDSN(ts.Config.DSN)
	panicErr(err)

	dbName := conf.DBName

	conf.DBName = ""
	dsn := conf.FormatDSN()

	db, err := sql.Open(driverName, dsn)
	panicErr(err)
	defer db.Close()

	_, err = db.Exec("CREATE DATABASE " + dbName)
	panicErr(err)
}

func (ts *TestSQL) createTables() {
	ts.sqlExec(ts.Config.TableSchemaPath)
	ts.recordTables()
}

func (ts *TestSQL) clearTable(tableName string) {
	sqlString := fmt.Sprintf("DELETE FROM %s", tableName)
	_, err := ts.DB.Exec(sqlString)
	panicErr(err)
}

func (ts *TestSQL) sqlExec(sqlFilePath string) {
	data, err := ioutil.ReadFile(sqlFilePath)
	panicErr(err)

	// database/sql/driver not support multiple queries/results
	sqlString := string(data)
	for _, s := range strings.Split(sqlString, ";") {
		if len(strings.TrimSpace(s)) == 0 {
			continue
		}
		_, err = ts.DB.Exec(s)
		panicErr(err)
	}
}

func (ts *TestSQL) recordTables() {
	data, err := ioutil.ReadFile(ts.Config.TableSchemaPath)
	panicErr(err)

	sqlString := string(data)
	for _, s := range strings.Split(sqlString, ";") {
		tableSQL := strings.TrimSpace(s)
		if len(tableSQL) == 0 {
			continue
		}

		tableName := ts.extraTableName(tableSQL)
		if tableName == "" {
			msg := fmt.Sprintf("extra table name failed, please check sql: %s", tableSQL)
			panic(msg)
		}

		ts.tables.Add(tableName)
	}
}

func (ts *TestSQL) extraTableName(tableSQL string) string {
	re, err := regexp.Compile("(?i)(create table)\\s+[`'\"]?(\\w+)")
	panicErr(err)

	lines := re.FindAllString(tableSQL, 1)
	if lines != nil {
		// lines[0]: "create table tableName"
		words := strings.Split(lines[0], " ")
		tableName := words[len(words)-1]
		if strings.HasPrefix(tableName, "`") || strings.HasPrefix(tableName, "'") || strings.HasPrefix(tableName, "\"") {
			tableName = tableName[1:]
		}
		return tableName
	}

	return ""
}

func replaceDBName(dsn string) string {
	conf, err := mysql.ParseDSN(dsn)
	panicErr(err)

	// check db name
	if !strings.HasPrefix(conf.DBName, "test_") {
		panic("DB name must starts with `test_`, e.g. test_some_db")
	}

	now := time.Now().UnixNano()
	num := rand.Intn(maxDBRandValue)
	conf.DBName = fmt.Sprintf(dbNameTemplate, conf.DBName, now, num)

	log.Printf("test db name: %s", conf.DBName)
	return conf.FormatDSN()
}

func panicErr(err error) {
	if err != nil {
		panic(err)
	}
}

func init() {
	rand.Seed(time.Now().UnixNano())
}
