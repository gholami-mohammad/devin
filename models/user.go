package models

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	jwt_request "github.com/dgrijalva/jwt-go/request"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"

	"devin/crypto"
	"devin/database"
	"devin/helpers"
)

func Init() {
	log.SetFlags(log.Lshortfile)
}

// User : model of all system users
type User struct {
	tableName     struct{} `sql:"public.users"`
	ID            uint64
	Username      string
	Email         string
	Password      string `json:"-"`
	PlainPassword string `json:"Password" sql:"-"`

	// 1: authenticatable user, 2: organization
	UserType uint `json:"-"`

	// Handle preload from a user object
	UserOrganizationMapping []*UserOrganization `gorm:"ForeignKey:UserID"`

	// Handle preload from an organization object
	OrganizationUserMapping []*UserOrganization `gorm:"ForeignKey:OrganizationID"`

	// OwnerID used for users of type organization
	OwnerID                *uint64
	Owner                  *User
	EmailVerified          bool
	EmailVerificationToken *string `json:"-"`
	IsRootUser             bool    `json:"-"`
	CreatedAt              time.Time
	UpdatedAt              time.Time
	DeletedAt              *time.Time `json:"-"`

	PublicProfile
}

// PublicProfile store data about profile of user or organization
type PublicProfile struct {
	FirstName *string
	LastName  *string
	FullName  *string `sql:"-"`
	Avatar    *string
	JobTitle  *string
	// FK to countries table to get localization settings
	LocalizationLanguageID *uint
	LocalizationLanguage   *Country

	// Default date formate to show dates in UI. List of date formates stored in 'date_formats' table, but for more DB performance, directly saved here.
	DateFormat *string

	// Default time format to show in UI. Time formats stored in 'time_formats' table, but for more DB performance, directly saved here.
	TimeFormat *string

	// FK to calendar_systems
	// Which calendar system will used to use in datepicker and showing dates
	CalendarSystemID *uint
	CalendarSystem   *CalendarSystem

	// FK to countries table
	OfficePhoneCountryCodeID *uint
	OfficePhoneCountryCode   *Country
	OfficePhoneNumber        *string

	// FK to countries table
	HomePhoneCountryCodeID *uint
	HomePhoneCountryCode   *Country
	HomePhoneNumber        *string

	// FK to countries table
	CellPhoneCountryCodeID *uint
	CellPhoneCountryCode   *Country
	CellPhoneNumber        *string

	// FK to countries table
	FaxCountryCodeID *uint
	FaxCountryCode   *Country
	FaxNumber        *string

	// #Address, FK to countries table. To improve database performance and ignore inner joings on SQL queries to load this data.
	CountryID *uint
	Country   *Country

	// #Address, FK to provinces table. To improve database performance and ignore inner joings on SQL queries to load this data.
	ProvinceID *uint
	Province   *Province

	// #Address, FK to cities table
	CityID *uint
	City   *City

	// Twitter username e.g 'm6devin' or full profile URL like 'https://twitter.com/m6devin'
	Twitter *string

	// Linkedin full profile URL
	Linkedin *string

	// Google plus full profile URL
	GooglePlus *string

	// Facebook username or full profile URL
	Facebook *string

	// Telegram username or full telegram profile URL
	Telegram *string

	// Personnal website URL
	Website *string
}

func (User) TableName() string {
	return "public.users"
}

// SetEncryptedPassword set new bcrypt password
func (user *User) SetEncryptedPassword(plainPassword string) {
	bts, _ := bcrypt.GenerateFromPassword([]byte(plainPassword), bcrypt.DefaultCost)

	user.Password = string(bts)
}

// SetNewEmailVerificationToken create new random string to verfy email address
func (user *User) SetNewEmailVerificationToken() {
	user.EmailVerified = false
	rndStr := helpers.RandomString(54)
	user.EmailVerificationToken = &rndStr
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
	// cookie.Secure = true
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
	// cookie.Secure = true
	cookie.Value = ""
	cookie.HttpOnly = true
	cookie.Expires = time.Now().Add(-10 * time.Hour)
	http.SetCookie(w, cookie)
}

func (user User) ExtractUserFromRequestContext(r *http.Request) (User, *Claim, error) {
	var clm *Claim
	if r.Context().Value("Authorization") == nil {
		_, ok := r.Header["Authorization"]
		if !ok {
			cookie, e := r.Cookie("Authorization")
			if e != nil {
				return User{}, nil, errors.New("No Authorization context(header or cookie) found")
			}
			r.Header.Set("Authorization", cookie.Value)
		}

		token, err := jwt_request.ParseFromRequestWithClaims(r, jwt_request.HeaderExtractor{"Authorization"}, &Claim{}, func(token *jwt.Token) (interface{}, error) {
			return crypto.GetJWTVerifyKey()
		})
		if err != nil || !token.Valid {
			return User{}, nil, errors.New("Invalid authentication token")
		}

		clm = token.Claims.(*Claim)

	} else {
		clm = r.Context().Value("Authorization").(*Claim)
	}
	jsonString, e := crypto.CBCDecrypter(clm.Payload)
	if e != nil {
		return User{}, nil, e
	}
	var u, dbUser User
	json.Unmarshal([]byte(jsonString), &u)

	db := database.NewGORMInstance()
	defer db.Close()
	e = db.Where("id=?", u.ID).First(&dbUser).Error

	return dbUser, clm, nil
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

// IsUniqueValue check duplication of value in given column of users table.
// ignoredID use for ignore given ID of checking. Set ignoredID to 0 if you want to check all records.
func (user User) IsUniqueValue(db *gorm.DB, columnName string, value string, ignoredID uint64) (isUnique bool, e error) {
	var cnt struct {
		Cnt uint64
	}
	sql := `SELECT count(*) as cnt FROM users WHERE ` + columnName + `=? `
	if ignoredID != 0 {
		sql += " AND id != ?"
		e = db.Raw(sql, value, ignoredID).Scan(&cnt).Error
	} else {
		e = db.Raw(sql, value).Scan(&cnt).Error
	}

	if e != nil {
		return false, e
	}

	if cnt.Cnt != 0 {
		return false, nil
	}
	return true, nil
}

func (user *User) SetFullName() {
	if user.FirstName == nil && user.LastName == nil {
		return
	}
	fn := ""
	if user.FirstName != nil {
		fn = *user.FirstName
	}
	ln := ""
	if user.LastName != nil {
		ln = *user.LastName
	}
	full := strings.TrimSpace(fmt.Sprintf("%v %v", fn, ln))
	user.FullName = &full
}

// Claim is claim structure of JWT
type Claim struct {
	jwt.StandardClaims
	Payload string `json:"payload" doc:"Hex string encrypted with AES-256. Decrypted of this string contains id, username and email of user"`
}
