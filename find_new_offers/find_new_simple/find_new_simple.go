package main

import (
	"flag"
	"fmt"
	"github.com/remotejob/gojobextractor/find_new_offers/find_new_simple/findalllinks"
	"github.com/remotejob/gojobextractor/find_new_offers/find_new_simple/jobdetails_simple"
	"gopkg.in/mgo.v2"
	"strconv"
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

	for i := 5; i < 6; i++ {

		navigstr := "http://stackoverflow.com/jobs?sort=p&pg=" + strconv.Itoa(i)
		links := findalllinks.FindAll(navigstr)

		for _, link := range links {

			newJobentry := jobdetails_simple.NewJobOffers()
			(*newJobentry).ParsePage(link)
			(*newJobentry).ExamDbRecord(*dbsession)

		}

	}

}
