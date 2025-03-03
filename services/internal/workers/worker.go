package workers

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"Jougan-0/distributed-task-scheduler/internal/api"
	"Jougan-0/distributed-task-scheduler/internal/kafka"
	"Jougan-0/distributed-task-scheduler/internal/metrics"
	"Jougan-0/distributed-task-scheduler/internal/redis"
	"Jougan-0/distributed-task-scheduler/internal/scheduler"

	"github.com/prometheus/client_golang/prometheus"
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
			api.BroadcastLog(fmt.Sprintf("Worker: Error fetching task: %v", err))
		}
		return
	}

	timer := prometheus.NewTimer(metrics.TaskProcessingTime)
	defer timer.ObserveDuration()

	api.BroadcastLog(fmt.Sprintf("Worker: Processing task ID=%s, Type=%s", task.ID.String(), task.Type))

	if err := scheduler.UpdateTaskStatus(db, task.ID, scheduler.StatusRunning); err != nil {
		log.Printf("Worker: Failed to lock task ID=%s: %v", task.ID.String(), err)
		api.BroadcastLog(fmt.Sprintf("Worker: Failed to lock task ID=%s: %v", task.ID.String(), err))
		return
	}

	log.Printf("Worker: Processing task ID=%s, Name=%s, Type=%s", task.ID.String(), task.Name, task.Type)

	success := executeTask(task)

	if success {
		_ = scheduler.UpdateTaskStatus(db, task.ID, scheduler.StatusCompleted)
		metrics.TotalTasksProcessed.WithLabelValues("completed").Inc()
		metrics.PendingTasksGauge.Dec()
		log.Printf("Worker: Task ID=%s completed.", task.ID.String())
		api.BroadcastLog(fmt.Sprintf("Worker: Task ID=%s completed.", task.ID.String()))
		eventMsg, _ := json.Marshal(kafka.KafkaEvent{
			Event:         "TaskCompleted",
			TaskID:        task.ID,
			TaskName:      task.Name,
			TaskType:      task.Type,
			Priority:      task.Priority,
			ScheduledTime: task.ScheduledTime.Format(time.RFC3339),
			CompletedAt:   time.Now().Format(time.RFC3339),
		})

		if err := kafka.PublishMessage("task-events", string(eventMsg)); err != nil {
			log.Printf("Failed to publish Kafka event for completed task %s: %v", task.ID.String(), err)
		}

	} else {
		task.Attempts++
		metrics.TaskRetries.WithLabelValues(task.Type).Inc()

		if task.Attempts >= task.MaxRetries {
			_ = db.Model(&scheduler.Task{}).
				Where("id = ?", task.ID).
				Updates(map[string]interface{}{
					"status":   scheduler.StatusFailed,
					"attempts": task.Attempts,
				}).Error
			metrics.TotalTasksProcessed.WithLabelValues("failed").Inc()
			metrics.PendingTasksGauge.Dec()
			log.Printf("Worker: Task ID=%s has failed after %d retries.", task.ID.String(), task.MaxRetries)
			api.BroadcastLog(fmt.Sprintf("Worker: Task ID=%s failed after %d retries.", task.ID.String(), task.MaxRetries))
			eventMsg, _ := json.Marshal(kafka.KafkaEvent{
				Event:         "TaskFailed",
				TaskID:        task.ID,
				TaskName:      task.Name,
				TaskType:      task.Type,
				Priority:      task.Priority,
				Attempts:      int(task.Attempts),
				MaxRetries:    int(task.MaxRetries),
				ScheduledTime: task.ScheduledTime.Format(time.RFC3339),
				FailedAt:      time.Now().Format(time.RFC3339),
			})

			if err := kafka.PublishMessage("task-events", string(eventMsg)); err != nil {
				log.Printf("Failed to publish Kafka event for failed task %s: %v", task.ID.String(), err)
			}

		} else {
			_ = db.Model(&scheduler.Task{}).
				Where("id = ?", task.ID).
				Updates(map[string]interface{}{
					"status":   scheduler.StatusPending,
					"attempts": task.Attempts,
				}).Error
			log.Printf("Worker: Retrying task ID=%s (Attempt %d of %d)", task.ID.String(), task.Attempts, task.MaxRetries)
			api.BroadcastLog(fmt.Sprintf("Worker: Retrying task ID=%s (Attempt %d of %d)", task.ID.String(), task.Attempts, task.MaxRetries))
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
	metrics.PendingTasksGauge.Set(float64(count))
}
