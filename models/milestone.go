package models

import (
	"time"
)

type Milestone struct {
	tableName        struct{} `sql:"public.milestones"`
	ID               uint64
	Name             string
	DueDate          time.Time                   `doc:"تاریخ دستیابی به هدف"`
	Description      string                      `doc:"Full description about the milestone"`
	ResponsibleUsers []*MilestoneResponsibleUser ``
	Followers        []*MilestoneFollower        ``
	Tags             []*TaggedObject             `doc:"A HasMany relation, where ModuleID = models.MODULE_MILESTONE"`
	Comments         []*MilestoneComment         `doc:"HasMany relation"`
	TaskLists        []*MilestoneTaskList        ``
	CreatedByID      uint64
	CreatedBy        *User
	CreatedAt        time.Time
	UpdatedAt        time.Time
	DeletedAt        *time.Time
}
