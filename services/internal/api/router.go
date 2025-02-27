package api

import (
	"net/http"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

func NewRouter(db *gorm.DB) *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/api/v1/tasks", createTaskHandler(db)).Methods("POST")
	r.HandleFunc("/api/v1/tasks", listTasksHandler(db)).Methods("GET")
	r.HandleFunc("/api/v1/tasks/{id}/status", updateTaskStatusHandler(db)).Methods("PATCH")
	r.HandleFunc("/api/v1/tasks/pending/count", getPendingTaskCountHandler(db)).Methods("GET")
	r.HandleFunc("/api/v1/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})
	// r.Handle("/metrics", MetricsHandler())
	return r
}
