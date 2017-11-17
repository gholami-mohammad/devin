package models

import (
	"time"
)

// User : model of all system users
type User struct {
	tableName struct{} `sql:"public.users"`
	ID        uint64
	FName     string
	LName     string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}
