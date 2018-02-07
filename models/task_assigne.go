package models

import "time"

type TaskAssigne struct {
	ID          uint64
	UserID      uint64
	User        *User
	TaskID      uint64
	Task        *Task
	CreatedByID uint64
	CreatedBy   *User
	CreatedAt   time.Time
}
