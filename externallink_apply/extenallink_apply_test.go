package externallink_apply

import (
	gm "github.com/onsi/gomega"
	"gopkg.in/mgo.v2"
	"testing"
	"github.com/remotejob/gojobextractor/dbhandler"
	"fmt"
)

var dbsession *mgo.Session

func TestExterllink(t *testing.T) {

	dbsession, err := mgo.Dial("127.0.0.1")
	if err != nil {
		panic(err)
	}
	defer dbsession.Close()
	gm.RegisterTestingT(t)
	results := dbhandler.ExternalEmploers(*dbsession)
	
	if len(results) > 0 {
		
		for _,result := range results{
			
			fmt.Println(result.Externallink,result.Email)
		}
		
	}	
	

}
