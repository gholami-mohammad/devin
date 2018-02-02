package models

import (
	"time"
)

type TaskSpentTimeTag struct {
	ID              uint64
	TagID           uint64
	Tag             *Tag
	TaskSpentTimeID uint64
	TaskSpentTime   *TaskSpentTime
	CreatedByID     uint64
	CreatedBy       *User
	CreatedAt       time.Time
}
