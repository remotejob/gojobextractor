package create_emails

import (
	"github.com/remotejob/gojobextractor/apply_for_job/handle_internal_link/mytags"
	"github.com/remotejob/gojobextractor/domains"
	"encoding/csv"
	"fmt"
	"os"
)

func Create(emplayers []domains.JobOffer) []domains.Email {

	var emailstosend []domains.Email

	csvfile, err := os.Open("coverletter.csv")
	if err != nil {
		fmt.Println(err)
	}
	reader := csv.NewReader(csvfile)
	reader.LazyQuotes = true

	records, err := reader.ReadAll()

	for _, joboffer := range emplayers {

		fmt.Println("Send to:",joboffer.Id,joboffer.Externallink,joboffer.Email)
		mytagstoinsert := mytags.GetMyTags("mytags.csv", joboffer.Tags)

		body := "My experience:\n\n"

		for _, tag := range mytagstoinsert {

			body = body + tag.Tag + " -> " + tag.Duration + "\n"
		}

		for _, record := range records {

			body = body + "\n" + record[0]
		}

		body = body + "\n\nThanks.\nAlex Mazurov"
		body = body + "\n\nAtt:mazurov_cv.pdf"

		emailtxt := domains.Email{joboffer.Email,joboffer.Id, body}

		emailstosend = append(emailstosend, emailtxt)

	}

	return emailstosend
}
