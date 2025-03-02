package api

import (
	"encoding/json"
	"log"
	"net/http"

	"Jougan-0/distributed-task-scheduler/internal/elasticsearch"

	"github.com/gorilla/mux"
)

func searchTasksHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	query := vars["query"]

	log.Printf("Received search query: %s", query)
	results, err := elasticsearch.SearchTasks("tasks", query)
	if err != nil {
		log.Printf("Error searching tasks: %v", err)
		http.Error(w, "Failed to search tasks", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}
