package models

import (
	"time"
)

// DiscussionChat :
type DiscussionChat struct {
	tableName    struct{} `sql:"project_management.discussion_chats"`
	ID           uint64
	DiscussionID uint64
	Discussion   *TaskDiscussion
	Message      string
	CreatedByID  uint64
	CreatedBy    *User
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    *time.Time
}
