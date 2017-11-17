package models

import (
	"time"
)

// TaskAssignment show all user Assigned to a task
type TaskAssignment struct {
	tableName   struct{} `sql:"project_management.task_assignments"`
	ID          uint64
	TaskID      uint64
	Task        *Task
	UserID      uint64
	User        *User
	CreatedByID uint64
	CreatedBy   *User
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time
}
