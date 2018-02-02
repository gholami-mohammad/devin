package models

import "time"

// TaskPriority : None , Low , Medium , High
type TaskPriority struct {
	ID        uint
	Title     string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}
