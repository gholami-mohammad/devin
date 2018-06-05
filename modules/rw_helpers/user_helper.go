package rw_helpers

import (
	"devin/helpers"
	"devin/models"
	"errors"
	"net/http"

	"github.com/jinzhu/gorm"

	user_repo "devin/modules/user/repository"
)

// UserExists check for existance of user with given ID
func UserExists(w http.ResponseWriter, db *gorm.DB, userID uint64) (e error) {
	var count uint64
	e = db.Model(&models.User{}).Where("id=?", userID).Count(&count).Error
	if e != nil {
		err := helpers.ErrorResponse{
			ErrorCode: http.StatusNotFound,
			Message:   "User not found",
		}
		helpers.NewErrorResponse(w, &err)
		return
	}

	if count == 0 {
		err := helpers.ErrorResponse{
			ErrorCode: http.StatusNotFound,
			Message:   "User not found",
		}
		e = errors.New(err.Message)
		helpers.NewErrorResponse(w, &err)
		return
	}

	return
}

// IsUserExists check existance of user in DB and handle http errors
func IsUserExists(w http.ResponseWriter, db *gorm.DB, userID uint64) (e error) {
	if user_repo.IsUserExists(db, userID) == true {
		return
	}

	err := helpers.ErrorResponse{
		Message:   "User Not Found",
		ErrorCode: http.StatusNotFound,
	}
	helpers.NewErrorResponse(w, &err)
	e = errors.New(err.Message)
	return
}
