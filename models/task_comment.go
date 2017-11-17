package models

import (
	"time"
)

// TaskComment shows a comment on a task
type TaskComment struct {
	tableName   struct{} `sql:"project_management.task_comments"`
	ID          uint64
	Comment     string
	TaskID      uint64
	Task        *Task
	CreatedByID uint64
	CreatedBy   *User
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time
}
