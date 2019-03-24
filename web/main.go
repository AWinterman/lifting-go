package main

import (
	"log"
	"net/http"

	"github.com/awinterman/lifting"
)

func main() {
	var storage, err = lifting.CreateStorage(".lift.sqlite", nil)

	if err != nil {
		panic(err)
	}

	mux := http.NewServeMux()
	fs := http.FileServer(http.Dir("static/"))

	handlers := Handlers{Storage: storage, step: 10}

	mux.Handle("/stylesheets/", fs)
	mux.HandleFunc("/", handlers.handle)

	port := ":9000"
	log.Printf("Listening http://localhost:%v", port)
	http.ListenAndServe(port, mux)
}
