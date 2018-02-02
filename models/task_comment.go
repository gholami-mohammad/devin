package models

import "time"

type TaskComment struct {
	ID             uint64
	TaskID         uint64
	Task           *Task
	ReplyToID      uint64
	ReplyTo        *TaskComment
	Comment        string
	AttachmentPath string
	CreatedByID    uint64
	CreatedBy      *User
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      *time.Time
}
