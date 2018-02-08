package models

import "time"

//ProjectStatus : active, archived, pending, etc
type ProjectStatus struct {
	ID        uint
	Status    string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}
