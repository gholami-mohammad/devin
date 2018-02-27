package models

import "time"

type TaskBoard struct {
	tableName   struct{} `sql:"public.task_boards"`
	ID          uint64
	Name        string
	ProjectID   uint64
	Project     *Project
	Color       string
	CreatedByID uint64
	CreatedBy   *User
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time
}
