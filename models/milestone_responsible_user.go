package models

import "time"

type MilestoneResponsibleUser struct {
	talbeName   struct{} `sql:"public.milestone_responsible_users"`
	ID          uint64
	MilestoneID uint64
	Milestone   *Milestone
	UserID      uint64
	User        *User
	CreatedByID uint64
	CreatedBy   *User
	CreatedAt   time.Time
}
