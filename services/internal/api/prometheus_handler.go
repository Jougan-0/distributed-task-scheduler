package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
)

func handlePrometheusQueryRange(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("query")
	start := r.URL.Query().Get("start")
	end := r.URL.Query().Get("end")
	step := r.URL.Query().Get("step")
	prom := os.Getenv("PROMETHEUS_URL")
	promURL := fmt.Sprintf("%s/api/v1/query_range", prom)
	promQueryParams := url.Values{}
	promQueryParams.Set("query", query)
	promQueryParams.Set("start", start)
	promQueryParams.Set("end", end)
	promQueryParams.Set("step", step)

	fullURL := promURL + "?" + promQueryParams.Encode()

	resp, err := http.Get(fullURL)
	if err != nil {
		log.Printf("Error querying Prometheus: %v", err)
		http.Error(w, "Error querying Prometheus", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	var promData interface{}
	if err := json.NewDecoder(resp.Body).Decode(&promData); err != nil {
		log.Printf("Error decoding Prometheus response: %v", err)
		http.Error(w, "Error decoding Prometheus response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(promData); err != nil {
		log.Printf("Error writing response: %v", err)
	}
}
