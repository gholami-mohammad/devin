package models

import (
	"time"
)

type Task struct {
	tableName               struct{} `sql:"public.tasks"`
	ID                      uint64
	Title                   string
	OrderID                 uint `doc:"شماره ترتیب قرارگیری در لیست"`
	Description             string
	ScheduledStartDate      time.Time
	ScheduledCompletionDate time.Time
	StartDate               time.Time
	CompletionDate          time.Time `doc:"تاریخ تکمیل شدن تسک"`
	CompletedSuccessfully   bool      `doc:"true for 'success' on completion, false for 'fail' on completion"`
	Attachments             []*TaskAttachment
	PriorityID              uint
	Priority                *TaskPriority
	FontColor               string
	BackgroundColor         string
	Progress                int
	EstimatedTime           time.Duration
	Followers               []*TaskFollower
	PrerequisiteTasks       []*TaskPrerequisite
	Reminders               []*TaskReminder
	TaskBoardID             uint
	TaskBoard               *TaskBoard
	Tags                    []*TaggedObject `doc:"A HasMany relation, where ModuleID = models.MODULE_TASK"`
	SpentTimes              []*TaskSpentTime
	Comments                []*TaskComment
	Assignments             []*TaskAssignment
	CreatedByID             uint64
	CreatedBy               *User
	CreatedAt               time.Time
	UpdatedAt               time.Time
	DeletedAt               *time.Time
}
