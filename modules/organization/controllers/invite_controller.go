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
	"devin/modules/organization/repository"
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

	if alreadyInvited(w, db, OrganizationID, authUser.ID) == true {
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

//AcceptOrRejectInvitation Accept or reject an invitation request
//@Route=/api/invitation/{id}/set_acceptance/{acceptance_status}
func AcceptOrRejectInvitation(w http.ResponseWriter, r *http.Request) {
	authUser, e := getAuthenticatedUser(w, r)
	if e != nil {
		return
	}

	acceptanceStatus, e := extractAcceptanceStatusFromURL(w, r)
	if e != nil {
		return
	}

	invitationID, e := extractIDFromURL(w, r)
	if e != nil {
		return
	}

	db := database.NewGORMInstance()
	defer db.Close()

	invitation, e := getPendingInvitaion(w, db, invitationID)
	if e != nil {
		return
	}

	if canUserChangeAcceptanceStatus(w, authUser, invitation) == false {
		return
	}

	e = repository.SetAcceptanceStatusOfInvitaion(db, invitationID, acceptanceStatus)
	if e != nil {
		err := helpers.ErrorResponse{
			ErrorCode: http.StatusInternalServerError,
			Message:   "Fail to update acceptance status!",
		}
		e = errors.New(err.Message)
		helpers.NewErrorResponse(w, &err)
		return
	}

	helpers.NewSuccessResponse(w, "Acceptance Status updated!")
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
