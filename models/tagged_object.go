package models

import "time"

type TaggedObject struct {
	ID          uint64
	TagID       uint64
	Tag         *Tag
	ObjectID    uint64 `doc:"آی دی رکورد در جدول مربوط به ماژول"`
	ModuleID    uint   `doc:"لیست ماژول ها از ثابت های تعریف شده در همین پکیج تغذیه میشود"`
	CreatedByID uint64
	CreatedBy   *User
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
