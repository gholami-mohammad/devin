package controllers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"

	"devin/database"
	"devin/helpers"
	"devin/models"

	"github.com/jinzhu/gorm"
)

func init() {
	log.SetFlags(log.Lshortfile)
}

// Signup handle user registration in application.
// Method: POST
// Content-Type: josn/application
//
// برای ثبت نام در سامانه اطلاعات ایمیل- نام کاربری-رمز عبور الزامی است
// تاییدیه ایمیل ارسال می شود
func Signup(w http.ResponseWriter, r *http.Request) {

	// Check content type
	if !helpers.HasJSONRequest(r) {
		err := helpers.ErrorResponse{
			Message:   "Invalid content type.",
			ErrorCode: http.StatusUnsupportedMediaType,
		}
		helpers.NewErrorResponse(w, &err)
		return
	}

	var user models.User
	e := json.NewDecoder(r.Body).Decode(&user)
	if e != nil {
		err := helpers.ErrorResponse{
			Message:   "Bad request body",
			ErrorCode: http.StatusBadRequest,
		}
		helpers.NewErrorResponse(w, &err)
		return
	}

	// lowercase username and email
	user.Username = strings.ToLower(user.Username)
	user.Email = strings.ToLower(user.Email)

	e, errorMessages := validateSignupInputs(user)
	if e != nil {
		err := helpers.ErrorResponse{
			Message:   e.Error(),
			ErrorCode: http.StatusUnprocessableEntity,
			Errors:    errorMessages,
		}
		helpers.NewErrorResponse(w, &err)
		return
	}

	db := database.NewGORMInstance()
	defer db.Close()
	// Check for duplication of email
	is, _ := user.IsUniqueValue(db, "email", user.Email, 0)
	if is == false {
		messages := make(map[string][]string)
		messages["Email"] = []string{"This email is already registered."}
		err := helpers.ErrorResponse{
			Message:   "Invalid Email address.",
			ErrorCode: http.StatusUnprocessableEntity,
			Errors:    messages,
		}
		helpers.NewErrorResponse(w, &err)
		return

	}
	// Check for duplication of username
	is, _ = user.IsUniqueValue(db, "username", user.Username, 0)
	if is == false {
		messages := make(map[string][]string)
		messages["Username"] = []string{"This username is already registered."}
		err := helpers.ErrorResponse{
			Message:   "Invalid username.",
			ErrorCode: http.StatusUnprocessableEntity,
			Errors:    messages,
		}
		helpers.NewErrorResponse(w, &err)
		return

	}

	user.UserType = 1
	user.SetEncryptedPassword(user.PlainPassword)
	user.SetNewEmailVerificationToken()

	// Saving data
	e = db.Save(&user).Error

	if e != nil {
		err := helpers.ErrorResponse{
			Message:   "Internal server error.",
			ErrorCode: http.StatusInternalServerError,
		}
		helpers.NewErrorResponse(w, &err)
		return
	}

	//omit password to be included in response
	user.PlainPassword = ""
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&user)
}

// validateSignupInputs check data validatin on user registration
func validateSignupInputs(user models.User) (e error, errMessages map[string][]string) {
	// Validate inputs
	hasError := false
	errMessages = make(map[string][]string)

	isValidEmail := helpers.Validator{}.IsValidEmailFormat(user.Email)
	if isValidEmail == false {
		hasError = true
		errMessages["Email"] = []string{"Invalid Email address"}
	}

	isValidUsername := helpers.Validator{}.IsValidUsernameFormat(user.Username)
	if isValidUsername == false {
		hasError = true
		errMessages["Username"] = []string{"Invalid username"}
	}

	if len(user.PlainPassword) < 6 {
		hasError = true
		errMessages["Password"] = []string{"Password length must be greater than 6 characters"}
	}

	if hasError {
		return errors.New("Validation failed"), errMessages
	}

	return nil, nil
}

//VerifySignup verify email address of registered user
// @Route: /api/signup/verify?token={token}
// Method: GET
func VerifySignup(w http.ResponseWriter, r *http.Request) {
	token, e := extractTokenFromURL(w, r)
	if e != nil {
		return
	}

	db := database.NewGORMInstance()
	defer db.Close()

	user, e := getUserByVerificationToken(w, db, token)
	if e != nil {
		return
	}

	e = activateUserAccount(w, db, user)
	if e != nil {
		return
	}

	helpers.NewSuccessResponse(w, "Congratulation, Your account has been activated!")
	return
}

//extractTokenFromURL get token parameter from query string e.g /api/signup/verify?token={token}
func extractTokenFromURL(w http.ResponseWriter, r *http.Request) (token string, e error) {
	token = r.URL.Query().Get("token")
	if strings.EqualFold(token, "") {
		err := helpers.ErrorResponse{}
		err.ErrorCode = http.StatusUnprocessableEntity
		err.Message = "Invalid verification token!"
		helpers.NewErrorResponse(w, &err)
		e = errors.New(err.Message)

		return
	}

	return
}

// getUserByVerificationToken load user object from DB using email_verification_token
// It handles http error responses
func getUserByVerificationToken(w http.ResponseWriter, db *gorm.DB, token string) (user models.User, e error) {
	db.Model(&user).Where("email_verification_token=?", token).First(&user)
	if user.ID == 0 {
		err := helpers.ErrorResponse{}
		err.ErrorCode = http.StatusNotFound
		err.Message = "Token not found!"
		helpers.NewErrorResponse(w, &err)
		e = errors.New(err.Message)

		return
	}

	return
}

//activateUserAccount will set EmailVerified to true and handle http error resoonses
func activateUserAccount(w http.ResponseWriter, db *gorm.DB, user models.User) (e error) {
	e = db.Model(&user).Where("id=?", user.ID).UpdateColumn(models.User{EmailVerified: true}).Error
	if e != nil {
		err := helpers.ErrorResponse{}
		err.ErrorCode = http.StatusNotFound
		err.Message = "Fail to activate account. Please try again!"
		err.Errors = make(map[string][]string)
		err.Errors["dev"] = []string{e.Error()}
		helpers.NewErrorResponse(w, &err)
		e = errors.New(err.Message)

		return
	}
	return
}
