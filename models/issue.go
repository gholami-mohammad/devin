package models

import "time"

type Issue struct {
	ID                  uint64
	ProjectID           uint64
	Project             *Project
	RepositoryID        uint64 `doc:"nullable column"`
	Repository          *Repository
	Message             string
	AttachementFilePath string
	Labels              []*IssueLabel
	AssineID            uint64
	Assigne             *User
	StatusID            uint
	Status              *IssueStatus
	Comments            []*IssueComment
	CreatedByID         uint64
	CreatedBy           *User
	CreatedAt           time.Time
	UpdatedAt           time.Time
	DeletedAt           *time.Time
}
