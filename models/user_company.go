package models

import (
	"time"
)

type UserCompany struct {
	tableName           struct{} `sql:"public.user_company"`
	ID                  uint64
	UserID              uint64 `doc:"ID of users record with type=1"`
	User                *User
	CompanyID           uint64 `doc:"ID of users record with type=2"`
	Company             *User
	IsAdminOfCompany    bool `doc:"If is_admin=true => user will has full access"`
	CanCreateProject    bool
	CanAddUserToCompany bool
	CreatedByID         uint64 `doc:"Who add this user to this company? OR How create this company for this user?"`
	CreatedBy           *User
	CreatedAt           time.Time
	UpdatedAt           time.Time
	DeletedAt           *time.Time
}
