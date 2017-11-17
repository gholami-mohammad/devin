package models

import (
	"time"
)

/**
 * Valid task statuses
 * 1. Not started
 * 2. In Progress
 * 3. Completed
 */

// TaskStatus shows all valid statuses of a task
type TaskStatus struct {
	tableName struct{} `sql:"project_management.task_statuses"`
	ID        uint
	Title     string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}
