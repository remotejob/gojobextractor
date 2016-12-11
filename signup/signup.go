package main

import (
	"fmt"
	"log"
	"time"

	"math/rand"

	"github.com/remotejob/gojobextractor/domains"
	"github.com/tebeka/selenium"
	gcfg "gopkg.in/gcfg.v1"
)

var login string
var pass string
var addrs []string
var database string
var username string
var password string
var mechanism string
var cvpdf string

var displayNames []string

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
		cvpdf = cfg.Cvpdf.File

		displayNames = []string{"Freelancer Development", "Programmer Development", "Programmer Freelancer", "Remote Development", "Support Development", "Programmer", "Free for Job"}

	}

}

func main() {

	rand.Seed(time.Now().UnixNano())

	randInt := rand.Perm(len(displayNames))
	displayName := displayNames[randInt[0]]

	caps := selenium.Capabilities{"browserName": "chrome"}
	//				caps := selenium.Capabilities{"browserName": "phantomjs"}
	wd, err := selenium.NewRemote(caps, "")
	if err != nil {
		fmt.Println(err.Error())
	}
	defer wd.Quit()

	wd.Get("https://stackoverflow.com/users/signup?ssrc=head&returnurl=%2fusers%2fstory%2fcurrent&utm_source=stackoverflow.com&utm_medium=dev-story&utm_campaign=signup-redirect")

	elem, err := wd.FindElement(selenium.ByID, "display-name")
	if err != nil {

		log.Fatalln(err.Error())
	}
	err = elem.SendKeys(displayName)
	if err != nil {
		fmt.Println(err.Error())
	}

	elemEmail, err := wd.FindElement(selenium.ByID, "email")
	if err != nil {

		log.Fatalln(err.Error())
	}
	err = elemEmail.SendKeys(login)
	if err != nil {
		fmt.Println(err.Error())
	}
	elemPass, err := wd.FindElement(selenium.ByID, "password")
	if err != nil {

		log.Fatalln(err.Error())
	}
	err = elemPass.SendKeys(pass)
	if err != nil {
		fmt.Println(err.Error())
	}

	fmt.Printf("Please enter an integer: ")

	// Read in an integer
	var i int
	_, err = fmt.Scanln(&i)
	if err != nil {
		fmt.Printf("Error: %s", err.Error())

		// If int read fails, read as string and forget
		var discard string
		fmt.Scanln(&discard)
		return
	}

	time.Sleep(time.Millisecond * 1500)
}
