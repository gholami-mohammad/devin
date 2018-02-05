package models

import (
	"time"
)

// User : model of all system users
type User struct {
	tableName          struct{}       `sql:"public.users"`
	ID                 uint64         ``
	Username           string         ``
	Email              string         ``
	UserType           uint           `doc:"1: authenticatable user, 2: company"`
	Firstname          string         ``
	Lastname           string         ``
	UserCompanyMapping []*UserCompany `doc:"نگاشت کاربران عضو در هر کمپانی"`
	Avatar             string         ``
	OwnerID            uint64         `doc:"کد یکتای مالک و سازنده ی یک کمپانی. این فیلد برای حساب کاربری افراد میتواند خالی باشد."`
	Owner              *User
	CreatedAt          time.Time
	UpdatedAt          time.Time
	DeletedAt          *time.Time
}
