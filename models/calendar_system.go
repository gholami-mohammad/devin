package models

import (
	"time"
)

type CalendarSystem struct {
	tableName     struct{} `sql:"public.calendar_systems"`
	ID            uint64
	Name          string
	ComponentName string `doc:"UI component name to use for rendering UI date picker. e.g jalali-datepicker"`
	FilterName    string `doc:"UI filter name to use in template rendering. e.g jdate"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     *time.Time `json:"-"`
}
