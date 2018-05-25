package controllers

import (
	"devin/database"
	"devin/helpers"
	"devin/models"
	"devin/policies"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"

	"devin/modules/organization/repository"
)

//extractOrganizationID Load organization ID from URL
func extractOrganizationID(w http.ResponseWriter, r *http.Request, paramName string) (uint64, error) {
	orgID, ok := mux.Vars(r)[paramName]
	if !ok {
		err := helpers.ErrorResponse{
			Message:   "Invalid Organization ID.",
			ErrorCode: http.StatusUnprocessableEntity,
		}
		helpers.NewErrorResponse(w, &err)
		return 0, errors.New(err.Message)
	}

	OrganizationID, e := strconv.ParseUint(orgID, 10, 64)
	if e != nil {
		err := helpers.ErrorResponse{
			Message:   "Invalid Organization ID. Just integer values accepted",
			ErrorCode: http.StatusUnprocessableEntity,
		}
		helpers.NewErrorResponse(w, &err)
		return 0, errors.New(err.Message)
	}

	return OrganizationID, nil
}

// fetchOrganizationFromDB Check DB for existance of organization
func fetchOrganizationFromDB(w http.ResponseWriter, db *gorm.DB, organizationID uint64) (organization models.User, e error) {
	db.Model(&organization).Where("id=? AND user_type=2", organizationID).First(&organization)
	if organization.ID == 0 {
		err := helpers.ErrorResponse{
			Message:   "Organization not found",
			ErrorCode: http.StatusNotFound,
		}
		helpers.NewErrorResponse(w, &err)
		e = errors.New(err.Message)
		return
	}
	return
}

//decodeInviteRequestBody decode request body as invitationReqModel
func decodeInviteRequestBody(w http.ResponseWriter, r *http.Request) (reqModel invitationReqModel, e error) {
	// Decode request json data to model
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

// isEmailAddressEmpty validate that invitation email is not empty
func isEmailAddressEmpty(w http.ResponseWriter, reqModel invitationReqModel) bool {
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

// canInviteUser check permission of authenticated user for invitation access
func canInviteUser(w http.ResponseWriter, authUser, organization models.User) bool {

	//Check permission of user to invite others
	if policies.CanInviteUserToOrganization(authUser, organization) == false {
		err := helpers.ErrorResponse{
			ErrorCode: http.StatusForbidden,
			Message:   "This request is not permitted for you.",
		}
		helpers.NewErrorResponse(w, &err)

		return false
	}

	return true
}

// isTargetUserRegistered check invited email already registered in the application
// It returns error if the target email address not registered already.
func isTargetUserRegistered(w http.ResponseWriter, db *gorm.DB, reqModel invitationReqModel) (targetUser models.User, registered bool) {
	// Check for invited user registration: already registered or not
	db.Model(&targetUser).Where("email=? OR username=?", reqModel.Identifier, reqModel.Identifier).First(&targetUser)
	if targetUser.ID == 0 {
		err := helpers.ErrorResponse{
			ErrorCode: http.StatusNotFound,
			Message:   "Not found",
		}
		err.Errors = make(map[string][]string)
		err.Errors["Identifier"] = []string{"No user found with this username or email "}
		helpers.NewErrorResponse(w, &err)
		registered = false
		return
	}
	registered = true
	return
}

//isUserMemberOfOrganization check the invited user membership status,
//If user is a memeber of given organization,
//This method will return true
func isUserMemberOfOrganization(w http.ResponseWriter, db *gorm.DB, organiztionID uint64, targetUserID uint64) bool {

	var count uint
	db.Model(&models.UserOrganization{}).Where("user_id=? AND organization_id=?", targetUserID, organiztionID).Count(&count)

	if count > 0 {
		err := helpers.ErrorResponse{
			ErrorCode: http.StatusUnprocessableEntity,
			Message:   "User exists",
		}
		err.Errors = make(map[string][]string)
		err.Errors["Identifier"] = []string{"A user with the given email/username already added to this organization"}
		helpers.NewErrorResponse(w, &err)
		return true
	}

	return false
}

// alreadyInvited check for previous invitations of user for the given organization
func alreadyInvited(w http.ResponseWriter, db *gorm.DB, organiztionID uint64, targetUserID uint64) bool {
	var cnt uint
	db.Model(&models.UserOrganizationInvitation{}).Where("user_id=? AND organization_id=?", targetUserID, organiztionID).Count(&cnt)

	if cnt > 0 {
		err := helpers.ErrorResponse{
			ErrorCode: http.StatusUnprocessableEntity,
			Message:   "User already invited",
		}
		err.Errors = make(map[string][]string)
		err.Errors["Identifier"] = []string{"A user with the given email/username already invited to this organization"}
		helpers.NewErrorResponse(w, &err)
		return true
	}

	return false
}

// saveInvitation save final data of invitation into the DB.
func saveInvitation(w http.ResponseWriter, db *gorm.DB, targetUser models.User, OrganizationID, CreatedByID uint64) error {
	var invitation models.UserOrganizationInvitation
	invitation.Email = &targetUser.Email
	invitation.UserID = &targetUser.ID
	invitation.OrganizationID = OrganizationID
	invitation.CreatedByID = CreatedByID

	e := db.Save(&invitation).Error
	if e != nil {
		err := helpers.ErrorResponse{
			ErrorCode: http.StatusInternalServerError,
			Message:   "Fail to save data",
		}
		helpers.NewErrorResponse(w, &err)
		return errors.New(err.Message)
	}

	return nil

}

//PendingInvitationRequests return a list of pending organization invitations
//User ID passed in URL
func PendingInvitationRequests(w http.ResponseWriter, r *http.Request) {
	authUser, e := getAuthenticatedUser(w, r)
	if e != nil {
		return
	}

	userID, e := extractIDFromURL(w, r)
	if e != nil {
		return
	}

	db := database.NewGORMInstance()
	defer db.Close()

	if userExists(w, db, userID) != nil {
		return
	}

	if canViewPendingInvitations(w, authUser, userID) == false {
		return
	}

	pendingRequests, e := repository.GetPendingInvitaionsByUserID(db, userID)
	if e != nil {
		err := helpers.ErrorResponse{
			ErrorCode: http.StatusInternalServerError,
			Message:   "Fail to load pending requests",
		}
		helpers.NewErrorResponse(w, &err)
		return
	}

	json.NewEncoder(w).Encode(&pendingRequests)
	return
}

// canViewPendingInvitations check permission and handle its error
func canViewPendingInvitations(w http.ResponseWriter, authUser models.User, targetUserID uint64) bool {
	if policies.CanViewPendingInvitations(authUser, targetUserID) == false {
		err := helpers.ErrorResponse{
			ErrorCode: http.StatusForbidden,
			Message:   "This request is not permitted for you.",
		}
		helpers.NewErrorResponse(w, &err)

		return false
	}

	return true
}

// userExists check for existance of user with given ID
func userExists(w http.ResponseWriter, db *gorm.DB, userID uint64) (e error) {
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

//extractAcceptanceStatusFromURL try to extract acceptance_status parameter from request URL path.
func extractAcceptanceStatusFromURL(w http.ResponseWriter, r *http.Request) (status bool, e error) {
	statusString, ok := mux.Vars(r)["acceptance_status"]
	if ok == false {
		err := helpers.ErrorResponse{
			Message:   "Invalid acceptance status.",
			ErrorCode: http.StatusUnprocessableEntity,
		}
		helpers.NewErrorResponse(w, &err)
		e = errors.New(err.Message)
		return false, e
	}

	switch statusString {
	case "accept":
		return true, nil
	case "reject":
		return false, nil
	default:
		err := helpers.ErrorResponse{
			Message:   "Invalid acceptance status.",
			ErrorCode: http.StatusUnprocessableEntity,
		}
		helpers.NewErrorResponse(w, &err)
		e = errors.New(err.Message)
		return false, e
	}
}

// canUserChangeAcceptanceStatus check edit permission and handle http errors
func canUserChangeAcceptanceStatus(w http.ResponseWriter, user models.User, invitation models.UserOrganizationInvitation) bool {
	if policies.CanUserChangeAcceptanceStatus(user, invitation) {
		return true
	}

	err := helpers.ErrorResponse{
		Message:   "Operation not permitted",
		ErrorCode: http.StatusForbidden,
	}
	helpers.NewErrorResponse(w, &err)

	return false
}

//getPendingInvitaion load pending invitaion from DB and handle http errors of failuer
func getPendingInvitaion(w http.ResponseWriter, db *gorm.DB, ID uint64) (invitation models.UserOrganizationInvitation, e error) {
	invitation, e = repository.GetPendingInvitaionsByID(db, ID)
	if e != nil {
		err := helpers.ErrorResponse{
			Message:   "No match found",
			ErrorCode: http.StatusNotFound,
		}
		helpers.NewErrorResponse(w, &err)
		return
	}

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

func extractUserIDFromURL(w http.ResponseWriter, r *http.Request, paramName string) (ID uint64, e error) {
	IDString, ok := mux.Vars(r)[paramName]
	if ok == false {
		err := helpers.ErrorResponse{
			Message:   "Invalid User ID.",
			ErrorCode: http.StatusUnprocessableEntity,
		}
		helpers.NewErrorResponse(w, &err)
		e = errors.New(err.Message)
		return
	}

	ID, e = strconv.ParseUint(IDString, 10, 64)
	if e != nil {
		err := helpers.ErrorResponse{
			Message:   "Invalid User ID. Just integer values accepted",
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
