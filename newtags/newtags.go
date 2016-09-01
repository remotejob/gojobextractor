package main

import (
	"encoding/csv"

	"github.com/remotejob/gojobextractor/dbhandler"
	"github.com/remotejob/gojobextractor/domains"
	//	"flag"
	"fmt"
	"log"
	"os"
	"sort"
	"time"

	"gopkg.in/gcfg.v1"
	"gopkg.in/mgo.v2"
)

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

	records := FromCSV("mytags.csv")

	mytags := make(map[string]struct{})

	for _, record := range records {

		mytags[record[0]] = struct{}{}

	}

	notmytagrecords := FromCSV("notmytags.csv")

	notmytags := make(map[string]struct{})

	for _, record := range notmytagrecords {

		notmytags[record[0]] = struct{}{}

	}

	var keys []string
	for k := range mytags {
		keys = append(keys, k)
	}

	employers := dbhandler.GetAllEmployers(*dbsession)

	employerstags := make(map[string]struct{})

	for _, employer := range employers {

		for _, tag := range employer.Tags {

			employerstags[tag] = struct{}{}
		}

	}

	var newtags []string

	for k, _ := range employerstags {

		_, ok := mytags[k]
		if !ok {

			newtags = append(newtags, k)

		}
	}

	sort.Strings(newtags)

	for _, newtag := range newtags {

		_, ok := notmytags[newtag]

		if !ok {

			fmt.Println(newtag)

		}

	}

}

func FromCSV(file string) [][]string {
	csvfile, err := os.Open(file)
	if err != nil {
		fmt.Println(err)
	}
	reader := csv.NewReader(csvfile)
	reader.LazyQuotes = true

	records, err := reader.ReadAll()
	if err != nil {
		fmt.Println(err)
	}

	return records

}
