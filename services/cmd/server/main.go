package main

import (
	"Jougan-0/distributed-task-scheduler/internal/api"
	"Jougan-0/distributed-task-scheduler/internal/config"
	"Jougan-0/distributed-task-scheduler/internal/db"
	"Jougan-0/distributed-task-scheduler/internal/elasticsearch"
	"Jougan-0/distributed-task-scheduler/internal/kafka"
	"Jougan-0/distributed-task-scheduler/internal/metrics"
	"Jougan-0/distributed-task-scheduler/internal/redis"
	"Jougan-0/distributed-task-scheduler/internal/scheduler"
	"Jougan-0/distributed-task-scheduler/internal/workers"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
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

	metrics.RegisterMetrics()

	esURL := os.Getenv("ELASTICSEARCH_URL")
	if err := elasticsearch.InitElasticsearch(esURL); err != nil {
		log.Fatalf("Failed to initialize Elasticsearch: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() {
		if err := kafka.StartConsumerGroup(ctx, "my-consumer-group", "task-events"); err != nil {
			log.Printf("Kafka Consumer error: %v", err)
		}
	}()

	go workers.StartWorker(dbConn)

	router := api.NewRouter(dbConn)
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.ServerPort),
		Handler: router,
	}

	go func() {
		log.Printf("Server running on port %s ...", cfg.ServerPort)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("HTTP server ListenAndServe error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Println("Received shutdown signal. Shutting down gracefully...")

	ctxShutdown, cancelShutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelShutdown()

	if err := srv.Shutdown(ctxShutdown); err != nil {
		log.Printf("HTTP server shutdown error: %v", err)
	} else {
		log.Println("HTTP server shut down gracefully.")
	}

	api.CloseSocketIOServer()

	cancel()

	log.Println("Application shutdown complete.")
}
