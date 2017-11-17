package models

import (
	"time"
)

/**
 * Valid task types
 * 1. Project
 * 2. Section OR module   (section = n * task + m * section)
 * 3. Mileston
 * 4. Task
 * 5. Sub task ( if task = n * sub => task is a section)
 */

// TaskType defines all valid task types
type TaskType struct {
	tableName struct{} `sql:"project_management.task_types"`
	ID        uint64
	Name      string
	CreatedAt time.Time
	UpdateAt  time.Time
	DeletedAt *time.Time
}
