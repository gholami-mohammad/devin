package models

import "time"

type TaskFollower struct {
	tableName   struct{} `sql:"public.task_followers"`
	ID          uint64
	TaskID      uint64
	Task        *Task
	UserID      uint64
	User        *User
	CreatedByID uint64
	CreatedBy   *User
	CreatedAt   time.Time
}
