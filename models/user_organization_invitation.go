package models

import "time"

type UserOrganizationInvitation struct {
	ID                 uint64
	UserID             *uint64
	Email              *string
	OrganizationID     uint64
	Organization       *User
	Accepted           *bool
	AcceptedRejectedAt *time.Time
	CreatedByID        uint64
	CreatedAt          time.Time
	UpdatedAt          time.Time
}
