package client

import "github.com/olivere/elastic/v7"

func NewEsClient() (*elastic.Client, error) {
	client, err := elastic.NewClient(
		elastic.SetURL("http://0.0.0.0:9200"),
		elastic.SetSniff(false),
	)
	if err != nil {
		return nil, err
	}

	return client, nil
}
