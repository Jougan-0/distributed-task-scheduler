package api

import (
	"fmt"
	"log"
	"net/http"

	"Jougan-0/distributed-task-scheduler/internal/config"

	"gorm.io/gorm"
)

type Server struct {
	httpServer *http.Server
}

func NewServer(cfg *config.Config, db *gorm.DB) *http.Server {
	router := NewRouter(db)
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.ServerPort),
		Handler: router,
	}

	return srv
}

func (s *Server) ListenAndServe() error {
	log.Printf("Starting HTTP server on %s", s.httpServer.Addr)
	return s.httpServer.ListenAndServe()
}
