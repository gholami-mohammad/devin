package models

import (
	"time"
)

type TaskReminder struct {
	ID                uint64
	TaskID            uint64
	Task              *Task
	Title             string
	RemindeOn         time.Time
	ReminderReceivers []*TaskReminderReceiver
	CreatedByID       uint64
	CreatedBy         *User
	CreatedAt         time.Time
	UpdatedAt         time.Time
}
