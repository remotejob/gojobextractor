package jobdetails_simple

import (
//	"fmt"
	"gopkg.in/mgo.v2"
	"testing"
)

func TestParsePage(t *testing.T) {
	dbsession, err := mgo.Dial("127.0.0.1")
	if err != nil {
		panic(err)
	}
	defer dbsession.Close()

	newJobentry := NewJobOffers()
	(*newJobentry).ParsePage("http://stackoverflow.com/jobs/117144/java-developer-number26-gmbh?offset=8&pg=10&sort=p")
	(*newJobentry).ExamDbRecord(*dbsession)

//	fmt.Println(newJobentry)
}
