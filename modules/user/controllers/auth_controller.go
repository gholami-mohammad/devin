package controllers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"golang.org/x/crypto/bcrypt"

	"devin/database"
	"devin/helpers"
	"devin/models"
)

type SigninReq struct {
	Email    string
	Password string
}

// Signin handle user login.
// Method: POST
// Content-Type: application/json
func Signin(w http.ResponseWriter, r *http.Request) {
	// Check content type
	if !helpers.HasJSONRequest(r) {
		err := helpers.ErrorResponse{
			Message:   "Invalid content type.",
			ErrorCode: http.StatusUnsupportedMediaType,
		}
		helpers.NewErrorResponse(w, &err)
		return
	}

	var userReq SigninReq

	e := json.NewDecoder(r.Body).Decode(&userReq)
	if e != nil {
		err := helpers.ErrorResponse{
			Message:   "Invalid request body.",
			ErrorCode: http.StatusBadRequest,
		}
		helpers.NewErrorResponse(w, &err)
		return
	}

	e, errorMessages := validateSignipInputs(userReq)
	if e != nil {
		err := helpers.ErrorResponse{
			Message:   e.Error(),
			ErrorCode: http.StatusUnprocessableEntity,
			Errors:    errorMessages,
		}
		helpers.NewErrorResponse(w, &err)
		return
	}

	db := database.NewPGInstance()
	defer db.Close()
	var user models.User

	isEmail := helpers.Validator{}.IsValidEmailFormat(userReq.Email)
	if isEmail {
		db.Model(&user).Where("email=?", userReq.Email).First()
	} else {
		db.Model(&user).Where("username=?", userReq.Email).First()
	}

	if user.ID == 0 {
		err := helpers.ErrorResponse{
			Message:   "Invalid email/username or password",
			ErrorCode: http.StatusUnauthorized,
		}
		helpers.NewErrorResponse(w, &err)
		return
	}

	if user.EmailVerified == false {
		err := helpers.ErrorResponse{
			Message:   "Please verify your email address then try to login.",
			ErrorCode: http.StatusUnauthorized,
		}
		helpers.NewErrorResponse(w, &err)
		return
	}

	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userReq.Password)) != nil {
		err := helpers.ErrorResponse{
			Message:   "Invalid email/username or password",
			ErrorCode: http.StatusUnauthorized,
		}
		helpers.NewErrorResponse(w, &err)
		return
	}

	claim := user.GenerateNewTokenClaim()

	tokenString, err := user.GenerateNewTokenString(claim)
	if err != nil {
		helpers.NewErrorResponse(w, err)
		return
	}
	user.SetAuthorizationCookieAndHeader(w, tokenString)

	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&user)
}

// validateSignipInputs check signin requirments
func validateSignipInputs(user SigninReq) (e error, errMessages map[string][]string) {
	// Validate inputs
	hasError := false
	errMessages = make(map[string][]string)

	if strings.EqualFold(user.Email, "") {
		hasError = true
		errMessages["email"] = []string{"Emain or username is required"}
	}

	if strings.EqualFold(user.Password, "") {
		hasError = true
		errMessages["password"] = []string{"Password is required"}
	}

	if hasError {
		return errors.New("Signin failed"), errMessages
	}

	return nil, nil
}
