package models

import (
	"time"
)

type Country struct {
	tableName   struct{} `sql:"public.address_countries"`
	ID          uint
	Name        string `doc:"The country name e.g Iran"`
	PhonePrefix string `doc:"Phone numbers prefix e.g +98 for IRAN"`
	Alpha2Code  string `doc:"ISO Alpha2 Code standard. e.g IR for IRAN"`
	Alpha3Code  string `doc:"ISO Alpha3 Code standard. e.g IRN for IRAN"`
	Flag        string `doc:"Base64 string of country flag"`
	LocaleCode  string `doc:"Localization (i18n) code e.g fa_IR for irannian persian language"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time `json:"-"`
}
