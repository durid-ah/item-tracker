package main

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func ConnectDatabase() *sql.DB {
	db, err := sql.Open("sqlite3", "something.db")
	if err != nil {
		panic(err)
	}

	return db
}

func CreateItemsTable(db *sql.DB) {

	tx,txErr := db.Begin()
	if txErr != nil {
		log.Fatal(txErr)
	}

	stmt, stmtErr := tx.Prepare(`
		CREATE TABLE items (
			id INTEGER PRIMARY KEY,
			user_id INTEGER,
			label TEXT NOT NULL,
			description TEXT NOT NULL,
			container_number INTEGER NOT NULL,
			expiration_date TEXT NOT NULL,
			image BLOB,
			FOREIGN KEY (user_id) REFERENCES users (id)
	);`)
	if stmtErr != nil {
		log.Fatal(stmtErr)
	}

	defer stmt.Close()

	_, err := stmt.Exec()
	if err != nil {
		log.Fatal(err)
	}

	tx.Commit()
}

func CreateUsersTable(db *sql.DB) {

	tx,txErr := db.Begin()
	if txErr != nil {
		log.Fatal(txErr)
	}

	stmt, stmtErr := tx.Prepare(`
		CREATE TABLE users (
			id INTEGER PRIMARY KEY,
			username TEXT NOT NULL UNIQUE,
			password TEXT NOT NULL
	);`)
	if stmtErr != nil {
		log.Fatal(stmtErr)
	}

	defer stmt.Close()

	_, err := stmt.Exec()
	if err != nil {
		log.Fatal(err)
	}

	tx.Commit()
}

func main() {
	db := ConnectDatabase()

	CreateUsersTable(db)
	CreateItemsTable(db)
}