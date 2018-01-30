package models

import (
	"time"
)

type UserCompanyMapping struct {
	tableName   struct{} `sql:"user_company_mappings"`
	ID          uint64
	UserID      uint64 `doc:"ID of users record with type=1"`
	User        *User
	CompanyID   uint64 `doc:"ID of users record with type=2"`
	Company     *User
	CreatedByID uint64 `doc:"Who add this user to this company? OR How create this company for this user?"`
	CreatedBy   *User
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time
}
