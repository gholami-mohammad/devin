package models

import (
	"time"
)

type TimeFormat struct {
	tableName struct{} `sql:"public.time_formates"`
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}
