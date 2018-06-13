package user

import (
	"os"
	"testing"
)

func TestGetUserByID(t *testing.T) {
	TS.Use("users.sql")
	defer TS.Clear()

	user, err := GetUserByID(1)
	if err != nil {
		t.Errorf(err.Error())
	}

	if user.ID != 1 || user.Name != "foo" {
		t.Errorf("")
	}

}

func TestUpdateUserNameByID(t *testing.T) {
	TS.Use("users.sql")
	defer TS.Clear()

	UpdateUserName(1, "bar")

	user, err := GetUserByID(1)
	if err != nil {
		t.Errorf(err.Error())
	}

	if user.ID != 1 || user.Name != "bar" {
		t.Errorf("")
	}
}

func TestMain(m *testing.M) {
	initTestDB()

	r := m.Run()

	TS.DropTestDB()

	os.Exit(r)
}
