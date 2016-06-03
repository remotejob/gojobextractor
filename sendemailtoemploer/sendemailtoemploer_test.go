package main

import (
	"gopkg.in/mgo.v2"
	"testing"
	//    "gopkg.in/gcfg.v1"
	"github.com/remotejob/gojobextractor/domains"
	gm "github.com/onsi/gomega"
	"github.com/remotejob/gojobextractor/sendemailtoemploer/create_emails"
	"github.com/remotejob/gojobextractor/sendemailtoemploer/find_emploers_for_email"
	"github.com/remotejob/gojobextractor/sendemailtoemploer/send_emails"
	//    "fmt"
)

var results []domains.JobOffer
var emailstosend []domains.Email

func TestAll(t *testing.T) {
	gm.RegisterTestingT(t)
	dbsession, err := mgo.Dial("127.0.0.1")
	if err != nil {
		panic(err)
	}
	defer dbsession.Close()

	results = find_emploers_for_email.FindEmpl(*dbsession)
	//	gm.Expect(len(results) >0).Should(gm.BeTrue())
	emailstosend = create_emails.Create(results)
	if len(emailstosend) > 0 {
		send_emails.SendAll(*dbsession, emailstosend, glogin, gpass)
	}

}
