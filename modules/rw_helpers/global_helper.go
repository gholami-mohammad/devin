package rw_helpers

import (
	"devin/helpers"
	"devin/models"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
)

// GetAuthenticatedUser get user who is now logged into the application
func GetAuthenticatedUser(w http.ResponseWriter, r *http.Request) (authUser models.User, e error) {
	authUser, _, e = models.User{}.ExtractUserFromRequestContext(r)
	if e != nil {
		err := helpers.ErrorResponse{
			ErrorCode: http.StatusUnauthorized,
			Message:   "Auhtentication failed.",
		}
		helpers.NewErrorResponse(w, &err)

		return
	}

	return
}

// IsValidEmail check validation of email
func IsValidEmail(w http.ResponseWriter, reqModel models.User) bool {
	// email validator
	isValidEmail := helpers.Validator{}.IsValidEmailFormat(reqModel.Email)
	if isValidEmail == false {

		err := helpers.ErrorResponse{
			ErrorCode: http.StatusUnprocessableEntity,
			Message:   "Fail to save",
		}
		err.Errors = make(map[string][]string)
		err.Errors["Email"] = []string{"Invalid email address"}
		helpers.NewErrorResponse(w, &err)

		return false
	}
	return true
}

//IsUniqueEmail check uniqueness of user/organization's email
func IsUniqueEmail(w http.ResponseWriter, db *gorm.DB, reqModel models.User) bool {
	// Check for duplication of email
	is, _ := reqModel.IsUniqueValue(db, "email", reqModel.Email, reqModel.ID)
	if is == false {

		err := helpers.ErrorResponse{
			Message:   "Invalid Email address.",
			ErrorCode: http.StatusUnprocessableEntity,
		}
		err.Errors = make(map[string][]string)
		err.Errors["Email"] = []string{"This email is already registered."}

		helpers.NewErrorResponse(w, &err)
		return false

	}

	return true
}

// IsEmailAddressEmpty validate that invitation email is not empty
func IsEmailAddressEmpty(w http.ResponseWriter, reqModel InvitationReqModel) bool {
	// Check email address of null data
	if strings.EqualFold(reqModel.Identifier, "") {
		err := helpers.ErrorResponse{
			ErrorCode: http.StatusUnprocessableEntity,
			Message:   "Invalid request data",
		}
		err.Errors = make(map[string][]string)
		err.Errors["Identifier"] = []string{"Username or email is required"}
		helpers.NewErrorResponse(w, &err)
		return true
	}
	return false
}

// IsUniqueUsername check uniqueness of organization's username
func IsUniqueUsername(w http.ResponseWriter, db *gorm.DB, reqModel models.User) bool {

	// Check for duplication of username
	is, _ := reqModel.IsUniqueValue(db, "username", reqModel.Username, reqModel.ID)
	if is == false {
		err := helpers.ErrorResponse{
			Message:   "Invalid username.",
			ErrorCode: http.StatusUnprocessableEntity,
		}
		err.Errors = make(map[string][]string)
		err.Errors["Username"] = []string{"This username is already registered."}

		helpers.NewErrorResponse(w, &err)

		return false
	}

	return true
}

// ExtractIDFromURL get ID variable form request URL
func ExtractIDFromURL(w http.ResponseWriter, r *http.Request) (ID uint64, e error) {
	IDString, ok := mux.Vars(r)["id"]
	if ok == false {
		err := helpers.ErrorResponse{
			Message:   "Invalid ID.",
			ErrorCode: http.StatusUnprocessableEntity,
		}
		helpers.NewErrorResponse(w, &err)
		e = errors.New(err.Message)
		return
	}

	ID, e = strconv.ParseUint(IDString, 10, 64)
	if e != nil {
		err := helpers.ErrorResponse{
			Message:   "Invalid ID. Just integer values accepted",
			ErrorCode: http.StatusUnprocessableEntity,
		}
		helpers.NewErrorResponse(w, &err)
		return
	}

	return ID, nil
}

// IsJSONRequest check request body for 'application/json' content type
func IsJSONRequest(w http.ResponseWriter, r *http.Request) bool {
	if !helpers.HasJSONRequest(r) {
		err := helpers.ErrorResponse{
			Message:   "Invalid content type.",
			ErrorCode: http.StatusUnsupportedMediaType,
		}
		helpers.NewErrorResponse(w, &err)
		return false
	}

	return true
}

// GetCurrectpage get the 'page' parameter form request query string
func GetCurrectpage(r *http.Request) uint64 {
	pageStr := r.URL.Query().Get("page")

	page, e := strconv.ParseUint(pageStr, 10, 64)
	if e != nil {
		page = 1
	}

	if page <= 0 {
		page = 1
	}

	return page
}

// GetPerPage get the 'per_page' parameter form request query string
func GetPerPage(r *http.Request) uint64 {
	perPageStr := r.URL.Query().Get("per_page")

	perPage, e := strconv.ParseUint(perPageStr, 10, 64)
	if e != nil {
		perPage = 10
	}

	if perPage <= 0 {
		perPage = 10
	}

	return perPage
}
