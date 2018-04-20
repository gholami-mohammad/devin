package models

import (
	"time"
)

type UserOrganization struct {
	tableName                struct{} `sql:"public.user_organization"`
	ID                       uint64
	UserID                   *uint64 `doc:"ID of users record with type=1"`
	User                     *User
	OrganizationID           *uint64 `doc:"ID of users record with type=2"`
	Organization             *User
	IsAdminOfOrganization    bool   `doc:"If is_admin=true => user will has full access"`
	CanCreateProject         bool   `doc:"Permission of adding new project to assigned organization by this user"`
	CanAddUserToOrganization bool   `doc:"Permission of adding new user to organization by this user"`
	CreatedByID              uint64 `doc:"Who add this user to this organization? OR How create this organization for this user?"`
	CreatedBy                *User
	CreatedAt                time.Time
	UpdatedAt                time.Time
	DeletedAt                *time.Time
}
