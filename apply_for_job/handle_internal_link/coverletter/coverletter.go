package coverletter

import (
	"github.com/remotejob/gojobextractor/domains"
	"encoding/csv"
	"fmt"
	"os"
)

var coverlettertxt string

func Create(tags []domains.Tags, coverlettercsv string) string {

	if len(tags) > 0 {
		coverlettertxt = "My experience:\n\n"

		for _, tag := range tags {

			coverlettertxt = coverlettertxt + tag.Tag + " -> " + tag.Duration + "\n"
		}

	} else {
		coverlettertxt = "Only for your consideration:\n\n"

	}

	csvfile, err := os.Open(coverlettercsv)
	if err != nil {
		fmt.Println(err)
	}
	reader := csv.NewReader(csvfile)
	reader.LazyQuotes = true

	records, err := reader.ReadAll()

	for _, record := range records {

		coverlettertxt = coverlettertxt + "\n" + record[0]
	}

	return coverlettertxt
}
