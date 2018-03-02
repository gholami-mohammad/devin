package models

import "time"

type TaskPrerequisite struct {
	tableName      struct{} `sql:"public.task_prerequisites"`
	ID             uint64
	TaskID         uint64
	Task           *Task
	PrerequisiteID uint64
	Prerequisite   *Task
	CreatedByID    uint64
	CreatedBy      *User
	CreatedAt      time.Time
}
