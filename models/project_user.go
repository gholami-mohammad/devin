package models

import (
	"time"
)

type ProjectUser struct {
	tableName   struct{} `sql:"peoject_users"`
	ID          uint64
	UserID      uint64 `doc:"ID of users record with type=1"`
	User        *User
	ProjectID   uint64 `doc:"ID of pm.peojects record"`
	Project     *Project
	CreatedByID uint64 `doc:"Who add this user to this project?"`
	CreatedBy   *User
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time
}
