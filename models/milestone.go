package models

import (
	"time"
)

type Milestone struct {
	ID               uint64
	Name             string
	DueDate          time.Time                   `doc:"تاریخ دستیابی به هدف"`
	Description      string                      `doc:"Full description about the milestone"`
	ResponsibleUsers []*MilestoneResponsibleUser ``
	Followers        []*MilestoneFollower        `pg:"many2many:milestone_followers"`
	Tags             []*TaggedObject             `doc:"A HasMany relation, where ModuleID = models.MODULE_MILESTONE"`
	Comments         []*MilestoneComment         `doc:"HasMany relation"`
	TaskLists        []*TaskList                 `pg:"many2many:milestone_tasklists"`
	CreatedByID      uint64
	CreatedBy        *User
	CreatedAt        time.Time
	UpdatedAt        time.Time
	DeletedAt        *time.Time
}
