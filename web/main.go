package main

import (
	"log"
	"net/http"

	"github.com/awinterman/lifting"
	"github.com/awinterman/lifting/postgres"
)

func main() {
	var storage, err = postgres.CreateStorage("user=andrew password=andrew dbname=andrew", nil)

	if err != nil {
		panic(err)
	}

	mux := http.NewServeMux()
	fs := http.FileServer(http.Dir("./static/"))

	handlers := lifting.Handlers{Storage: storage, Step: 10}

	mux.Handle("/stylesheets/", fs)
	mux.HandleFunc("/", handlers.Handle)

	port := ":9000"
	log.Printf("Listening http://localhost:%v", port)
	http.ListenAndServe(port, mux)
}
