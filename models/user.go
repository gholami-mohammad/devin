package models

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"

	"devin/crypto"
	"devin/helpers"
)

func Init() {
	log.SetFlags(log.Lshortfile)
}

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
	IsRootUser             bool

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
func (user User) SetAuthorizationCookieAndHeader(w http.ResponseWriter, value string) {
	cookie := &http.Cookie{}
	cookie.Name = "Authorization"
	cookie.Secure = true
	cookie.Value = value
	cookie.HttpOnly = true
	cookie.Expires = time.Now().Add(user.CookieLifetime())
	cookie.Path = "/"
	http.SetCookie(w, cookie)
	w.Header().Add("Authorization", value)
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

func (user User) ExtractUserFromRequestContext(r *http.Request) (User, *Claim, error) {
	clm := r.Context().Value("Authorization").(*Claim)
	jsonString, e := crypto.CBCDecrypter(clm.Payload)
	if e != nil {
		return User{}, nil, e
	}
	var u User
	json.Unmarshal([]byte(jsonString), &u)

	return u, clm, nil
}

func (user User) ExtractUserFromClaimPayload(payload string) (User, error) {

	jsonString, e := crypto.CBCDecrypter(payload)
	if e != nil {
		return User{}, e
	}
	var u User
	e = json.Unmarshal([]byte(jsonString), &u)
	return u, e
}

func (user User) GenerateNewTokenClaim() Claim {
	var claimPayload struct {
		ID       uint64 `json:"id"`
		Username string `json:"username"`
		Email    string `json:"email"`
	}
	claimPayload.ID = user.ID
	claimPayload.Username = user.Username
	claimPayload.Email = user.Email

	bts, _ := json.Marshal(&claimPayload)
	payload, _ := crypto.CBCEncrypter(string(bts))

	claim := Claim{
		Payload: payload,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(user.TokenLifetime()).Unix(),
			Issuer:    "devin",
		},
	}

	return claim
}

func (user User) GenerateNewTokenClaimWithCustomLifetime(duration time.Duration) Claim {
	var claimPayload struct {
		ID       uint64 `json:"id"`
		Username string `json:"username"`
		Email    string `json:"email"`
	}
	claimPayload.ID = user.ID
	claimPayload.Username = user.Username
	claimPayload.Email = user.Email
	bts, _ := json.Marshal(&claimPayload)

	payload, _ := crypto.CBCEncrypter(string(bts))
	claim := Claim{
		Payload: payload,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(duration).Unix(),
			Issuer:    "devin",
		},
	}

	return claim
}

func (user User) GenerateNewTokenString(claim Claim) (string, *helpers.ErrorResponse) {
	t := jwt.NewWithClaims(jwt.SigningMethodRS512, claim)

	sk, e := crypto.GetJWTSignKey()
	if e != nil {
		err := helpers.ErrorResponse{
			Message:   "Internal server error(load jwt)",
			ErrorCode: http.StatusInternalServerError,
		}

		return "", &err
	}
	tokenString, err := t.SignedString(sk)
	if err != nil {
		err := helpers.ErrorResponse{
			Message:   "Internal server error(sign jwt)",
			ErrorCode: http.StatusInternalServerError,
		}

		return "", &err
	}

	return tokenString, nil
}

// Claim is claim structure of JWT
type Claim struct {
	jwt.StandardClaims
	Payload string `json:"payload" doc:"Hex string encrypted with AES-256. Decrypted of this string contains id, username and email of user"`
}
