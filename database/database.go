package database

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func ConnectDatabase() *sql.DB {
	db, err := sql.Open("sqlite3", "/var/db/item-tracker.db")
	if err != nil {
		panic(err)
	}

	return db
}

func createItemsTable(db *sql.DB) {
	stmt, stmtErr := db.Prepare(`
		CREATE TABLE IF NOT EXISTS items (
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
}

func createUsersTable(db *sql.DB) {
	stmt, stmtErr := db.Prepare(`
		CREATE TABLE IF NOT EXISTS users (
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
}

func SetupDb(db *sql.DB) {

	createUsersTable(db)
	createItemsTable(db)

	// userSvc := services.UserService{ Db: db }
	// userSvc.Add(&services.User{ Username: "user1", Password: []byte("password")})
	// userSvc.Add(&services.User{ Username: "user2", Password: []byte("password")})
}