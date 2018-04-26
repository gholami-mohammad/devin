package models

import (
	"time"
)

type City struct {
	tableName  struct{} `sql:"public.address_cities"`
	ID         uint
	Name       string
	ProvinceID uint      `doc:"FK to provinces table"`
	Province   *Province `doc:"Belongs to Province model"`
	CountryID  uint      `doc:"FK to countries table for fastest DB data loading"`
	Country    *Country  `doc:"Belongs to Country model(FK)"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  *time.Time
}

func (City) TableName() string {
	return "public.address_cities"
}
