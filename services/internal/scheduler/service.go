package scheduler

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(&Task{})
}

func CreateTask(db *gorm.DB, t *Task) (*Task, error) {
	t.Status = StatusPending
	if err := db.Create(t).Error; err != nil {
		return nil, err
	}
	return t, nil
}

func ListTasks(db *gorm.DB) ([]Task, error) {
	var tasks []Task
	if err := db.Find(&tasks).Error; err != nil {
		return nil, err
	}
	return tasks, nil
}

func UpdateTaskStatus(db *gorm.DB, taskID uuid.UUID, newStatus TaskStatus) error {
	return db.Model(&Task{}).
		Where("id = ?", taskID).
		Update("status", newStatus).
		Error
}
