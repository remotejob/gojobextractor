package main

import (
	"github.com/remotejob/gojobextractor/domains"
	//	"flag"
	//	"fmt"
	"gopkg.in/gcfg.v1"
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

//	send("hello there")

}
//func send(body string) {
//	from := glogin
//	pass := gpass
//	to := "aleksander.mazurov@gmail.com"
//
//	myfrom := "support@mazurov.eu"
//
//	msg := "From: " + myfrom + "\n" +
//		"To: " + to + "\n" +
//		"Subject: Hello there\n\n" +
//		body
//
//	err := smtp.SendMail("smtp.gmail.com:587",
//		smtp.PlainAuth("", from, pass, "smtp.gmail.com"),
//		from, []string{to}, []byte(msg))
//
//	if err != nil {
//		log.Printf("smtp error: %s", err)
//		return
//	}
//
//	log.Print("sent, visit http://foobarbazz.mailinator.com")
//}
