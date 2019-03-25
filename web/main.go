package main

import (
	"log"
	"net/http"

	"github.com/awinterman/lifting"
	"github.com/awinterman/lifting/web"
)

func main() {
	var storage, err = lifting.CreateStorage(".lift.sqlite", nil)

	if err != nil {
		panic(err)
	}

	mux := http.NewServeMux()
	fs := http.FileServer(http.Dir("static/"))

	handlers := web.Handlers{Storage: storage, Step: 10}

	mux.Handle("../stylesheets/", fs)
	mux.HandleFunc("/", handlers.Handle)

	port := ":9000"
	log.Printf("Listening http://localhost:%v", port)
	http.ListenAndServe(port, mux)
}
