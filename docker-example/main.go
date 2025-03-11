package main

import (
	"fmt"
	"net/http"
)

const (
	port = ":8080"
)

func main() {
	http.HandleFunc("GET /ping", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-type", "text/plain")
		fmt.Fprintf(w, "pong")
	})

	fmt.Println("Listening on ", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		panic(err)
	}
}
