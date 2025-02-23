package api

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"gorm.io/gorm"

	"Jougan-0/distributed-task-scheduler/internal/scheduler"
)

type CreateTaskRequest struct {
	Name          string `json:"name"`
	Type          string `json:"type"`
	Payload       string `json:"payload"`
	MaxRetries    uint   `json:"max_retries"`
	ScheduledTime string `json:"scheduled_time"`
	Priority      int    `json:"priority"`
}

func createTaskHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req CreateTaskRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		var scheduledTime time.Time
		if req.ScheduledTime != "" {
			parsedTime, err := time.Parse(time.RFC3339, req.ScheduledTime)
			if err != nil {
				http.Error(w, "Invalid scheduled_time format", http.StatusBadRequest)
				return
			}
			scheduledTime = parsedTime
		} else {
			scheduledTime = time.Now()
		}

		task := &scheduler.Task{
			Name:          req.Name,
			Type:          req.Type,
			Payload:       req.Payload,
			MaxRetries:    req.MaxRetries,
			Status:        scheduler.StatusPending,
			ScheduledTime: scheduledTime,
			Priority:      req.Priority,
		}

		if err := db.Create(task).Error; err != nil {
			log.Printf("Error creating task: %v", err)
			http.Error(w, "Failed to create task", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(task)
	}
}

func listTasksHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tasks, err := scheduler.ListTasks(db)
		if err != nil {
			log.Printf("Error listing tasks: %v", err)
			http.Error(w, "Failed to list tasks", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(tasks)
	}
}

type UpdateStatusRequest struct {
	Status scheduler.TaskStatus `json:"status"`
}

func updateTaskStatusHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		taskIDStr := vars["id"]

		taskID, err := strconv.Atoi(taskIDStr)
		if err != nil {
			http.Error(w, "Invalid task ID", http.StatusBadRequest)
			return
		}

		var req UpdateStatusRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		if err := scheduler.UpdateTaskStatus(db, uint(taskID), req.Status); err != nil {
			log.Printf("Error updating task status: %v", err)
			http.Error(w, "Failed to update task status", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Task status updated"))
	}
}
