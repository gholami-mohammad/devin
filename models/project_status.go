package models

import "time"

type ProjectStatus struct {
	ID        uint
	Status    string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}