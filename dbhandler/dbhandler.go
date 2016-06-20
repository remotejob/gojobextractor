package dbhandler

import (
	"github.com/remotejob/gojobextractor/domains"
	"fmt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
)

func UpdateExtEmploerEmail(dbsession mgo.Session, email domains.Email) {

	dbsession.SetMode(mgo.Monotonic, true)

	c := dbsession.DB("cv_employers").C("employers")

	var joboffer domains.JobOffer

	err := c.Find(bson.M{"id": email.Subject}).Limit(1).One(&joboffer)
	if err != nil {

		log.Fatal(err)
	}

	joboffer.Applied = true

	err = c.Update(bson.M{"id": email.Subject}, joboffer)
	if err != nil {

		log.Fatal(err)
	}

}

func ExternalEmploers(dbsession mgo.Session) []domains.JobOffer {

	dbsession.SetMode(mgo.Monotonic, true)

	c := dbsession.DB("cv_employers").C("employers")

	var results []domains.JobOffer
//	err := c.Find(bson.M{"externallink": bson.M{"$ne": ""}, "location": bson.RegEx{Pattern: "Sweden", Options: "i"}, "applied": false}).All(&results)
	err := c.Find(bson.M{"externallink": bson.M{"$ne": ""},"applied": false}).All(&results)

	if err != nil {

		log.Fatal(err)
	}
	return results
}

func GetAllEmployers(dbsession mgo.Session) []domains.JobOffer {

	dbsession.SetMode(mgo.Monotonic, true)

	c := dbsession.DB("cv_employers").C("employers")
	var results []domains.JobOffer
	err := c.Find(nil).All(&results)
	if err != nil {

		log.Fatal(err)
	}

	return results
}

func InsertRecord(dbsession mgo.Session, joboffer domains.JobOffer) {

	dbsession.SetMode(mgo.Monotonic, true)

	c := dbsession.DB("cv_employers").C("employers")

	count, err := c.Find(bson.M{"id": joboffer.Id}).Limit(1).Count()
	if err != nil {

		log.Fatal(err)
	}

	if count == 0 {

		err := c.Insert(joboffer)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		fmt.Println("EXIST", count,joboffer.Id)

	}

}

func FindNotApplyedEmployers(dbsession mgo.Session) []domains.JobOffer {

	dbsession.SetMode(mgo.Monotonic, true)

	c := dbsession.DB("cv_employers").C("employers")

	var results []domains.JobOffer
	err := c.Find(bson.M{"externallink": "", "applied": false}).All(&results)
	if err != nil {

		log.Fatal(err)
	}

	return results

}

func UpdateEmployer(dbsession mgo.Session, joboffer domains.JobOffer) {

	dbsession.SetMode(mgo.Monotonic, true)

	c := dbsession.DB("cv_employers").C("employers")

	err := c.Update(bson.M{"id": joboffer.Id}, joboffer)
	if err != nil {

		log.Fatal(err)
	}

}

func UpdateEmployerById(dbsession mgo.Session, id string) {

	dbsession.SetMode(mgo.Monotonic, true)

	c := dbsession.DB("cv_employers").C("employers")

	err := c.Update(bson.M{"id": id}, bson.M{"$set": bson.M{"applied": true}})
	if err != nil {

		log.Fatal(err)
	}

}

