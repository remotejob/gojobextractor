package loader

import (
	"log"
	"testing"
	"time"

	"github.com/remotejob/gojobextractor/domains"

	gcfg "gopkg.in/gcfg.v1"
	mgo "gopkg.in/mgo.v2"
)

var login string
var pass string
var addrs []string
var database string
var username string
var password string
var mechanism string
var cvpdf string
var mongoDBDialInfo mgo.DialInfo

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

		mongoDBDialInfo = mgo.DialInfo{
			Addrs:     addrs,
			Timeout:   60 * time.Second,
			Database:  database,
			Username:  username,
			Password:  password,
			Mechanism: mechanism,
		}
	}

}

func TestLoad(t *testing.T) {

	type args struct {
		dbsession mgo.Session
	}

	dbsession, err := mgo.DialWithInfo(&mongoDBDialInfo)
	if err != nil {
		log.Fatalln(err.Error())
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
		{"test0", args{*dbsession}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Load(tt.args.dbsession)
		})
	}
}
