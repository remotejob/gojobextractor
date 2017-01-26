package pdfcv

import (
	"fmt"
	"log"
	"testing"
	"time"

	gcfg "gopkg.in/gcfg.v1"
	mgo "gopkg.in/mgo.v2"

	"github.com/remotejob/gojobextractor/apply_for_job/handle_internal_link/mytags"
	"github.com/remotejob/gojobextractor/dbhandler"
	"github.com/remotejob/gojobextractor/domains"
)

var login string
var pass string
var addrs []string
var database string
var username string
var password string
var mechanism string
var cvpdf string
var joboffer domains.JobOffer
var results []domains.JobOffer
var _mytagstoinsert []domains.Tags

func init() {

	var cfg domains.ServerConfig
	if err := gcfg.ReadFileInto(&cfg, "/home/juno/gowork/src/github.com/remotejob/gojobextractor/config.gcfg"); err != nil {
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

	}

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

	results = dbhandler.FindNotApplyedEmployers(*dbsession)

	fmt.Println("implouers to apply", len(results))
	joboffer = results[10]
	_mytagstoinsert = mytags.GetMyTags("mytags.csv", joboffer.Tags)
}

func TestCreateCV(t *testing.T) {
	type args struct {
		emplayer       domains.JobOffer
		mytagstoinsert []domains.Tags
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
		{"test0", args{joboffer, _mytagstoinsert}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			CreateCV(tt.args.emplayer, tt.args.mytagstoinsert)
		})
	}
}
