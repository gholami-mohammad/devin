package models

import "time"

type RepositoryContributer struct {
	ID           uint64
	RepositoryID uint64
	Repository   *Repository
	UserID       uint64
	User         *User
	RoleID       uint
	Role         *GitRole
	CreatedByID  uint64
	CreatedBy    *User
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    *time.Time
}
