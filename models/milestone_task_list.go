package models

import "time"

type MilestoneTaskList struct {
	tableName   struct{} `sql:"public.milestone_task_lists"`
	ID          uint64
	MilestoneID uint64
	Milestone   *Milestone
	TaskListID  uint64
	TaskList    *TaskList
	CreatedByID uint64
	CreatedBy   *User
	CreatedAt   time.Time
}
