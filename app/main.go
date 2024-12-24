package main

import (
	"database/sql"
	"log"
	"net/http"
	
	"embed"
	"io/fs"
	"os"

	itemendpoints "github.com/durid-ah/item-tracker/app/item_endpoints"
	userendpoints "github.com/durid-ah/item-tracker/app/user_endpoints"
	"github.com/durid-ah/item-tracker/helpers"

	_ "github.com/mattn/go-sqlite3"
)

//go:embed public
var public embed.FS


func ConnectDatabase() *sql.DB {
	db, err := sql.Open("sqlite3", "./something.db")
	if err != nil {
		log.Fatal(err)
	}

	return db
}

func main() {
	
	path, _ := os.Getwd()
	log.Printf("Executable: %s \n", path)
	
	publicFS, err := fs.Sub(public, "public")
	if err != nil {
		log.Fatal(err)
	}
	
	log.Println("Directory setup")
	http.Handle("/", http.FileServer(http.FS(publicFS)))

	db := ConnectDatabase()
	log.Println("Database connected...")

	log.Println("Setting up user endpoints")
	http.Handle("/api/signin", userendpoints.Signin(db))
	http.Handle("/api/refresh", userendpoints.Refresh(db))
	http.HandleFunc("/api/signout", userendpoints.Signout)

	log.Println("Setting up item endpoints")
	http.Handle("/api/items/add", helpers.WithAuth(itemendpoints.AddItemHandler(db)))
	http.Handle("/api/items/list", helpers.WithAuth(itemendpoints.GetItemsHandler(db)))
	http.Handle("/api/items/update", helpers.WithAuth(itemendpoints.UpdateItemHandler(db)))
	http.Handle("/api/items/delete", helpers.WithAuth(itemendpoints.DeleteItemHandler(db)))
	
	log.Println("Listening localhost:8080")

	http.ListenAndServe(":8080", nil)
}
