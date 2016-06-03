package main

import (
	"github.com/remotejob/gojobextractor/dbhandler"
	"encoding/csv"
	"flag"
	"fmt"
	"gopkg.in/mgo.v2"
	"os"
	"sort"
)

const APP_VERSION = "0.1"

// The flag package provides a default help printer via -h switch
var versionFlag *bool = flag.Bool("v", false, "Print the version number.")

func main() {
	flag.Parse() // Scan the arguments list

	if *versionFlag {
		fmt.Println("Version:", APP_VERSION)
	}
	dbsession, err := mgo.Dial("127.0.0.1")
	if err != nil {
		panic(err)
	}
	defer dbsession.Close()

	//	csvfile, err := os.Open("/home/juno/git/jobprotractor/gojobextractor/mytags.csv")
	//	if err != nil {
	//		fmt.Println(err)
	//	}
	//	reader := csv.NewReader(csvfile)
	//	reader.LazyQuotes = true
	//
	//	records, err := reader.ReadAll()
	records := FromCSV("/home/juno/git/jobprotractor/gojobextractor/mytags.csv")

	mytags := make(map[string]struct{})

	for _, record := range records {

		mytags[record[0]] = struct{}{}

	}

	notmytagrecords := FromCSV("/home/juno/git/jobprotractor/gojobextractor/notmytags.csv")

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
//		fmt.Println(newtag)

		_, ok := notmytags[newtag]

		if !ok {

			fmt.Println(newtag)	
			//			newtags = append(newtags, k)

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
