package find_emploers_for_email

import (
	"github.com/remotejob/gojobextractor/dbhandler"
	"github.com/remotejob/gojobextractor/domains"
	"fmt"
	"gopkg.in/mgo.v2"
	"net/url"
	"strings"
)

func FindEmpl(dbsession mgo.Session) []domains.JobOffer {

	results := dbhandler.ExternalEmploers(dbsession)

	var to_email_emploers []domains.JobOffer

	for _, result := range results {

		if strings.HasPrefix(result.Externallink, "mailto:") {

			u, err := url.Parse(result.Externallink)
			if err != nil {
				panic(err)
			}
			
			cleanemail := strings.TrimLeft(strings.TrimRight(result.Externallink,u.RawQuery),"mailto")

			cleanemail =cleanemail[1:len(cleanemail)-1]
			result.Email=cleanemail
			
			to_email_emploers = append(to_email_emploers, result)

		} else if result.Email != "" && !strings.HasPrefix(result.Externallink, "mailto:") {

			fmt.Print("get from mail record:",result.Email)
			to_email_emploers = append(to_email_emploers, result)

		}
	}

	return to_email_emploers
}
