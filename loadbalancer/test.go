package main

import (
	"fmt"
	"net/http"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run test.go <port_number>")
		os.Exit(1)
	}

	port := os.Args[1]
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(r.URL.Host)
		fmt.Fprintf(w, "Hello World! from "+r.URL.Host+r.URL.Path)
	})
	fmt.Println("Listening on port " + port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		panic(err)
	}
}
