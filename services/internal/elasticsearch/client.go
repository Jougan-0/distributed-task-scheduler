package elasticsearch

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	elasticsearch "github.com/elastic/go-elasticsearch/v7"
)

var ESClient *elasticsearch.Client

func InitElasticsearch(esURL string) error {
	cfg := elasticsearch.Config{
		Addresses: []string{esURL},
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
		DisableMetaHeader: true,
	}
	client, err := elasticsearch.NewClient(cfg)
	if err != nil {
		return err
	}
	ESClient = client
	res, err := ESClient.Indices.Exists([]string{"tasks"})
	if err != nil {
		log.Printf("Failed to check if index exists: %v", err)
	} else if res.StatusCode == 404 {
		log.Println("Index 'tasks' not found. Creating it now...")
		_, err := ESClient.Indices.Create("tasks",
			ESClient.Indices.Create.WithBody(strings.NewReader(`{
				"settings": { "number_of_shards": 1, "number_of_replicas": 1 },
				"mappings": {
					"properties": {
						"Name": { "type": "text" },
						"Status": { "type": "keyword" },
						"Priority": { "type": "integer" },
						"CreatedAt": { "type": "date" }
					}
				}
			}`)),
		)
		if err != nil {
			log.Printf("Failed to create tasks index: %v", err)
		} else {
			log.Println("Tasks index created successfully")
		}
	}

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
	log.Print("Result of elasticsearch: ", result)
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
