package models

import (
	"time"
)

type Milestone struct {
	ID               uint64
	Name             string
	DueDate          time.Time           `doc:"تاریخ دستیابی به هدف"`
	Description      string              `doc:"Full description about the milestone"`
	ResponsibleUsers []*User             `pg:"many2many:milestone_responsibale_users"`
	Followers        []*User             `pg:"many2many:milestone_followers"`
	Tags             []*TaggedObject     `doc:"A HasMany relation, where ModuleID = models.MODULE_MILESTONE"`
	Comments         []*MilestoneComment `doc:"hasMany"`
	TaskLists        []*TaskList         `pg:"many2many:milestone_tasklists"`
	CreatedAt        time.Time
	UpdatedAt        time.Time
	DeletedAt        *time.Time
}
