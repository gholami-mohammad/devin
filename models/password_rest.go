package models

import "time"

type PasswordReset struct {
	ID           uint64
	UserID       uint64
	User         *User
	Token        string
	UsedForReset bool
	CreatedAt    time.Time
	UpdatedAt    time.Time
	ExpiresAt    time.Time
}

func (PasswordReset) TableName() string {
	return "public.password_resets"
}
