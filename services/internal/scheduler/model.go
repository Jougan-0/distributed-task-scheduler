package scheduler

import (
	"time"
)

type Task struct {
	ID            uint       `gorm:"primaryKey"`
	Name          string     `gorm:"type:varchar(255)"`
	Type          string     `gorm:"type:varchar(50)"`
	Payload       string     `gorm:"type:text"`
	Status        TaskStatus `gorm:"type:varchar(50)"`
	Attempts      uint       `gorm:"default:0"`
	MaxRetries    uint       `gorm:"default:3"`
	ScheduledTime time.Time  `gorm:"index"`
	Priority      int        `gorm:"default:5"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type TaskStatus string

const (
	StatusPending   TaskStatus = "PENDING"
	StatusRunning   TaskStatus = "RUNNING"
	StatusCompleted TaskStatus = "COMPLETED"
	StatusFailed    TaskStatus = "FAILED"
)
