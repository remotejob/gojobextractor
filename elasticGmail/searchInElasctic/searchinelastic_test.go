package searchInElasctic

import (
	"log"
	"testing"

	"gopkg.in/olivere/elastic.v5"
)

func TestSearch(t *testing.T) {
	client, err := elastic.NewClient()
	if err != nil {
		// Handle error

		log.Fatalln(err.Error())
	}
	type args struct {
		client *elastic.Client
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
		{"test0", args{client}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Search(tt.args.client)
		})
	}
}
