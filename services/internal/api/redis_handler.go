package api

import (
	"encoding/json"
	"log"
	"net/http"

	"Jougan-0/distributed-task-scheduler/internal/redis"
)

func getRedisKeysHandler(w http.ResponseWriter, r *http.Request) {
	keys, err := redis.Client.Keys(redis.Ctx, "*").Result()
	if err != nil {
		log.Printf("Error fetching Redis keys: %v", err)
		http.Error(w, "Failed to fetch Redis keys", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(keys)
}
