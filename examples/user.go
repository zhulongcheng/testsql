package user

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql" // init
)

// User represents a user
type User struct {
	ID   int
	Name string
}

// DB represents a sql.DB obj
var DB *sql.DB

func init() {
	dsn := "root:@/testsql"
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}

	DB = db
	log.Printf("init db: %+v", DB)
}

// GetUserByID returns a user
func GetUserByID(id int) (*User, error) {
	user := User{}
	err := DB.QueryRow("select id, name from users where id=?", id).Scan(&user.ID, &user.Name)
	if err != nil {
		return &user, err
	}

	log.Printf("use db: %+v", DB)
	return &user, nil
}

// UpdateUserName update user's name
func UpdateUserName(id int, name string) error {
	_, err := DB.Exec("update users set name=? where id=?", name, id)
	if err != nil {
		return err
	}

	return nil
}
