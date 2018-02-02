package models

import "time"

type MilestoneComment struct {
	ID             uint64
	MilestoneID    uint64
	Milestone      *Milestone
	ReplyToID      uint64
	ReplyTo        *MilestoneComment
	Comment        string
	AttachmentPath string
	CreatedByID    uint64
	CreatedBy      *User
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      *time.Time
}
