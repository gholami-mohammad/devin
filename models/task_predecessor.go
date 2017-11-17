package models

import (
	"time"
)

// TaskPredecessor contains preceding tasks for a task
type TaskPredecessor struct {
	tableName       struct{} `sql:"project_management.task_predecessors"`
	ID              uint64
	TaskID          uint64
	Task            *Task
	PrecedingTaskID uint64
	PrecedingTask   *Task
	CreatedByID     uint64
	CreatedBy       *User
	CreatedAt       time.Time
	UpdatedAt       time.Time
	DeletedAt       *time.Time
}
