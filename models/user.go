package models

import (
	"time"
)

// User : model of all system users
type User struct {
	tableName struct{} `sql:"public.users"`
	ID        uint64
	Username  string ``
	Email     string ``
	UserType  uint   `doc:"1: authenticatable user, 2: company"`
	Firstname string ``
	Lastname  string ``
	Companies []User ``
	Avatar    string ``
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}
