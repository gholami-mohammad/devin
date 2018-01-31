package models

import (
	"time"
)

// MilestoneTag contains all tags assigned to a Milestone.
type MilestoneTag struct {
	ID          uint64
	TagID       uint64
	Tag         *Tag
	MilestoneID uint64
	Milestone   *Milestone
	CreatedByID uint64
	CreatedBy   *User
	CreatedAt   time.Time
}
