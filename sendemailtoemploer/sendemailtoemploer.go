package main

import (
	"github.com/remotejob/gojobextractor/domains"
	"github.com/remotejob/gojobextractor/sendemailtoemploer/create_emails"
	"github.com/remotejob/gojobextractor/sendemailtoemploer/find_emploers_for_email"
	"github.com/remotejob/gojobextractor/sendemailtoemploer/send_emails"
	"gopkg.in/gcfg.v1"
	"gopkg.in/mgo.v2"
	"log"
	"time"
)

var glogin string
var gpass string
var addrs []string
var database string
var username string
var password string
var mechanism string

func init() {

	var cfg domains.ServerConfig
	if err := gcfg.ReadFileInto(&cfg, "config.gcfg"); err != nil {
		log.Fatalln(err.Error())

	} else {

		glogin = cfg.Login.Glogin
		gpass = cfg.Pass.Gpass
		addrs = cfg.Dbmgo.Addrs
		database = cfg.Dbmgo.Database
		username = cfg.Dbmgo.Username
		password = cfg.Dbmgo.Password
		mechanism = cfg.Dbmgo.Mechanism

	}

}

func main() {
	
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

	results := find_emploers_for_email.FindEmpl(*dbsession)

	emailstosend := create_emails.Create(results)
	if len(emailstosend) > 0 {
		send_emails.SendAll(*dbsession, emailstosend, glogin, gpass)
	}


}
