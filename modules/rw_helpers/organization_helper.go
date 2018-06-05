// Package rw_helpers contains request helpers
// Most of this functions handle response errors
package rw_helpers

import (
	"devin/helpers"
	"devin/models"
	"devin/policies"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"

	"devin/modules/organization/repository"
)

// InvitationReqModel strcuct of invitation request
type InvitationReqModel struct {
	Identifier string `doc:"username or email of user"`
}

// OrganizationPermissionUpdatableData request model of organization permissions
type OrganizationPermissionUpdatableData struct {
	UserID                   uint64
	OrganizationID           uint64
	IsAdminOfOrganization    bool
	CanCreateProject         bool
	CanAddUserToOrganization bool
}

// ExtractOrganizationID Load organization ID from URL
func ExtractOrganizationID(w http.ResponseWriter, r *http.Request, paramName string) (uint64, error) {
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

// FetchOrganizationFromDB Check DB for existance of organization
func FetchOrganizationFromDB(w http.ResponseWriter, db *gorm.DB, organizationID uint64) (organization models.User, e error) {
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

//DecodeInviteRequestBody decode request body as InvitationReqModel
func DecodeInviteRequestBody(w http.ResponseWriter, r *http.Request) (reqModel InvitationReqModel, e error) {
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

// CanInviteUser check permission of authenticated user for invitation access
func CanInviteUser(w http.ResponseWriter, authUser, organization models.User) bool {

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

// IsTargetUserRegistered check invited email already registered in the application
// It returns error if the target email address not registered already.
func IsTargetUserRegistered(w http.ResponseWriter, db *gorm.DB, reqModel InvitationReqModel) (targetUser models.User, registered bool) {
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

//IsUserMemberOfOrganization check the invited user membership status,
//If user is a memeber of given organization,
//This method will return true
func IsUserMemberOfOrganization(w http.ResponseWriter, db *gorm.DB, organiztionID uint64, targetUserID uint64) bool {

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

// AlreadyInvited check for previous invitations of user for the given organization
func AlreadyInvited(w http.ResponseWriter, db *gorm.DB, organiztionID uint64, targetUserID uint64) bool {
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

// SaveInvitation save final data of invitation into the DB.
func SaveInvitation(w http.ResponseWriter, db *gorm.DB, targetUser models.User, OrganizationID, CreatedByID uint64) error {
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

// CanViewPendingInvitations check permission and handle its error
func CanViewPendingInvitations(w http.ResponseWriter, authUser models.User, targetUserID uint64) bool {
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

// ExtractAcceptanceStatusFromURL try to extract acceptance_status parameter from request URL path.
func ExtractAcceptanceStatusFromURL(w http.ResponseWriter, r *http.Request) (status bool, e error) {
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

// CanUserChangeAcceptanceStatus check edit permission and handle http errors
func CanUserChangeAcceptanceStatus(w http.ResponseWriter, user models.User, invitation models.UserOrganizationInvitation) bool {
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

// GetPendingInvitaion load pending invitaion from DB and handle http errors of failuer
func GetPendingInvitaion(w http.ResponseWriter, db *gorm.DB, ID uint64) (invitation models.UserOrganizationInvitation, e error) {
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

// ExtractUserIDFromURL extract a parameter passed in URL, convert to uint64 and returns as a user ID
func ExtractUserIDFromURL(w http.ResponseWriter, r *http.Request, paramName string) (ID uint64, e error) {
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

// DecodeOrganizationRequestModel decode request body to object of models.User
func DecodeOrganizationRequestModel(w http.ResponseWriter, r *http.Request) (reqModel models.User, e error) {

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

// CanCreateOrganization check permission of user to create new organization.
func CanCreateOrganization(w http.ResponseWriter, authUser, reqModel models.User) bool {
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

// CanViewOrganizationsOfUser check permission of authenticated user to view an organization
func CanViewOrganizationsOfUser(w http.ResponseWriter, authUser models.User, userID uint64) bool {
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

// IsOrganizationUsernameValid check validations of username
func IsOrganizationUsernameValid(w http.ResponseWriter, reqModel models.User) bool {
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

// SaveOrganization save final data to the DB
func SaveOrganization(w http.ResponseWriter, db *gorm.DB, reqModel models.User) (org models.User, e error) {
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

// CanUpdateUserOrganizationPermissions check permission of authenticated user
//to update permissions of given user on an organization
func CanUpdateUserOrganizationPermissions(w http.ResponseWriter, db *gorm.DB, authUser, organization models.User) bool {
	can := policies.CanUpdateUserOrganizationPermissions(db, authUser, organization)
	if can == true {
		return true
	}

	err := helpers.ErrorResponse{
		ErrorCode: http.StatusForbidden,
		Message:   "Operation not permitted.",
	}
	helpers.NewErrorResponse(w, &err)
	return false

}

// UpdateOrganizationPermissions handle updating of permissions and http errors
func UpdateOrganizationPermissions(w http.ResponseWriter, db *gorm.DB, reqModel OrganizationPermissionUpdatableData) (e error) {
	e = db.Model(&models.UserOrganization{}).
		Where("user_id=? AND organization_id=?", reqModel.UserID, reqModel.OrganizationID).
		UpdateColumns(&reqModel).Error
	if e != nil {
		err := helpers.ErrorResponse{
			Message:   "Fail to update",
			ErrorCode: http.StatusInternalServerError,
		}
		helpers.NewErrorResponse(w, &err)
		return
	}
	return
}
