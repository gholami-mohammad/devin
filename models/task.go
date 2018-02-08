package models

import (
	"time"
)

type Task struct {
	ID                      uint64
	Title                   string
	OrderID                 uint `doc:"شماره ترتیب قرارگیری در لیست"`
	Description             string
	ScheduledStartDate      time.Time
	ScheduledCompletionDate time.Time
	StartDate               time.Time
	CompletionDate          time.Time
	Attachments             []*TaskAttachment
	PriorityID              uint
	Priority                *TaskPriority
	FontColor               string
	BackgroundColor         string
	Progress                float32
	EstimatedTime           time.Duration
	Followers               []*User `pg:"many2many:task_followers"`
	PrerequisiteTasks       []*Task `pg:"many2many:task_prerequisites"`
	Reminders               []*TaskReminder
	TaskBoardID             uint
	TaskBoard               *TaskBoard
	Tags                    []*TaggedObject `doc:"A HasMany relation, where ModuleID = models.MODULE_TASK"`
	SpentTimes              []*TaskSpentTime
	Comments                []*TaskComment
	Assignes                []*TaskAssigne
	CreatedByID             uint64
	CreatedBy               *User
	CreatedAt               time.Time
	UpdatedAt               time.Time
	DeletedAt               *time.Time
}
