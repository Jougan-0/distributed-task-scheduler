package api

import (
	"encoding/json"
	"log"
	"net/http"

	"Jougan-0/distributed-task-scheduler/internal/redis"
)

type RedisKeyValue struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func getRedisKeysHandler(w http.ResponseWriter, r *http.Request) {
	keys, err := redis.Client.Keys(redis.Ctx, "*").Result()
	if err != nil {
		log.Printf("Error fetching Redis keys: %v", err)
		http.Error(w, "Failed to fetch Redis keys", http.StatusInternalServerError)
		return
	}

	var keyValues []RedisKeyValue

	for _, key := range keys {
		value, err := redis.Client.Get(redis.Ctx, key).Result()
		if err != nil {
			log.Printf("Error fetching value for key %s: %v", key, err)
			value = "Error fetching value"
		}
		keyValues = append(keyValues, RedisKeyValue{Key: key, Value: value})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(keyValues)
}
