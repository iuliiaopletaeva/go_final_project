package main

import (
	"database/sql"
	"main/pkg/api"
	"main/pkg/db"
	"net/http"
)

var database *sql.DB

func main() {
	err := db.Init("scheduler.db")
	if err != nil {
		panic(err)
	}
	defer database.Close()

	api.Init()

	http.Handle("/", http.FileServer(http.Dir("web")))
	err = http.ListenAndServe(":7540", nil)
	if err != nil {
		panic(err)
	}
}
