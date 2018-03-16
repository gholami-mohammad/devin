package controllers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"

	"devin/auth"
	"devin/database"
	"devin/helpers"
	"devin/models"
)

type SigninReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
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

	claims := models.Claim{
		ID:       user.ID,
		Email:    user.Email,
		Username: user.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(user.TokenLifetime()).Unix(),
			Issuer:    "devin",
		},
	}

	// create a signer for rsa 512
	t := jwt.NewWithClaims(jwt.SigningMethodRS512, claims)

	sk, e := auth.GetJWTSignKey()
	if e != nil {
		err := helpers.ErrorResponse{
			Message:   "Internal server error(load jwt)",
			ErrorCode: http.StatusInternalServerError,
		}
		helpers.NewErrorResponse(w, &err)
		return
	}
	tokenString, err := t.SignedString(sk)
	if err != nil {
		err := helpers.ErrorResponse{
			Message:   "Internal server error(sign jwt)",
			ErrorCode: http.StatusInternalServerError,
		}
		helpers.NewErrorResponse(w, &err)
		return
	}
	user.SetAuthorizationCookie(w, tokenString)
	db.Update(&user)

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
