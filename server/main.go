package main

import (
	"github.com/remotejob/gojobextractor/dbhandler"
	"github.com/remotejob/gojobextractor/domains"
	"gopkg.in/gcfg.v1"
	"gopkg.in/mgo.v2"
	"log"
	"net/http"
	"path/filepath"
	"sync"
	"text/template"
	"time"
)

type templateHandler struct {
	once     sync.Once
	filename string
	templ    *template.Template
}

type updateHandler struct{}

var addrs []string
var database string
var username string
var password string
var mechanism string

//var mongoDBDialInfo mgo.DialInfo

//var dbsession mgo.Session

func init() {

	var cfg domains.ServerConfig
	if err := gcfg.ReadFileInto(&cfg, "config.gcfg"); err != nil {
		log.Fatalln(err.Error())

	} else {

		addrs = cfg.Dbmgo.Addrs
		database = cfg.Dbmgo.Database
		username = cfg.Dbmgo.Username
		password = cfg.Dbmgo.Password
		mechanism = cfg.Dbmgo.Mechanism

	}


}

func (t *updateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	if r.Method == "POST" {

		mongoDBDialInfo := &mgo.DialInfo{
			Addrs:     addrs,
			Timeout:   60 * time.Second,
			Database:  database,
			Username:  username,
			Password:  password,
			Mechanism: mechanism,
		}
		dbsession, err := mgo.DialWithInfo(mongoDBDialInfo)
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

	mongoDBDialInfo := &mgo.DialInfo{
		Addrs:     addrs,
		Timeout:   60 * time.Second,
		Database:  database,
		Username:  username,
		Password:  password,
		Mechanism: mechanism,
	}
	dbsession, err := mgo.DialWithInfo(mongoDBDialInfo)
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
