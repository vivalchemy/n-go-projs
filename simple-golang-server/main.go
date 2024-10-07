package main

import (
	"html/template"
	"log"
	"net/http"
)

func main() {
	fileServer := http.FileServer(http.Dir("./templates"))
	http.Handle("/", fileServer)

	http.HandleFunc("GET /hello", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"message": "Hello World", "name": "` + r.URL.Query().Get("name") + `"}`))
	})

	// aint got the time to figure out the solution
	// http.HandleFunc("GET /hello/{name}", func(w http.ResponseWriter, r *http.Request) {
	// 	w.Header().Set("Content-Type", "application/json")
	// })

	http.HandleFunc("GET /form", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		t, _ := template.ParseFiles("./templates/form.html")
		t.Execute(w, nil)
	})

	http.HandleFunc("POST /form", func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			log.Println(err)
			return
		}

		log.Println(r.Form)

		name := r.FormValue("name")
		age := r.FormValue("age")
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"name": "` + name + `", "age": "` + age + `"}`))
	})

	log.Println("Listening on port 3000")
	if err := http.ListenAndServe(":3000", nil); err != nil {
		log.Fatal(err)
	}
}
