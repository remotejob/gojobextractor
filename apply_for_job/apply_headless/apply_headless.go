package main

import (
	"fmt"
	"log"
	"time"

	"github.com/remotejob/gojobextractor/apply_for_job/handle_internal_link"

	"github.com/remotejob/gojobextractor/domains"
	"github.com/remotejob/gojobextractor/elasticLoader/loader/dbhandler"
	"github.com/tebeka/selenium"
	"gopkg.in/gcfg.v1"
	"gopkg.in/mgo.v2"
)

var login string
var pass string
var addrs []string
var database string
var username string
var password string
var mechanism string
var cvpdf string

func init() {

	var cfg domains.ServerConfig
	if err := gcfg.ReadFileInto(&cfg, "config.gcfg"); err != nil {
		log.Fatalln(err.Error())

	} else {

		login = cfg.Login.Slogin
		pass = cfg.Pass.Spass
		addrs = cfg.Dbmgo.Addrs
		database = cfg.Dbmgo.Database
		username = cfg.Dbmgo.Username
		password = cfg.Dbmgo.Password
		mechanism = cfg.Dbmgo.Mechanism
		// cvpdf = cfg.Cvpdf.File
		cvpdf = "/tmp/my_cv.pdf"
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

	results := dbhandler.FindNotApplyedEmployers(*dbsession)

	fmt.Println("implouers to apply", len(results))

	if len(results) > 0 {

		caps := selenium.Capabilities{"browserName": "chrome"}
		//				caps := selenium.Capabilities{"browserName": "phantomjs"}
		wd, err := selenium.NewRemote(caps, "")
		if err != nil {
			fmt.Println(err.Error())
		}
		defer wd.Quit()

		wd.Get("https://stackoverflow.com/users/login?ssrc=head&returnurl=http%3a%2f%2fstackoverflow.com%2fjobs")

		time.Sleep(time.Millisecond * 1500)

		elem, err := wd.FindElement(selenium.ByID, "email")
		if err != nil {
			fmt.Println(err.Error())
		}
		spass, err := wd.FindElement(selenium.ByID, "password")
		if err != nil {
			fmt.Println(err.Error())
		}
		time.Sleep(time.Millisecond * 1000)

		err = elem.SendKeys(login)
		if err != nil {
			fmt.Println(err.Error())
		}
		err = spass.SendKeys(pass)
		if err != nil {
			fmt.Println(err.Error())
		}
		btm, err := wd.FindElement(selenium.ByID, "submit-button")
		if err != nil {
			fmt.Println(err.Error())
		}
		btm.Click()
		time.Sleep(time.Millisecond * 4000)

		// for i := 0; i < len(results); i++ {
		for i := 0; i < 20; i++ {

			fmt.Println("id", results[i].Id)

			employer := handle_internal_link.NewInternalJobOffers(results[i])
			reCaph := (*employer).Apply_headless(*dbsession, wd, results[i].Id, cvpdf)

			if reCaph {

				log.Println("ReCaph Present Stop loop")
				break

			} else {
				log.Println("ReCaph NOT present Continue loop")
			}

		}

		wd.Quit()

	}

}
