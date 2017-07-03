package main

import (
	"log"

	"github.com/blevesearch/bleve"
)

func main() {

	// bleve, err := bleve.Open("employers.bleve")

	// if err == bleve.ErrorIndexPathDoesNotExist {

	// 	log.Println("db don't exists")

	// } else if err != nil {
	// 	log.Fatal(err)
	// } else {
	// 	log.Printf("Opening existing index...")
	// }
	// query := bleve.NewQueryStringQuery("bleve")
	// searchRequest := bleve.NewSearchRequest(query)
	// searchResult, _ := index.Search(searchRequest)

	index, _ := bleve.Open("employers.bleve")
	query := bleve.NewQueryStringQuery("golang")
	// searchRequest := bleve.NewSearchRequest(query)
	searchRequest := bleve.NewSearchRequestOptions(query, 1000, 0, false)

	// searchRequest.Fields = []string{"_id", "id", "company", "title", "location", "tags", "externallink", "email", "hits", "created_at", "applied", "description"}
	searchRequest.Fields = []string{"*"}
	searchResult, _ := index.Search(searchRequest)

	// log.Println(searchResult)
	// log.Println(searchResult.Hits[0].Fields["id"])

	for i, hit := range searchResult.Hits {

		log.Println(hit.Fields["id"])
		log.Println("--------------------------------------")

		log.Println(hit.Fields["description"])
		log.Println("--------------------------------------", i)

	}

}
