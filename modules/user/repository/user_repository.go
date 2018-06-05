package repository

import (
	"devin/models"
	"errors"

	"github.com/jinzhu/gorm"
)

// IsUserExists check existance of userID in the users table
func IsUserExists(db *gorm.DB, userID uint64) bool {
	var count uint64
	db.Model(&models.User{}).Where("id=? AND user_type=1", userID).Count(&count)

	if count == 0 {
		return false
	}
	return true
}

// GetUserByEmail Load related user by given email address
func GetUserByEmail(db *gorm.DB, email string) (user models.User, e error) {
	db.Model(&models.User{}).Where("email=? AND user_type=1", email).First(&user)
	if user.ID == 0 {
		e = errors.New("User not found")
		return
	}

	return
}
