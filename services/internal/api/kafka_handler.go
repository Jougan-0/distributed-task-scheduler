package api

import (
	"encoding/json"
	"net/http"

	"Jougan-0/distributed-task-scheduler/internal/kafka"
)

func getKafkaEventsHandler(w http.ResponseWriter, r *http.Request) {
	events := kafka.GetEvents()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(events)
}
