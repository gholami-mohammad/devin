package models

import (
	"time"
)

type TaskListUser struct {
	tableName   struct{} `sql:"task_list_users"`
	ID          uint64   ``
	UserID      uint64   `doc:"ID of users record with type=1"`
	User        *User    ``
	TaskListID  uint64   `doc:"ID of tasklist record"`
	TaskList    *TaskList
	CreatedByID uint64
	CreatedBy   *User
	CreatedAt   time.Time
}
