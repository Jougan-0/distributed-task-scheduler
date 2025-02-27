package main

import (
	"Jougan-0/distributed-task-scheduler/internal/api"
	"Jougan-0/distributed-task-scheduler/internal/config"
	"Jougan-0/distributed-task-scheduler/internal/db"
	"Jougan-0/distributed-task-scheduler/internal/kafka"
	"Jougan-0/distributed-task-scheduler/internal/redis"
	"Jougan-0/distributed-task-scheduler/internal/scheduler"
	"Jougan-0/distributed-task-scheduler/internal/workers"
	"context"

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

	if err := scheduler.AutoMigrate(dbConn); err != nil {
		log.Fatalf("DB Migration failed: %v", err)
	}

	if err := redis.Init(); err != nil {
		log.Fatalf("Failed to initialize Redis: %v", err)
	}
	log.Println("Redis initialized successfully.")

	if err := kafka.InitProducer(); err != nil {
		log.Fatalf("Failed to initialize Kafka: %v", err)
	}
	log.Println("Kafka Producer initialized successfully.")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() {
		if err := kafka.StartConsumerGroup(ctx, "my-consumer-group", "task-events"); err != nil {
			log.Printf("Kafka Consumer error: %v", err)
		}
	}()

	go workers.StartWorker(dbConn)
	server := api.NewServer(cfg, dbConn)
	log.Printf("Server running on port %s ...", cfg.ServerPort)
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
