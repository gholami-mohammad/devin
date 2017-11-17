package models

import (
	"time"
)

// TaskAttachment loads all Attachments of a task
type TaskAttachment struct {
	tableName   struct{} `sql:"project_management.task_attachments"`
	ID          uint64
	FilePath    string
	TaskID      uint64
	Task        *Task
	CreatedByID uint64
	CreatedBy   *User
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time
}
