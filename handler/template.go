package handler

import (
	"net/http"
	"path/filepath"
	"sync"
	"text/template"
)

type Template struct {
	Once     sync.Once
	Filename string
	Templ    *template.Template
}

func (t *Template) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.Once.Do(func() {
		t.Templ = template.Must(template.ParseFiles(filepath.Join("templates", t.Filename)))
	})
	t.Templ.Execute(w, r)
}
