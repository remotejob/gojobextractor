package loader

import (
	"context"
	"fmt"
	"strconv"

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
	// info, code, err := client.Ping("http://localhost:9200").Do()
	// if err != nil {
	// 	// Handle error
	// 	panic(err)
	// }
	// fmt.Printf("Elasticsearch returned with code %d and version %s", code, info.Version.Number)
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

	// _, err = client.CreateIndex("emploers").Do()
	// if err != nil {
	// 	// Handle error
	// 	panic(err)
	// }
	for i, result := range results {
		// log.Println(result)
		put1, err := client.Index().
			Index("twitter").
			Type("tweet").
			Id(strconv.Itoa(i)).
			BodyJson(result).
			Do(context.TODO())
		if err != nil {
			// Handle error
			panic(err)
		}
		fmt.Printf("Indexed tweet %s to index %s, type %s\n", put1.Id, put1.Index, put1.Type)

	}

}
