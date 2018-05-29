package controllers

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"

	"devin/database"
	"devin/helpers"
	"devin/models"
	"devin/modules/user/repository"
)

//RequestPasswordReset generate token to reset user's password
// @Mehtod: POST
// @Route: /api/password_reset/request?email={email}
func RequestPasswordReset(w http.ResponseWriter, r *http.Request) {
	email, e := extractEmailFromURL(w, r)
	if e != nil {
		return
	}
	db := database.NewGORMInstance()
	defer db.Close()

	user, e := getUserByEmail(w, db, email)
	if e != nil {
		return
	}

	var reset models.PasswordReset
	reset.UserID = user.ID
	reset.UsedForReset = false
	reset.Token = helpers.RandomString(26)
	reset.ExpiresAt = time.Now().Add(24 * time.Hour)

	db.Model(&models.PasswordReset{}).Create(&reset)

	helpers.NewSuccessResponse(w, "Password reset link sent to your email, please click to reset your new password!")
}

// extractEmailFromURL gt value of email address in query string
func extractEmailFromURL(w http.ResponseWriter, r *http.Request) (email string, e error) {
	email = r.URL.Query().Get("email")
	isValid := helpers.Validator{}.IsValidEmailFormat(email)
	if isValid == true {
		return
	}

	err := helpers.ErrorResponse{
		Message:   "Invalid email address!",
		ErrorCode: http.StatusUnprocessableEntity,
	}
	e = errors.New(err.Message)
	helpers.NewErrorResponse(w, &err)

	return

}

func getUserByEmail(w http.ResponseWriter, db *gorm.DB, email string) (user models.User, e error) {
	user, e = repository.GetUserByEmail(db, email)
	if e != nil {
		err := helpers.ErrorResponse{
			Message:   "No valid account found!",
			ErrorCode: http.StatusNotFound,
		}
		err.Errors = make(map[string][]string)
		err.Errors["dev"] = []string{e.Error()}
		e = errors.New(err.Message)
		helpers.NewErrorResponse(w, &err)

		return
	}

	return
}

//ValidatePasswordResetLink validate the token of password reset request
// @Method: GET
// @Route: /api/password_reset/validate/{token}
func ValidatePasswordResetLink(w http.ResponseWriter, r *http.Request) {

}

// extractTokenFromURL extract token parameter from the URL e.g /api/password_reset/{token}
func extractTokenFromURL(w http.ResponseWriter, r *http.Request) (token string, e error) {
	token = mux.Vars(r)["token"]
	if !strings.EqualFold(token, "") {
		return
	}

	err := helpers.ErrorResponse{
		Message:   "Invalid content type.",
		ErrorCode: http.StatusUnsupportedMediaType,
	}
	e = errors.New(err.Message)
	helpers.NewErrorResponse(w, &err)
	return

}
