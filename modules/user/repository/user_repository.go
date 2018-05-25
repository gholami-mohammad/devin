package repository

import (
	"devin/models"

	"github.com/jinzhu/gorm"
)

func IsUserExists(db *gorm.DB, userID uint64) bool {
	var count uint64
	db.Debug().Model(&models.User{}).Where("id=? AND user_type=1", userID).Count(&count)

	if count == 0 {
		return false
	}
	return true
}
