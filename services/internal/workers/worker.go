package workers

import (
	"encoding/json"
	"log"
	"time"

	"Jougan-0/distributed-task-scheduler/internal/scheduler"

	"gorm.io/gorm"
)

// StartWorker runs the task processing loop.
func StartWorker(db *gorm.DB) {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		processPendingTasks(db)
	}

}

// processPendingTasks selects and processes tasks.
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

	err = scheduler.UpdateTaskStatus(db, task.ID, scheduler.StatusRunning)
	if err != nil {
		log.Printf("Worker: Failed to lock task ID=%d: %v", task.ID, err)
		return
	}

	log.Printf("Worker: Processing task ID=%d, Name=%s, Type=%s", task.ID, task.Name, task.Type)

	success := executeTask(task)

	if success {
		_ = scheduler.UpdateTaskStatus(db, task.ID, scheduler.StatusCompleted)
	} else {
		task.Attempts++

		if task.Attempts >= task.MaxRetries {
			_ = db.Model(&scheduler.Task{}).
				Where("id = ?", task.ID).
				Updates(map[string]interface{}{
					"status":   scheduler.StatusFailed,
					"attempts": task.Attempts,
				}).Error
			log.Printf("Worker: Task ID=%d has failed after %d retries.", task.ID, task.MaxRetries)
		} else {
			_ = db.Model(&scheduler.Task{}).
				Where("id = ?", task.ID).
				Updates(map[string]interface{}{
					"status":   scheduler.StatusPending,
					"attempts": task.Attempts,
				}).Error
			log.Printf("Worker: Retrying task ID=%d (Attempt %d of %d)", task.ID, task.Attempts, task.MaxRetries)
		}
	}
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
