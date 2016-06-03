package main

import (
	"github.com/remotejob/gojobextractor/dbhandler"
	"gopkg.in/mgo.v2"
	"log"
	"net/http"
	"path/filepath"
	"sync"
	"text/template"
)

type templateHandler struct {
	once     sync.Once
	filename string
	templ    *template.Template
}

type updateHandler struct{}

func (t *updateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	if r.Method == "POST" {

		dbsession, err := mgo.Dial("127.0.0.1")
		if err != nil {
			panic(err)
		}
		defer dbsession.Close()

		r.ParseForm()

		id := r.Form["Id"][0]
		log.Println("Update POST", id)
		dbhandler.UpdateEmployerById(*dbsession, id)

		http.Redirect(w, r, "/", http.StatusFound)

	}

}

// ServeHTTP handles the HTTP request.
func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	t.once.Do(func() {
		t.templ = template.Must(template.ParseFiles(filepath.Join("templates", t.filename)))
	})

	dbsession, err := mgo.Dial("127.0.0.1")
	if err != nil {
		panic(err)
	}
	defer dbsession.Close()

	jobs := dbhandler.ExternalEmploers(*dbsession)

	t.templ.Execute(w, jobs)
}

func main() {
	http.Handle("/update", &updateHandler{})
	http.Handle("/", &templateHandler{filename: "extjobs.html"})

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
