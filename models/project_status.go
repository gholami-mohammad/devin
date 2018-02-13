package models

import "time"

//ProjectStatus : active, archived, pending, etc
type ProjectStatus struct {
	tableName struct{} `sql:"project_statuses"`
	ID        uint
	Status    string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}
