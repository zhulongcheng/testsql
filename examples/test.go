package user

import (
	"github.com/zhulongcheng/testsql"
)

// TS represents a TestSQL obj
var TS *testsql.TestSQL

func newTS() *testsql.TestSQL {
	dsn := "testsql:password@tcp(localhost:3306)/test_testsql"
	tableSchemaPath := "../db/schema.sql"
	dirPath := "../db/fixtures"

	ts := testsql.New(dsn, tableSchemaPath, dirPath)
	return ts
}

func initTestDB() {
	TS = newTS()
	DB = TS.DB
}
