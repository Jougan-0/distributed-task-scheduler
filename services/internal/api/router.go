package api

import (
	"log"
	"net/http"

	socketio "github.com/googollee/go-socket.io"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/cors"
	"gorm.io/gorm"
)

var server *socketio.Server

func NewRouter(db *gorm.DB) http.Handler {

	r := mux.NewRouter()

	r.HandleFunc("/ws", WsHandler)

	r.HandleFunc("/api/v1/tasks", createTaskHandler(db)).Methods("POST")
	r.HandleFunc("/api/v1/tasks", listTasksHandler(db)).Methods("GET")
	r.HandleFunc("/api/v1/tasks/{id}/status", updateTaskStatusHandler(db)).Methods("PATCH")
	r.HandleFunc("/api/v1/tasks/pending/count", getPendingTaskCountHandler(db)).Methods("GET")
	r.HandleFunc("/api/v1/tasks/search/{query}", searchTasksHandler).Methods("GET")

	r.HandleFunc("/api/v1/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	r.Handle("/metrics", promhttp.Handler())

	r.HandleFunc("/kafka/events", getKafkaEventsHandler).Methods("GET")
	r.HandleFunc("/redis/keys", getRedisKeysHandler).Methods("GET")

	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: false,
	})

	return corsHandler.Handler(r)
}

func CloseSocketIOServer() {
	if server != nil {
		log.Println("Closing Socket.IO server...")
		server.Close()
	}
}
