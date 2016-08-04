package find_emploers_for_email

import (
	"fmt"
	"net/url"
	"strings"
	"testing"
	//    "gopkg.in/mgo.v2"
	//    gm "github.com/onsi/gomega"
	//    "github.com/remotejob/gojobextractor/domains"
)

//var results []domains.JobOffer

func TestUrlstr(t *testing.T) {

	urlstr := "mailto:meenan@microsoft.com?subject=Senior%20Software%20"

	u, err := url.Parse(urlstr)
	if err != nil {
		panic(err)
	}
	
	fmt.Println(urlstr,u.RawQuery)
	fmt.Println(strings.TrimRight(urlstr, u.RawQuery))
	
	
	
	cleanemail := strings.TrimLeft(strings.TrimRight(urlstr, u.RawQuery), "mailto")

	fmt.Println(cleanemail)	
	
	cleanemail = cleanemail[1 : len(cleanemail)-1]
	//			result.Email=cleanemail

	fmt.Println(cleanemail)
	//	gm.RegisterTestingT(t)
	//	dbsession, err := mgo.Dial("127.0.0.1")
	//	if err != nil {
	//		panic(err)
	//	}
	//	defer dbsession.Close()
	//
	//	results = FindEmpl(*dbsession)
	//	gm.Expect(len(results) >0).Should(gm.BeTrue())

}
