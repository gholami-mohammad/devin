package models

import "time"

type TaskFollower struct {
	ID          uint64
	TaskID      uint64
	Task        *Task
	UserID      uint64
	User        *User
	CreatedByID uint64
	CreatedBy   *User
	CreatedAt   time.Time
}
