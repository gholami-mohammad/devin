package models

import (
	"time"
)

// TaskReminder show all Reminders of a task
type TaskReminder struct {
	tableName   struct{} `sql:"project_management.task_reminders"`
	ID          uint64
	TaskID      uint64
	Task        *Task
	RemindOn    time.Time
	CreatedByID uint64
	CreatedBy   *User
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time
}
