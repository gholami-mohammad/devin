package models

import (
	"time"
)

type UserProfile struct {
	tableName               struct{} `sql:"public.user_profile"`
	JobTitle                string
	OfficePhoneCountyCodeID uint            `doc:"FK to countries table"`
	OfficePhoneCountryCode  *Country        `doc:"Belogs to Country"`
	HomePhoneCountyCodeID   uint            `doc:"FK to countries table"`
	HomePhoneCountryCode    *Country        `doc:"Belogs to Country"`
	CellPhoneCountyCodeID   uint            `doc:"FK to countries table"`
	CellPhoneCountryCode    *Country        `doc:"Belogs to Country"`
	FaxCountyCodeID         uint            `doc:"FK to countries table"`
	FaxCountryCode          *Country        `doc:"Belogs to Country"`
	AddressCountryID        uint            `doc:"FK to countries table. To improve database performance and ignore inner joings on SQL queries to load this data."`
	Country                 *Country        `doc:"Belogs to Country"`
	ProvinceID              uint            `doc:"FK to provinces table. To improve database performance and ignore inner joings on SQL queries to load this data."`
	Province                *Province       `doc:"Belogs to Province"`
	CityID                  uint            `doc:"FK to cities table"`
	City                    *City           `doc:"Belogs to City"`
	Twitter                 string          `doc:"Twitter username e.g 'm6devin' or full profile URL like 'https://twitter.com/m6devin'"`
	Linkedin                string          `doc:"Linkedin full profile URL "`
	GooglePlus              string          `doc:"Google plus full profile URL"`
	Facebook                string          `doc:"Facebook username or full profile URL"`
	TelegramID              string          `doc:"Telegram username or full telegram profile URL"`
	Website                 string          `doc:"Personnal website URL"`
	LocalizationLanguageID  uint            `doc:"FK to countries table to get localization settings"`
	LocalizationLanguage    *Country        `doc:"Belongs to Country model to load i18n settings"`
	DateFormat              string          `doc:"Default date formate to show dates in UI. List of date formates stored in 'date_formats' table, but for more DB performance, directly saved here."`
	TimeFormat              string          `doc:"Default time format to show in UI. Time formats stored in 'time_formats' table, but for more DB performance, directly saved here."`
	CalendarSystemID        uint            `doc:"FK to calendar_systems"`
	CalendarSystem          *CalendarSystem `doc:"Which calendar system will used to use in datepicker and showing dates "`
	CreatedAt               time.Time
	UpdatedAt               time.Time
	DeletedAt               *time.Time
}
