package main

import (
	"github.com/remotejob/gojobextractor/domains"
	"github.com/remotejob/gojobextractor/sendemailtoemploer/create_emails"
	"github.com/remotejob/gojobextractor/sendemailtoemploer/find_emploers_for_email"
	"github.com/remotejob/gojobextractor/sendemailtoemploer/send_emails"
	//	"flag"
	//	"fmt"
	"gopkg.in/gcfg.v1"
	"gopkg.in/mgo.v2"	
	"log"
	//	"net/smtp"
	//	"os"
	//    "strconv"
)

var glogin = ""
var gpass = ""

func init() {

	var cfg domains.ServerConfig
	if err := gcfg.ReadFileInto(&cfg, "/home/juno/neonworkspace/gojobextractor/config.gcfg"); err != nil {
		log.Fatalln(err.Error())

	} else {

		glogin = cfg.Login.Glogin
		gpass = cfg.Pass.Gpass

	}

}

func main() {
	dbsession, err := mgo.Dial("127.0.0.1")
	if err != nil {
		panic(err)
	}
	defer dbsession.Close()

	results := find_emploers_for_email.FindEmpl(*dbsession)

	emailstosend := create_emails.Create(results)
	if len(emailstosend) > 0 {
		send_emails.SendAll(*dbsession, emailstosend, glogin, gpass)
	}

	//	send("hello there")

}


