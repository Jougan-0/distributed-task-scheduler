package api

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
)

func PrometheusQueryHandler(w http.ResponseWriter, r *http.Request) {
	prometheusURL := os.Getenv("PROMETHEUS_URL")
	if prometheusURL == "" {
		prometheusURL = "http://localhost:9090"
	}

	queryParams := r.URL.Query()
	query := queryParams.Get("query")
	start := queryParams.Get("start")
	end := queryParams.Get("end")
	step := queryParams.Get("step")

	promURL := fmt.Sprintf("%s/api/v1/query_range?query=%s&start=%s&end=%s&step=%s",
		prometheusURL, url.QueryEscape(query), start, end, step)

	resp, err := http.Get(promURL)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to query Prometheus: %v", err), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Failed to read Prometheus response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(body)
}
