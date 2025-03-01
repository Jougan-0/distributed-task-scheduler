package elasticsearch

import (
	"bytes"
	"context"
	"encoding/json"
	"log"

	elasticsearch "github.com/elastic/go-elasticsearch/v8"
)

var ESClient *elasticsearch.Client

func InitElasticsearch(esURL string) error {
	cfg := elasticsearch.Config{
		Addresses: []string{esURL},
	}
	client, err := elasticsearch.NewClient(cfg)
	if err != nil {
		return err
	}
	ESClient = client
	log.Println("Elasticsearch initialized successfully.")
	return nil
}

func IndexTask(index string, task interface{}) error {
	body, _ := json.Marshal(task)

	req := bytes.NewReader(body)
	_, err := ESClient.Index(index, req)
	if err != nil {
		log.Printf("Failed to index task: %v", err)
	}
	return err
}

// Search Tasks in Elasticsearch
func SearchTasks(index, query string) ([]map[string]interface{}, error) {
	var buf bytes.Buffer
	queryBody := map[string]interface{}{
		"query": map[string]interface{}{
			"match": map[string]string{
				"name": query,
			},
		},
	}
	json.NewEncoder(&buf).Encode(queryBody)

	res, err := ESClient.Search(
		ESClient.Search.WithContext(context.Background()),
		ESClient.Search.WithIndex(index),
		ESClient.Search.WithBody(&buf),
		ESClient.Search.WithPretty(),
	)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var result map[string]interface{}
	json.NewDecoder(res.Body).Decode(&result)

	hits := result["hits"].(map[string]interface{})["hits"].([]interface{})
	var tasks []map[string]interface{}
	for _, hit := range hits {
		source := hit.(map[string]interface{})["_source"].(map[string]interface{})
		tasks = append(tasks, source)
	}

	return tasks, nil
}
