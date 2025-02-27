package workers

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"Jougan-0/distributed-task-scheduler/internal/kafka"
	"Jougan-0/distributed-task-scheduler/internal/redis"
	"Jougan-0/distributed-task-scheduler/internal/scheduler"

	"gorm.io/gorm"
)

func StartWorker(db *gorm.DB) {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		processPendingTasks(db)
	}
}

func processPendingTasks(db *gorm.DB) {
	var task scheduler.Task

	err := db.Where("status = ? AND scheduled_time <= ?", scheduler.StatusPending, time.Now()).
		Order("priority ASC, scheduled_time ASC").
		First(&task).Error

	if err != nil {
		if err != gorm.ErrRecordNotFound {
			log.Printf("Worker: Error fetching task: %v", err)
		}
		return
	}

	if err := scheduler.UpdateTaskStatus(db, task.ID, scheduler.StatusRunning); err != nil {
		log.Printf("Worker: Failed to lock task ID=%s: %v", task.ID.String(), err)
		return
	}

	log.Printf("Worker: Processing task ID=%s, Name=%s, Type=%s", task.ID.String(), task.Name, task.Type)

	success := executeTask(task)

	if success {
		_ = scheduler.UpdateTaskStatus(db, task.ID, scheduler.StatusCompleted)
		log.Printf("Worker: Task ID=%s completed.", task.ID.String())
		kafka.PublishMessage("task-events", fmt.Sprintf("TaskCompleted:%s", task.ID.String()))
	} else {
		task.Attempts++
		if task.Attempts >= task.MaxRetries {
			_ = db.Model(&scheduler.Task{}).
				Where("id = ?", task.ID).
				Updates(map[string]interface{}{
					"status":   scheduler.StatusFailed,
					"attempts": task.Attempts,
				}).Error
			log.Printf("Worker: Task ID=%s has failed after %d retries.", task.ID.String(), task.MaxRetries)
			kafka.PublishMessage("task-events", fmt.Sprintf("TaskFailed:%s", task.ID.String()))
		} else {
			_ = db.Model(&scheduler.Task{}).
				Where("id = ?", task.ID).
				Updates(map[string]interface{}{
					"status":   scheduler.StatusPending,
					"attempts": task.Attempts,
				}).Error
			log.Printf("Worker: Retrying task ID=%s (Attempt %d of %d)", task.ID.String(), task.Attempts, task.MaxRetries)
		}
	}
	updatePendingTaskCount(db)
}

func executeTask(task scheduler.Task) bool {
	switch task.Type {
	case "EMAIL":
		return sendEmail(task.Payload)
	case "REPORT_GENERATION":
		return generateReport(task.Payload)
	default:
		log.Printf("Worker: Unknown task type: %s", task.Type)
		return false
	}
}

func sendEmail(payload string) bool {
	var data map[string]string
	json.Unmarshal([]byte(payload), &data)
	log.Printf("Worker: Sending email to %s", data["email"])
	time.Sleep(2 * time.Second)
	return true
}

func generateReport(payload string) bool {
	log.Println("Worker: Generating report with payload:", payload)
	time.Sleep(3 * time.Second)
	return true
}

func updatePendingTaskCount(db *gorm.DB) {
	var count int64
	if err := db.Model(&scheduler.Task{}).Where("status = ?", scheduler.StatusPending).Count(&count).Error; err != nil {
		log.Printf("Worker: Failed to update pending task count: %v", err)
		return
	}
	redis.Client.Set(redis.Ctx, "pending_tasks_count", count, 10*time.Second)
}
