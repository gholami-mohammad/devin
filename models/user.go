package models

import (
	"net/http"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"

	"devin/helpers"
)

// User : model of all system users
type User struct {
	tableName              struct{} `sql:"public.users"`
	ID                     uint64
	Username               string
	Email                  string
	Password               string `json:"-"`
	PlainPassword          string `json:"password" sql:"-"`
	UserType               uint   `json:"-" doc:"1: authenticatable user, 2: company"`
	FirstName              string
	LastName               string
	UserCompanyMapping     []*UserCompany `doc:"نگاشت کاربران عضو در هر کمپانی"`
	Avatar                 string
	OwnerID                uint64 `doc:"کد یکتای مالک و سازنده ی یک کمپانی. این فیلد برای حساب کاربری افراد میتواند خالی باشد."`
	Owner                  *User
	EmailVerified          bool
	EmailVerificationToken string `json:"-"`

	/**
	 * Profile properties
	 */
	JobTitle                 string          `doc:"User's job title in a company"`
	LocalizationLanguageID   uint            `doc:"FK to countries table to get localization settings"`
	LocalizationLanguage     *Country        `doc:"Belongs to Country model to load i18n settings"`
	DateFormat               string          `doc:"Default date formate to show dates in UI. List of date formates stored in 'date_formats' table, but for more DB performance, directly saved here."`
	TimeFormat               string          `doc:"Default time format to show in UI. Time formats stored in 'time_formats' table, but for more DB performance, directly saved here."`
	CalendarSystemID         uint            `doc:"FK to calendar_systems"`
	CalendarSystem           *CalendarSystem `sql:"-" doc:"Which calendar system will used to use in datepicker and showing dates "`
	OfficePhoneCountryCodeID uint            `doc:"FK to countries table"`
	OfficePhoneCountryCode   *Country        `doc:"Belogs to Country"`
	HomePhoneCountryCodeID   uint            `doc:"FK to countries table"`
	HomePhoneCountryCode     *Country        `doc:"Belogs to Country"`
	CellPhoneCountryCodeID   uint            `doc:"FK to countries table"`
	CellPhoneCountryCode     *Country        `doc:"Belogs to Country"`
	FaxCountryCodeID         uint            `doc:"FK to countries table"`
	FaxCountryCode           *Country        `doc:"Belogs to Country"`
	CountryID                uint            `doc:"#Address, FK to countries table. To improve database performance and ignore inner joings on SQL queries to load this data."`
	Country                  *Country        `doc:"Belogs to Country"`
	ProvinceID               uint            `doc:"#Address, FK to provinces table. To improve database performance and ignore inner joings on SQL queries to load this data."`
	Province                 *Province       `doc:"Belogs to Province"`
	CityID                   uint            `doc:"#Address, FK to cities table"`
	City                     *City           `doc:"Belogs to City"`
	Twitter                  string          `doc:"Twitter username e.g 'm6devin' or full profile URL like 'https://twitter.com/m6devin'"`
	Linkedin                 string          `doc:"Linkedin full profile URL "`
	GooglePlus               string          `doc:"Google plus full profile URL"`
	Facebook                 string          `doc:"Facebook username or full profile URL"`
	Telegram                 string          `doc:"Telegram username or full telegram profile URL"`
	Website                  string          `doc:"Personnal website URL"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `json:"-"`
}

// SetEncryptedPassword set new bcrypt password
func (user *User) SetEncryptedPassword(plainPassword string) {
	bts, _ := bcrypt.GenerateFromPassword([]byte(plainPassword), bcrypt.DefaultCost)

	user.Password = string(bts)
}

// SetNewEmailVerificationToken create new random string to verfy email address
func (user *User) SetNewEmailVerificationToken() {
	user.EmailVerified = false
	user.EmailVerificationToken = helpers.RandomString(54)
}

// CookieLifetime get the max time of Authorization cookie.
func (user User) CookieLifetime() time.Duration {
	return 13 * time.Hour
}

// TokenLifetime get the max time of Authorization token
func (user User) TokenLifetime() time.Duration {
	return 13 * time.Hour
}

// SetAuthorizationCookie set a cookie with `Authorization` name
func (user User) SetAuthorizationCookie(w http.ResponseWriter, value string) {
	cookie := &http.Cookie{}
	cookie.Name = "Authorization"
	cookie.Secure = true
	cookie.Value = value
	cookie.HttpOnly = true
	cookie.Expires = time.Now().Add(user.CookieLifetime())
	cookie.Path = "/"
	http.SetCookie(w, cookie)
}

// ExpireAuthorizationCookie exipre `Authorization` cookie if exists
func (user User) ExpireAuthorizationCookie(w http.ResponseWriter) {
	cookie := &http.Cookie{}
	cookie.Name = "Authorization"
	cookie.Secure = true
	cookie.Value = ""
	cookie.HttpOnly = true
	cookie.Expires = time.Now().Add(-10 * time.Hour)
	http.SetCookie(w, cookie)
}

// Token is toke structure of JWT
type Token struct {
	ID    uint64 `json:"user_id"`
	Token string `json:"Token"`
}

// Claim is claim structure of JWT
type Claim struct {
	jwt.StandardClaims

	ID       uint64 `json:"id"`
	Email    string `json:"email"`
	Username string `json:"username"`
}
