package main

import (
	"database/sql"
	"log"
	"net/http"

	itemendpoints "github.com/durid-ah/item-tracker/app/item_endpoints"
	userendpoints "github.com/durid-ah/item-tracker/app/user_endpoints"
	"github.com/durid-ah/item-tracker/helpers"

	_ "github.com/mattn/go-sqlite3"
)

func ConnectDatabase() *sql.DB {
	db, err := sql.Open("sqlite3", "./something.db")
	if err != nil {
		log.Fatal(err)
	}

	return db
}

func main() {
	db := ConnectDatabase()
	log.Println("Database connected...")

	http.Handle("/signin", userendpoints.Signin(db))

	http.Handle("/items", helpers.WithAuth(itemendpoints.GetItemsHandler(db)))

	log.Println("Listening localhost:8080")
	http.ListenAndServe("localhost:8080", nil)
}
