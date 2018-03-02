package models

import (
	"time"
)

type TaskAttachment struct {
	tableName   struct{} `sql:"public.task_attachments"`
	ID          uint64
	FilePath    string
	TaskID      uint64
	Task        *Task
	CreatedByID uint64
	CreatedBy   *User
	CreatedAt   time.Time
}
