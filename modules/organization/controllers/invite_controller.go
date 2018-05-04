package controllers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"

	"devin/database"
	"devin/helpers"
	"devin/models"
	"devin/policies"
)

type invitationReqModel struct {
	Identifier string `doc:"username or email of user"`
}

// InviteUser send invitaion request to given user
// This method invite already registered users using
// their username or email
func InviteUser(w http.ResponseWriter, r *http.Request) {
	if isJsonRequest(w, r) == false {
		return
	}

	OrganizationID, e := extractOrganizationID(w, r)
	if e != nil {
		return
	}

	db := database.NewGORMInstance()
	defer db.Close()

	organization, e := fetchOrganizationFromDB(w, db, OrganizationID)
	if e != nil {
		return
	}

	authUser, e := getAuthenticatedUser(w, r)
	if e != nil {
		return
	}

	if isRequestBodyNil(w, r) == true {
		return
	}

	reqModel, e := decodeInviteRequestBody(w, r)
	if e != nil {
		return
	}

	if isEmailAddressEmpty(w, reqModel) == true {
		return
	}

	if canInviteUser(w, authUser, organization) == false {
		return
	}

	targetUser, ex := isTargetUserRegistered(w, db, reqModel)
	if ex == false {
		return
	}

	if isUserMemberOfOrganization(w, db, OrganizationID, targetUser.ID) == true {
		return
	}

	if saveInvitation(w, db, targetUser, OrganizationID, authUser.ID) != nil {
		return
	}

	helpers.NewSuccessResponse(w, "Invitation sent successfully")
}

//extractOrganizationID Load organization ID from URL
func extractOrganizationID(w http.ResponseWriter, r *http.Request) (uint64, error) {
	orgID, ok := mux.Vars(r)["id"]
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
