package models

import (
	"time"
)

type TaskList struct {
	tableName    struct{} `sql:"public.task_lists"`
	ID           uint64
	Name         string
	Description  string
	AllowedUsers []*TaskListUser
	MilestoneID  uint64 `doc:"This task is belong to which milestone"`
	Milestone    *Milestone
	CreatedByID  uint64
	CreatedBy    *User
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    *time.Time
}
