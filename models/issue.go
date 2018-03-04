package models

import "time"

type Issue struct {
	tableName               struct{} `sql:"public.issues"`
	ID                      uint64
	ProjectID               uint64
	Project                 *Project
	RepositoryID            uint64 `doc:"nullable column"`
	Repository              *Repository
	Message                 string
	AttachmentFilePath      string
	Labels                  []*IssueLabel
	IssueAssignments        []*IssueAssignment
	StatusID                uint
	Status                  *IssueStatus
	SetAsInProgressDateTime *time.Time
	SetAsResolvedDateTime   *time.Time
	Comments                []*IssueComment
	CreatedByID             uint64
	CreatedBy               *User
	CreatedAt               time.Time
	UpdatedAt               time.Time
	DeletedAt               *time.Time
}
