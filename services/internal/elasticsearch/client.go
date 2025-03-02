package elasticsearch

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
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

	res, err := ESClient.Index(index, req)
	if err != nil {
		log.Printf("Failed to index task: %v", err)
		return err
	}
	defer res.Body.Close()

	bodyBytes, _ := io.ReadAll(res.Body)
	log.Printf("Elasticsearch Response: %s", string(bodyBytes))

	if res.IsError() {
		return fmt.Errorf("elasticsearch error: %s", string(bodyBytes))
	}

	log.Println("Task indexed successfully")
	return nil
}

func SearchTasks(index, query string) ([]map[string]interface{}, error) {
	var buf bytes.Buffer

	queryBody := map[string]interface{}{
		"query": map[string]interface{}{
			"term": map[string]interface{}{
				"Name.keyword": query,
			},
		},
	}

	log.Printf("Elasticsearch Query: %+v", query)

	if err := json.NewEncoder(&buf).Encode(queryBody); err != nil {
		return nil, fmt.Errorf("error encoding query: %w", err)
	}

	res, err := ESClient.Search(
		ESClient.Search.WithContext(context.Background()),
		ESClient.Search.WithIndex(index),
		ESClient.Search.WithBody(&buf),
		ESClient.Search.WithPretty(),
	)
	if err != nil {
		return nil, fmt.Errorf("elasticsearch search error: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		bodyBytes, _ := io.ReadAll(res.Body)
		return nil, fmt.Errorf("search returned status %d: %s", res.StatusCode, string(bodyBytes))
	}

	var result map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("error parsing search response: %w", err)
	}

	rawHits, ok := result["hits"].(map[string]interface{})
	if !ok {
		return nil, nil
	}

	innerHits, ok := rawHits["hits"].([]interface{})
	if !ok {
		return nil, nil
	}

	var tasks []map[string]interface{}
	for _, hit := range innerHits {
		h, ok := hit.(map[string]interface{})
		if !ok {
			continue
		}

		source, ok := h["_source"].(map[string]interface{})
		if !ok {
			continue
		}

		tasks = append(tasks, source)
	}

	return tasks, nil
}
