package controllers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"

	"devin/database"
	"devin/helpers"
	"devin/models"
	"devin/modules/organization/repository"
	"devin/policies"
)

func init() {
	log.SetFlags(log.Lshortfile)
}

//UserOrganizationsIndex return a list of organizations that the given user is a member of that.
func UserOrganizationsIndex(w http.ResponseWriter, r *http.Request) {

	authUser, e := getAuthenticatedUser(w, r)
	if e != nil {
		return
	}

	var searchObject repository.OrganizationSearchable
	json.Unmarshal([]byte(r.URL.Query().Get("filters")), &searchObject)
	if authUser.IsRootUser == false {

		if searchObject.UserID == nil || *searchObject.UserID == 0 {
			err := helpers.ErrorResponse{
				Message:   "Invalid User ID.",
				ErrorCode: http.StatusUnprocessableEntity,
			}
			helpers.NewErrorResponse(w, &err)
			e = errors.New(err.Message)
			return
		}

		if canViewOrganizationsOfUser(w, authUser, *searchObject.UserID) == false {
			return
		}
	}

	db := database.NewGORMInstance()
	defer db.Close()

	orgs, e := repository.LoadOrganizationsFilter(db, searchObject)

	if e != nil {
		err := helpers.ErrorResponse{
			Message:   "Fail to load organizations data",
			ErrorCode: http.StatusInternalServerError,
		}
		helpers.NewErrorResponse(w, &err)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&orgs)
}

// Save process inserting and updating of organizations
func Save(w http.ResponseWriter, r *http.Request) {

	if isJsonRequest(w, r) == false {
		return
	}

	if isRequestBodyNil(w, r) == true {
		return
	}

	ownerID, e := extractIDFromURL(w, r)
	if e != nil {
		return
	}

	authUser, e := getAuthenticatedUser(w, r)
	if e != nil {
		return
	}

	reqModel, e := decodeOrganizationRequestModel(w, r)

	reqModel.OwnerID = &ownerID
	reqModel.FirstName = reqModel.FullName
	reqModel.Username = strings.ToLower(reqModel.Username)
	reqModel.Email = strings.ToLower(reqModel.Email)
	reqModel.UserType = 2 //2 for organizations

	if canCreateOrganization(w, authUser, reqModel) == false {
		return
	}

	if isOrganizationUsernameValid(w, reqModel) == false {
		return
	}

	if isOrganizationEmailValid(w, reqModel) == false {
		return
	}

	db := database.NewGORMInstance()
	defer db.Close()

	if isUniqueEmail(w, db, reqModel) == false {
		return
	}

	if isUniqueUsername(w, db, reqModel) == false {
		return
	}

	reqModel, e = saveOrganization(w, db, reqModel)
	if e != nil {
		return
	}

	json.NewEncoder(w).Encode(&reqModel)
	return
}

//isJsonRequest check request body for 'application/json' content type
func isJsonRequest(w http.ResponseWriter, r *http.Request) bool {
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

// isRequestBodyNil check request body to being not nil
func isRequestBodyNil(w http.ResponseWriter, r *http.Request) bool {
	// Check request boby
	if r.Body == nil {
		err := helpers.ErrorResponse{
			ErrorCode: http.StatusInternalServerError,
			Message:   "Request body cant be empty",
		}
		helpers.NewErrorResponse(w, &err)
		return true
	}

	return false
}

//extractIDFromURL get ID variable form request URL
func extractIDFromURL(w http.ResponseWriter, r *http.Request) (ID uint64, e error) {
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

// getAuthenticatedUser get user who is now logged into the application
func getAuthenticatedUser(w http.ResponseWriter, r *http.Request) (authUser models.User, e error) {
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

// decodeOrganizationRequestModel decode request body to object of models.User
func decodeOrganizationRequestModel(w http.ResponseWriter, r *http.Request) (reqModel models.User, e error) {

	e = json.NewDecoder(r.Body).Decode(&reqModel)
	if e != nil {
		err := helpers.ErrorResponse{
			ErrorCode: http.StatusInternalServerError,
			Message:   "Invalid request body",
		}
		helpers.NewErrorResponse(w, &err)

		return
	}

	return
}

// canCreateOrganization check permission of user to create new organization.
func canCreateOrganization(w http.ResponseWriter, authUser, reqModel models.User) bool {
	if !policies.CanCreateOrganization(authUser, reqModel) {
		err := helpers.ErrorResponse{
			ErrorCode: http.StatusForbidden,
			Message:   "This action is not allowed for you.",
		}
		helpers.NewErrorResponse(w, &err)
		return false
	}

	return true
}

func canViewOrganizationsOfUser(w http.ResponseWriter, authUser models.User, userID uint64) bool {
	if !policies.CanViewOrganizationsOfUser(authUser, userID) {
		err := helpers.ErrorResponse{
			ErrorCode: http.StatusForbidden,
			Message:   "This action is not allowed for you.",
		}
		helpers.NewErrorResponse(w, &err)
		return false
	}

	return true
}

// isOrganizationUsernameValid check validations of username
func isOrganizationUsernameValid(w http.ResponseWriter, reqModel models.User) bool {
	// username validator
	isValidUsername := helpers.Validator{}.IsValidUsernameFormat(reqModel.Username)
	if isValidUsername == false {

		err := helpers.ErrorResponse{
			ErrorCode: http.StatusUnprocessableEntity,
			Message:   "Fail to save",
		}
		err.Errors = make(map[string][]string)
		err.Errors["Username"] = []string{"Invalid username"}
		helpers.NewErrorResponse(w, &err)

		return false
	}

	return true
}

// isOrganizationEmailValid check validation of email
func isOrganizationEmailValid(w http.ResponseWriter, reqModel models.User) bool {
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

//isUniqueEmail check uniqueness of organization's email
func isUniqueEmail(w http.ResponseWriter, db *gorm.DB, reqModel models.User) bool {
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

//isUniqueUsername check uniqueness of organization's username
func isUniqueUsername(w http.ResponseWriter, db *gorm.DB, reqModel models.User) bool {

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

// saveOrganization save final data to the DB
func saveOrganization(w http.ResponseWriter, db *gorm.DB, reqModel models.User) (org models.User, e error) {
	e = db.Model(&reqModel).Save(&reqModel).Error
	if e != nil {
		err := helpers.ErrorResponse{
			ErrorCode: http.StatusInternalServerError,
			Message:   "Fail to in save in DB.",
		}
		helpers.NewErrorResponse(w, &err)
		return
	}

	return reqModel, nil
}
