package main

import (
	"main/pkg/api"
	"main/pkg/db"
	"net/http"
)

func main() {
	err := db.Init("scheduler.db")
	if err != nil {
		panic(err)
	}

	api.Init()

	http.Handle("/", http.FileServer(http.Dir("web")))
	err = http.ListenAndServe(":7540", nil)
	if err != nil {
		panic(err)
	}
}
