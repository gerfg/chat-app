package main

import (
	"log"
	"net/http"
	"path/filepath"
	"sync"
	"text/template"
)

func main() {
	http.Handle("/", &templateHandler{Filename: "chat.html"})

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

type templateHandler struct {
	Once     sync.Once
	Filename string
	Templ    *template.Template
}

func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.Once.Do(func() {
		t.Templ = template.Must(template.ParseFiles(filepath.Join("templates", t.Filename)))
	})
	t.Templ.Execute(w, nil)
}
