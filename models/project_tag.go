package models

import (
	"time"
)

// ProjectTag contains all tags assigned to a project.
type ProjectTag struct {
	ID          uint64
	TagID       uint64
	Tag         *Tag
	ProjectID   uint64
	Project     *Project
	CreatedByID uint64
	CreatedBy   *User
	CreatedAt   time.Time
}
