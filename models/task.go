package models

import (
	"time"
)

// Task model
type Task struct {
	tableName               struct{} `sql:"project_management.tasks"`
	ID                      uint64
	Title                   string
	Description             string
	Duration                time.Duration
	ScheduledStartDate      time.Time
	ScheduledCompletionDate time.Time
	StartedOn               time.Time
	CompletedOn             time.Time
	CompletionPercentage    float32
	StatusID                uint
	Status                  *TaskStatus
	TypeID                  uint
	Type                    *TaskType
	Predecessors            []*TaskPredecessor
	Assignments             []*TaskAssignment
	Attachments             []*TaskAttachment
	Reminders               []*TaskReminder
	Discussions             []*TaskDiscussion
	CreatedByID             uint64
	CreatedBy               *User
	CreatedAt               time.Time
	UpdatedAt               time.Time
	DeletedAt               *time.Time
}
