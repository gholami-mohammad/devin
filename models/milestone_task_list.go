package models

import "time"

type MilestoneTaskList struct {
	ID          uint64
	MilestoneID uint64
	Milestone   *Milestone
	TaskListID  uint64
	TaskList    *TaskList
	CreatedByID uint64
	CreatedBy   *User
	CreatedAt   time.Time
}
