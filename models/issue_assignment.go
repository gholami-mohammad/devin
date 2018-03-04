package models

import "time"

type IssueAssignment struct {
	tableName   struct{} `sql:"public.issue_assignments"`
	ID          uint64
	IssueID     uint64
	Issue       *Issue
	UserID      uint64
	User        *User
	CreatedByID uint64
	CreatedBy   *User
	CreatedAt   time.Time
}
