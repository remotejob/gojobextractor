package loader

import (
	"context"
	"fmt"

	"log"

	"github.com/remotejob/gojobextractor/dbhandler"
	mgo "gopkg.in/mgo.v2"
	elastic "gopkg.in/olivere/elastic.v5"
)

func Load(dbsession mgo.Session) {
	results := dbhandler.FindNotApplyedEmployers(dbsession)

	fmt.Println("implouers to apply", len(results))
	client, err := elastic.NewClient()
	if err != nil {
		// Handle error

		log.Fatalln(err.Error())
	}

	esversion, err := client.ElasticsearchVersion("http://127.0.0.1:9200")
	if err != nil {
		// Handle error
		panic(err)
	}
	fmt.Printf("Elasticsearch version %s", esversion)
	exists, err := client.IndexExists("twitter").Do(context.TODO())
	if err != nil {
		// Handle error
		panic(err)
	}

	if !exists {
		// Create a new index.
		createIndex, err := client.CreateIndex("twitter").Do(context.TODO())
		if err != nil {
			// Handle error
			panic(err)
		}
		if !createIndex.Acknowledged {
			// Not acknowledged

			log.Println("Not acknowledged")

		}
	}

	for i, result := range results {
		// log.Println(result)
		put1, err := client.Index().
			Index("employers").
			Type("employer").
			Id(results[i].Id).
			BodyJson(result).
			Do(context.TODO())
		if err != nil {
			// Handle error
			panic(err)
		}
		fmt.Printf("Indexed tweet %s to index %s, type %s\n", put1.Id, put1.Index, put1.Type)

	}

}
