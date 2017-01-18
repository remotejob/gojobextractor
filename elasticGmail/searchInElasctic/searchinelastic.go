package searchInElasctic

import (
	"context"
	"fmt"
	"reflect"

	"github.com/remotejob/gojobextractor/domains"

	"log"

	"gopkg.in/olivere/elastic.v5"
)

func Search(client *elastic.Client) {
	get1, err := client.Get().
		Index("twitter").
		Type("tweet").
		Id("1").Pretty(true).
		Do(context.TODO())
	if err != nil {
		// Handle error
		panic(err)
	}
	if get1.Found {

		log.Println(get1)
		fmt.Printf("Got document %s in version %d from index %s, type %s\n", get1.Id, get1.Version, get1.Index, get1.Type)

	}
	_, err = client.Flush().Index("twitter").Do(context.TODO())
	if err != nil {
		panic(err)
	}

	termQuery := elastic.NewTermQuery("Hits", 0)

	searchResult, err := client.Search().
		Index("twitter"). // search in index "twitter"
		Query(termQuery). // specify the query
		// Sort("id", true).  // sort by "user" field, ascending
		From(0).Size(10).  // take documents 0-9
		Pretty(true).      // pretty print request and response JSON
		Do(context.TODO()) // execute
	if err != nil {
		// Handle error
		panic(err)
	}

	fmt.Printf("Query took %d milliseconds\n", searchResult.TookInMillis)

	log.Println("res ", searchResult.Hits.TotalHits)

	var ttyp domains.JobOffer
	for _, item := range searchResult.Each(reflect.TypeOf(ttyp)) {
		if t, ok := item.(domains.JobOffer); ok {
			fmt.Printf("Tweet by %s: %s\n", t.Id, t.Title)
		}
	}

}
