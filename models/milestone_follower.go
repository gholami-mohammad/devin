package models

import "time"

type MilestoneFollower struct {
	ID          uint64
	MilestoneID uint64
	Milestone   *Milestone
	UserID      uint64
	User        *User
	CreatedByID uint64
	CreatedBy   *User
	CreatedAt   time.Time
}
