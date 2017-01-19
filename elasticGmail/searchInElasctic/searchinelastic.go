package searchInElasctic

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/remotejob/gojobextractor/domains"

	"log"

	"gopkg.in/olivere/elastic.v5"
)

func Search(client *elastic.Client) {
	// get1, err := client.Get().
	// 	Index("employers").
	// 	Type("employer").
	// 	Id("1").Pretty(true).
	// 	Do(context.TODO())
	// if err != nil {
	// 	// Handle error
	// 	panic(err)
	// }
	// if get1.Found {

	// 	log.Println(get1)
	// 	fmt.Printf("Got document %s in version %d from index %s, type %s\n", get1.Id, get1.Version, get1.Index, get1.Type)

	// }
	_, err := client.Flush().Index("employers").Do(context.TODO())
	if err != nil {
		panic(err)
	}

	// termQuery := elastic.NewTermQuery("Hits", 1)

	searchResult, err := client.Search().
		Index("employers"). // search in index "twitter"
		// Query(termQuery).   // specify the query
		// Sort("id", true).  // sort by "user" field, ascending
		From(0).Size(10000). // take documents 0-9
		Pretty(true).        // pretty print request and response JSON
		Do(context.TODO())   // execute
	if err != nil {
		// Handle error
		panic(err)
	}

	fmt.Printf("Query took %d milliseconds\n", searchResult.TookInMillis)

	// var ttyp domains.JobOffer
	// for i, item := range searchResult.Each(reflect.TypeOf(ttyp)) {
	// 	if t, ok := item.(domains.JobOffer); ok {

	// 		log.Println(i, t.Id, t.Title)
	// 	}

	// }
	if searchResult.Hits.TotalHits > 0 {

		log.Println("res ", searchResult.Hits.TotalHits)

		for i, hit := range searchResult.Hits.Hits {
			// hit.Index contains the name of the index

			// Deserialize hit.Source into a Tweet (could also be just a map[string]interface{}).
			var t domains.JobOffer
			err := json.Unmarshal(*hit.Source, &t)
			if err != nil {
				// Deserialization failed
				log.Fatalln(err.Error())
			}

			// Work with tweet
			// fmt.Printf("Tweet by %s: %s\n", t.User, t.Message)
			log.Println(i, t.Id, t.Title)
		}
	} else {
		// No hits
		fmt.Print("Found no tweets\n")
	}
}
