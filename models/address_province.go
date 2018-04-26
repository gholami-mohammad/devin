package models

import (
	"time"
)

type Province struct {
	tableName struct{} `sql:"public.address_provinces"`
	ID        uint
	Name      string
	CountryID uint     `doc:"FK to countries table"`
	Country   *Country `doc:"Belongs to Country model(FK)"`
	Cities    []*City
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

func (Province) TableName() string {
	return "public.address_provinces"
}
