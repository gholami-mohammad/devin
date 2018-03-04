package models

import "time"

type IssueComment struct {
	tableName      struct{} `sql:"public.issue_comments"`
	ID             uint64
	IssueID        uint64
	Issue          *Issue
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
