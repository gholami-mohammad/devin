package models

import "time"

type TaskSpentTime struct {
	ID          uint64
	SpentByID   uint64 `doc:"automatically set to current user"`
	SpentBy     *User
	TaskID      uint64
	Task        *Task
	StartDate   time.Time
	EndDate     *time.Time
	IsBillable  bool
	IsBilled    bool
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time
}
