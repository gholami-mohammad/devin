package models

import (
	"time"
)

type DateFormat struct {
	tableName struct{} `sql:"public.date_formats"`
	ID        uint
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `json:"-"`
}
