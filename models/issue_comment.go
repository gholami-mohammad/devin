package models

import "time"

type IssueComment struct {
	ID             uint64
	MilestoneID    uint64
	Milestone      *Issue
	ReplyToID      uint64
	ReplyTo        *IssueComment
	Comment        string
	AttachmentPath string
	CreatedByID    uint64
	CreatedBy      *User
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      *time.Time
}
