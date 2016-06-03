package find_emploers_for_email

import (
    "testing"
    "gopkg.in/mgo.v2"
    gm "github.com/onsi/gomega"
    "github.com/remotejob/gojobextractor/domains" 
)

var results []domains.JobOffer

func TestFind(t *testing.T) {
	gm.RegisterTestingT(t)	
	dbsession, err := mgo.Dial("127.0.0.1")
	if err != nil {
		panic(err)
	}
	defer dbsession.Close()

	results = FindEmpl(*dbsession)
	gm.Expect(len(results) >0).Should(gm.BeTrue())


}

