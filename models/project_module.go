package models

import "time"

type ProjectModule struct {
	tableName struct{} `sql:"project_modules"`
	ID        uint64
	Name      string `doc:"Module name"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}
