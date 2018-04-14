package models

import (
	"time"
)

type Tag struct {
	tableName   struct{} `sql:"public.tags"`
	ID          uint64
	OrganizationID   uint64 `doc:"This tag created on this organization. All tags of a organization will be shared in all of its projects and modules"`
	Organization     *User
	Title       string
	CreatedByID uint64 `doc:"Who creates this tag?"`
	CreatedBy   *User
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time
}
