package models

import (
	"time"
)

// TaskDiscussion :
type TaskDiscussion struct {
	tableName   struct{} `sql:"project_management.task_discussions"`
	ID          uint64
	Title       string
	TaskID      uint64
	Task        *Task
	CreatedByID uint64
	CreatedBy   *User
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time
	Chats       []*DiscussionChat
}
