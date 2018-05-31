package controllers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/jinzhu/gorm"

	"devin/database"
	"devin/helpers"
	"devin/models"
	"devin/modules/user/repository"
)

type passwordResetReqModel struct {
	Token          string
	Password       string
	PasswordVerify string
}

//RequestPasswordReset generate token to reset user's password
// @Mehtod: POST
// @Route: /api/password_reset/request
// @PostParams: email={email}
func RequestPasswordReset(w http.ResponseWriter, r *http.Request) {
	if helpers.IsRequestBodyNil(w, r) {
		return
	}

	email, e := extractEmailFromRequest(w, r)
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
	reset.Token = helpers.RandomString(64)
	reset.ExpiresAt = time.Now().Add(24 * time.Hour)

	db.Model(&models.PasswordReset{}).Create(&reset)

	helpers.NewSuccessResponse(w, "Password reset link sent to your email, please click to reset your new password!")
}

// extractEmailFromRequest gt value of email address in query string
func extractEmailFromRequest(w http.ResponseWriter, r *http.Request) (email string, e error) {
	var reqModel struct {
		Email string
	}

	json.NewDecoder(r.Body).Decode(&reqModel)
	defer r.Body.Close()
	email = reqModel.Email

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
// @Route: /api/password_reset/validate?token={token}
func ValidatePasswordResetLink(w http.ResponseWriter, r *http.Request) {
	token, e := extractTokenFromURLQS(w, r)
	if e != nil {
		return
	}
	db := database.NewGORMInstance()
	defer db.Close()

	_, e = getUserByResetPasswordToken(w, db, token)
	if e != nil {
		return
	}

	helpers.NewSuccessResponse(w, "Token is valid!")

	return
}

func getUserByResetPasswordToken(w http.ResponseWriter, db *gorm.DB, token string) (user models.User, e error) {
	var reset models.PasswordReset
	db.Model(&models.PasswordReset{}).
		Preload("User").
		Where("token=? AND used_for_reset=false", token).
		First(&reset)
	if reset.ID == 0 {
		err := helpers.ErrorResponse{}
		err.ErrorCode = http.StatusUnprocessableEntity
		err.Message = "Invalid password reset link!"
		e = errors.New(err.Message)
		helpers.NewErrorResponse(w, &err)
		return
	}

	if reset.ExpiresAt.Before(time.Now()) {
		err := helpers.ErrorResponse{}
		err.ErrorCode = http.StatusUnprocessableEntity
		err.Message = "Expired token! Request new one."
		e = errors.New(err.Message)
		helpers.NewErrorResponse(w, &err)
		return
	}

	if reset.User == nil {
		err := helpers.ErrorResponse{}
		err.ErrorCode = http.StatusNotFound
		err.Message = "No matching account!"
		e = errors.New(err.Message)
		helpers.NewErrorResponse(w, &err)
		return
	}

	user = *reset.User

	return
}

//ResetPassword reset assosiated user's password using reset token
func ResetPassword(w http.ResponseWriter, r *http.Request) {

	if helpers.IsRequestBodyNil(w, r) {
		return
	}

	reqModel, e := unmarshalPasswordResetRequest(w, r)
	if e != nil {
		return
	}

	if isValidResetPassword(w, reqModel) == false {
		return
	}

	db := database.NewGORMInstance()
	defer db.Close()

	user, e := getUserByResetPasswordToken(w, db, reqModel.Token)
	if e != nil {
		return
	}

	user.SetEncryptedPassword(reqModel.Password)
	db.Model(&user).UpdateColumn("password", user.Password)
	expireResetPasswordToken(db, reqModel.Token)
	helpers.NewSuccessResponse(w, "Your password updated successfully!")
	return

}

// unmarshalPasswordResetRequest decode json request to passwordResetReqModel object
func unmarshalPasswordResetRequest(w http.ResponseWriter, r *http.Request) (reqModel passwordResetReqModel, e error) {
	e = json.NewDecoder(r.Body).Decode(&reqModel)
	if e != nil {
		err := helpers.ErrorResponse{}
		err.ErrorCode = http.StatusBadRequest
		err.Message = "Invalid request!"
		err.Errors = make(map[string][]string)
		err.Errors["dev"] = []string{e.Error()}
		helpers.NewErrorResponse(w, &err)
		return
	}
	return
}

//isValidResetPassword validate new password strength
func isValidResetPassword(w http.ResponseWriter, reqModel passwordResetReqModel) bool {
	if strings.EqualFold(strings.Trim(reqModel.Password, " "), "") {
		err := helpers.ErrorResponse{}
		err.ErrorCode = http.StatusUnprocessableEntity
		err.Message = "An error accoured!"
		err.Errors = make(map[string][]string)
		err.Errors["Password"] = []string{"New Password can't be empty"}
		helpers.NewErrorResponse(w, &err)
		return false
	}

	if len(strings.Trim(reqModel.Password, " ")) < 6 {
		err := helpers.ErrorResponse{}
		err.ErrorCode = http.StatusUnprocessableEntity
		err.Message = "An error accoured!"
		err.Errors = make(map[string][]string)
		err.Errors["Password"] = []string{"Password length must be greater than 6 characters"}
		helpers.NewErrorResponse(w, &err)
		return false
	}

	if !strings.EqualFold(reqModel.Password, reqModel.PasswordVerify) {
		err := helpers.ErrorResponse{}
		err.ErrorCode = http.StatusUnprocessableEntity
		err.Message = "An error accoured!"
		err.Errors = make(map[string][]string)
		err.Errors["Password"] = []string{"New Password and its verify doesn't match!"}
		helpers.NewErrorResponse(w, &err)
		return false
	}

	return true
}

func expireResetPasswordToken(db *gorm.DB, token string) {
	db.Model(&models.PasswordReset{}).
		Where("token=?", token).
		UpdateColumns(models.PasswordReset{UsedForReset: true, UpdatedAt: time.Now()})
}
