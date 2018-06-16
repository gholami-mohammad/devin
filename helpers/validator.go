package helpers

import (
	"regexp"
	"strings"
)

// Validator validators struct
type Validator struct{}

// IsValidEmailFormat check email address format.
func (Validator) IsValidEmailFormat(email string) bool {
	if &email == nil || strings.EqualFold(email, "") {
		return false
	}
	emailRegexp := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

	return emailRegexp.MatchString(email)
}

// IsValidUsernameFormat check the given string to use as username.
func (Validator) IsValidUsernameFormat(username string) bool {
	pattern := regexp.MustCompile(`^[a-z0-9\_]{3,100}$`)

	return pattern.MatchString(username)
}

// IsNilOrEmptyString check give string to be empty or nil
func IsNilOrEmptyString(str *string) bool {
	if str == nil || strings.EqualFold(strings.Trim(*str, " "), "") == true {
		return true
	}

	return false
}
