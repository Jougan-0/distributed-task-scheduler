package main

import (
	"Jougan-0/distributed-task-scheduler/internal/api"
	"Jougan-0/distributed-task-scheduler/internal/config"
	"Jougan-0/distributed-task-scheduler/internal/db"

	"log"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	dbConn, err := db.InitDB(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}

	server := api.NewServer(cfg, dbConn)
	log.Printf("Server running on port %s ...", cfg.ServerPort)
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
