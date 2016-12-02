package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"time"

	gcfg "gopkg.in/gcfg.v1"
	mgo "gopkg.in/mgo.v2"

	"github.com/remotejob/gojobextractor/domains"
	"gitlab.com/remotejob/gojobextractor/dbhandler"
)

var addrs []string
var database string
var username string
var password string
var mechanism string
var emplouers []string

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

func main() {

	mongoDBDialInfo := &mgo.DialInfo{
		Addrs:     addrs,
		Timeout:   60 * time.Second,
		Database:  database,
		Username:  username,
		Password:  password,
		Mechanism: mechanism,
	}

	//	dbsession, err := mgo.Dial("127.0.0.1")
	dbsession, err := mgo.DialWithInfo(mongoDBDialInfo)
	if err != nil {
		panic(err)
	}
	defer dbsession.Close()

	f, err := os.Open("cant_apply.txt")

	if err != nil {

		log.Fatalln(err.Error())
	}
	// Create a new Scanner for the file.
	scanner := bufio.NewScanner(f)
	// Loop over all lines in the file and print them.
	for scanner.Scan() {
		line := scanner.Text()
		fmt.Println(line)

		emplouers = append(emplouers, line)
	}

	for _, emplId := range emplouers {

		dbhandler.UpdateEmployerById(*dbsession, emplId)

	}

}
