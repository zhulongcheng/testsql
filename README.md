# TestSQL

[![Build Status](https://api.travis-ci.com/zhulongcheng/testsql.svg?branch=master)](https://travis-ci.org/zhulongcheng/testsql)
[![codecov.io](https://codecov.io/github/zhulongcheng/testsql/branch/master/graph/badge.svg)](https://codecov.io/github/zhulongcheng/testsql)
[![Go Report Card](https://goreportcard.com/badge/github.com/zhulongcheng/testsql)](https://goreportcard.com/report/github.com/zhulongcheng/testsql)
[![GoDoc](https://godoc.org/github.com/zhulongcheng/testsql?status.svg)](https://godoc.org/github.com/zhulongcheng/testsql)


Generate test data from SQL files before testing and clear it after finished.

## Contents

- [Installation](#installation)
- [Usage](#usage)
- [API Reference](#api-reference)
- [See also](#see-also)


## Installation
```
go get github.com/zhulongcheng/testsql
```

## Usage
Create a folder for the table-schema file and sql files:
```bash
testsql
├── fixtures
│   └── users.sql
└── schema.sql
```

 The table-schema file include all tables schema, it would be like this:
```sql
CREATE TABLE `users` (
    `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
    `name` varchar(32) CHARACTER SET utf8mb4 NOT NULL,
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
```

The sql file would be like this:
```sql
INSERT INTO `users` (`id`, `name`)
VALUES
    (1, 'foo');
```

Your tests would be like this:
```go
var TS *testsql.TestSQL

func newTS() *testsql.TestSQL {
    dsn := "user:password@tcp(host:port)/test_db_name"
    tableSchemaPath := "testsql/schema.sql"
    dirPath := "testsql/fixtures"
    ts := testsql.New(dsn, tableSchemaPath, dirPath)
    return ts
}

func initTestDB() {
    TS = newTS()
    
    // set sql-driver/orm to read/write data from TS's DSN
    // Driver = sql.Open(TS.Config.DSN)
    // ORM = ORM.New(TS.Config.DSN)
}

func TestMain(m *testing.M) {
    initTestDB()
    r := m.Run()
    TS.DropTestDB()
    os.Exit(r)
}

func TestUser(t *testing.T) {
    TS.Use("users.sql")
    defer TS.Clear()
    
    // user := GetUserByID(1)
    // if user.name != "foo" {
    //    t.Errorf("not equal, expected: %s, actual: %s", "foo", user.name) 
    // }
}
```

## API Reference 
`TestSQL.Exec` generates test data from sql string.

`TestSQL.Use` generates test data from sql files.

`TestSQL.Clear` deletes data of all tables

`TestSQL.DropTestDB` drops the database created by TestSQL.

## FQA
### Set sql-driver/orm to read/write data from TestSQL's DSN

```go
    db := yourSQLDriver.Open(TS.Config.DSN)
    // or db := yourORM.Open(TS.Config.DSN)

    yourModel.SetDB(db)  // reset db
```


## See also
[Examples](https://github.com/zhulongcheng/testsql/tree/master/examples)

