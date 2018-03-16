package helpers

import (
	"regexp"
)

type Validator struct{}

// ValidateEmailFormat check email address format.
func (Validator) IsValidEmailFormat(email string) bool {
	emailRegexp := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

	return emailRegexp.MatchString(email)
}

func (Validator) IsValidUsernameFormat(username string) bool {
	pattern := regexp.MustCompile(`^[a-z0-9\_]{3,100}$`)

	return pattern.MatchString(username)
}
