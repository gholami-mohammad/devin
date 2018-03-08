package models

import (
	"time"
)

// User : model of all system users
type User struct {
	tableName          struct{} `sql:"public.users"`
	ID                 uint64
	Username           string
	Email              string
	UserType           uint `doc:"1: authenticatable user, 2: company"`
	Firstname          string
	Lastname           string
	UserCompanyMapping []*UserCompany `doc:"نگاشت کاربران عضو در هر کمپانی"`
	Avatar             string
	OwnerID            uint64 `doc:"کد یکتای مالک و سازنده ی یک کمپانی. این فیلد برای حساب کاربری افراد میتواند خالی باشد."`
	Owner              *User

	/**
	 * Profile properties
	 */
	JobTitle                string          `doc:"User's job title in a company"`
	LocalizationLanguageID  uint            `doc:"FK to countries table to get localization settings"`
	LocalizationLanguage    *Country        `doc:"Belongs to Country model to load i18n settings"`
	DateFormat              string          `doc:"Default date formate to show dates in UI. List of date formates stored in 'date_formats' table, but for more DB performance, directly saved here."`
	TimeFormat              string          `doc:"Default time format to show in UI. Time formats stored in 'time_formats' table, but for more DB performance, directly saved here."`
	CalendarSystemID        uint            `doc:"FK to calendar_systems"`
	CalendarSystem          *CalendarSystem `doc:"Which calendar system will used to use in datepicker and showing dates "`
	OfficePhoneCountyCodeID uint            `doc:"FK to countries table"`
	OfficePhoneCountryCode  *Country        `doc:"Belogs to Country"`
	HomePhoneCountyCodeID   uint            `doc:"FK to countries table"`
	HomePhoneCountryCode    *Country        `doc:"Belogs to Country"`
	CellPhoneCountyCodeID   uint            `doc:"FK to countries table"`
	CellPhoneCountryCode    *Country        `doc:"Belogs to Country"`
	FaxCountyCodeID         uint            `doc:"FK to countries table"`
	FaxCountryCode          *Country        `doc:"Belogs to Country"`
	CountryID               uint            `doc:"#Address, FK to countries table. To improve database performance and ignore inner joings on SQL queries to load this data."`
	Country                 *Country        `doc:"Belogs to Country"`
	ProvinceID              uint            `doc:"#Address, FK to provinces table. To improve database performance and ignore inner joings on SQL queries to load this data."`
	Province                *Province       `doc:"Belogs to Province"`
	CityID                  uint            `doc:"#Address, FK to cities table"`
	City                    *City           `doc:"Belogs to City"`
	Twitter                 string          `doc:"Twitter username e.g 'm6devin' or full profile URL like 'https://twitter.com/m6devin'"`
	Linkedin                string          `doc:"Linkedin full profile URL "`
	GooglePlus              string          `doc:"Google plus full profile URL"`
	Facebook                string          `doc:"Facebook username or full profile URL"`
	Telegram                string          `doc:"Telegram username or full telegram profile URL"`
	Website                 string          `doc:"Personnal website URL"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}
