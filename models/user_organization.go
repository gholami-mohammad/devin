package models

import (
	"time"
)

// UserOrganization mapping of users and organization membership
type UserOrganization struct {
	ID uint64

	// ID of users record with type=1
	UserID *uint64
	User   *User

	// ID of users record with type=2
	OrganizationID *uint64
	Organization   *User

	// If IsAdminOfOrganization equals true , then user will has full access
	IsAdminOfOrganization bool

	// Permission of adding new project to assigned organization by this user
	CanCreateProject bool

	// Permission of editing a project of assigned organizaiton
	CanUpateProject bool

	// Permission of adding new user to organization by this user
	CanAddUserToOrganization bool

	// The creator of this record in the DB
	CreatedByID uint64
	CreatedBy   *User

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

// TableName return table name
func (UserOrganization) TableName() string {
	return "public.user_organization"
}
