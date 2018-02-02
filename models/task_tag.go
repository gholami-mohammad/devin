package models

import (
	"time"
)

// TaskTag contains all tags assigned to a task.
type TaskTag struct {
	ID          uint64
	TagID       uint64
	Tag         *Tag
	TaskID      uint64
	Task        *Task
	CreatedByID uint64
	CreatedBy   *User
	CreatedAt   time.Time
}
