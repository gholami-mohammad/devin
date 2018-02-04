package models

import "time"

type Wiki struct {
	ID           uint64
	Name         string
	ProjectID    uint64
	Project      *Project
	RepositoryID uint64
	Repository   *Repository
	Pages        []*WikiPage
	CreatedByID  uint64
	CreatedBy    *User
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    *time.Time
}
