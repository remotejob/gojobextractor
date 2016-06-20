package main

import (
	//	"flag"
	//	"fmt"
	"gopkg.in/gcfg.v1"	
	"github.com/remotejob/gojobextractor/domains"
	"github.com/remotejob/gojobextractor/find_new_offers/find_new_simple/findalllinks"
	"github.com/remotejob/gojobextractor/find_new_offers/find_new_simple/jobdetails_simple"
	"gopkg.in/mgo.v2"
	"log"
	"strconv"
	"time"
)

var addrs []string
var database string
var username string
var password string
var mechanism string
var startpage int
var stoppage int

func init() {

	var cfg domains.ServerConfig
	if err := gcfg.ReadFileInto(&cfg, "config.gcfg"); err != nil {
		log.Fatalln(err.Error())

	} else {

		addrs = cfg.Dbmgo.Addrs
		database = cfg.Dbmgo.Database
		username =	cfg.Dbmgo.Username
		password = 	cfg.Dbmgo.Password
		mechanism = cfg.Dbmgo.Mechanism
		startpage = cfg.Pages.Startpage
		stoppage = 	cfg.Pages.Stoppage							
	}

}

func main() {

	mongoDBDialInfo := &mgo.DialInfo{
	Addrs:    addrs,
	Timeout:  60 * time.Second,
	Database: database,
	Username: username,
	Password: password,
	Mechanism: mechanism, 	
}

//	dbsession, err := mgo.Dial("127.0.0.1")
	dbsession, err := mgo.DialWithInfo(mongoDBDialInfo)
	if err != nil {
		panic(err)
	}
	defer dbsession.Close()

	for i :=startpage ; i < stoppage; i++ {

		navigstr := "http://stackoverflow.com/jobs?sort=p&pg=" + strconv.Itoa(i)
		links := findalllinks.FindAll(navigstr)

		for _, link := range links {
			
			newJobentry := jobdetails_simple.NewJobOffers()
			(*newJobentry).ParsePage(link)
			(*newJobentry).ExamDbRecord(*dbsession)

		}

	}

}
