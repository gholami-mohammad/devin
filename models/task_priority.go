package models

import "time"

// TaskPriority : None , Low , Medium , High
type TaskPriority struct {
	tableName struct{} `sql:"public.task_priopities"`
	ID        uint
	Title     string
	Color     string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}
