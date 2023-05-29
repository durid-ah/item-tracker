package services

import (
	"database/sql"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func TestItemServiceAdd(t *testing.T) {
	db, err := sql.Open("sqlite3", "../something.db")
	if err != nil {
		t.Error(err)
	}

	svc := ItemService{Db: db}
	testItem := Item{
		Label:          "Some Label",
		Description:    "Some description",
		ExpirationDate: time.Now(),
	}

	_, err = svc.Add(&testItem)
	if err != nil {
		t.Error(err)
	}
}

func TestItemService(t *testing.T) {
	db, connErr := sql.Open("sqlite3", "../db/something.db")
	if connErr != nil {
		t.Error(connErr)
	}

	svc := ItemService{Db: db}

	_, err := svc.GetAll()
	if err != nil {
		t.Error(err)
	}
}
