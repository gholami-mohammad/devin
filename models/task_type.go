package models

import (
	"time"
)

// TaskType defines all valid task types
type TaskType struct {
	tableName struct{} `sql:"project_management.task_types"`
	ID        uint64
	Name      string
	CreatedAt time.Time
	UpdateAt  time.Time
	DeletedAt *time.Time
}
