package api

import (
	"encoding/json"
	"net/http"

	"Jougan-0/distributed-task-scheduler/internal/elasticsearch"

	"github.com/gorilla/mux"
)

func searchTasksHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		query := vars["query"]

		results, err := elasticsearch.SearchTasks("tasks", query)
		if err != nil {
			http.Error(w, "Failed to search tasks", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(results)
	}
}
