package controllers

import (
	"encoding/json"
	"net/http"

	"devin/database"
	"devin/helpers"
	"devin/modules/organization/repository"
	"devin/modules/rw_helpers"
)

// InviteUser send invitaion request to given user
// This method invite already registered users using
// their username or email
func InviteUser(w http.ResponseWriter, r *http.Request) {
	if rw_helpers.IsJSONRequest(w, r) == false {
		return
	}

	OrganizationID, e := rw_helpers.ExtractOrganizationID(w, r, "id")
	if e != nil {
		return
	}

	db := database.NewGORMInstance()
	defer db.Close()

	organization, e := rw_helpers.FetchOrganizationFromDB(w, db, OrganizationID)
	if e != nil {
		return
	}

	authUser, e := rw_helpers.GetAuthenticatedUser(w, r)
	if e != nil {
		return
	}

	if helpers.IsRequestBodyNil(w, r) == true {
		return
	}

	reqModel, e := rw_helpers.DecodeInviteRequestBody(w, r)
	if e != nil {
		return
	}

	if rw_helpers.IsEmailAddressEmpty(w, reqModel) == true {
		return
	}

	if rw_helpers.CanInviteUser(w, authUser, organization) == false {
		return
	}

	targetUser, ex := rw_helpers.IsTargetUserRegistered(w, db, reqModel)
	if ex == false {
		return
	}

	if rw_helpers.IsUserMemberOfOrganization(w, db, OrganizationID, targetUser.ID) == true {
		return
	}

	if rw_helpers.AlreadyInvited(w, db, OrganizationID, authUser.ID) == true {
		return
	}

	if rw_helpers.SaveInvitation(w, db, targetUser, OrganizationID, authUser.ID) != nil {
		return
	}

	helpers.NewSuccessResponse(w, "Invitation sent successfully")
}

//AcceptOrRejectInvitation Accept or reject an invitation request
//@Route=/api/invitation/{id}/set_acceptance/{acceptance_status}
func AcceptOrRejectInvitation(w http.ResponseWriter, r *http.Request) {
	authUser, e := rw_helpers.GetAuthenticatedUser(w, r)
	if e != nil {
		return
	}

	acceptanceStatus, e := rw_helpers.ExtractAcceptanceStatusFromURL(w, r)
	if e != nil {
		return
	}

	invitationID, e := rw_helpers.ExtractIDFromURL(w, r)
	if e != nil {
		return
	}

	db := database.NewGORMInstance()
	defer db.Close()

	invitation, e := rw_helpers.GetPendingInvitaion(w, db, invitationID)
	if e != nil {
		return
	}

	if rw_helpers.CanUserChangeAcceptanceStatus(w, authUser, invitation) == false {
		return
	}

	e = repository.SetAcceptanceStatusOfInvitaion(db, invitationID, acceptanceStatus)
	if e != nil {
		err := helpers.ErrorResponse{
			ErrorCode: http.StatusInternalServerError,
			Message:   "Fail to update acceptance status!",
		}
		err.Errors = make(map[string][]string)
		err.Errors["dev"] = []string{e.Error()}
		helpers.NewErrorResponse(w, &err)
		return
	}

	if acceptanceStatus == true {
		e = repository.AddUserToOrganziation(db, *invitation.UserID, authUser.ID, invitation.OrganizationID)
		if e != nil {
			err := helpers.ErrorResponse{
				ErrorCode: http.StatusUnprocessableEntity,
				Message:   "Fail to add user to organization!",
			}
			err.Errors = make(map[string][]string)
			err.Errors["dev"] = []string{e.Error()}
			helpers.NewErrorResponse(w, &err)
			return
		}
	}

	helpers.NewSuccessResponse(w, "Acceptance Status updated!")
}

//PendingInvitationRequests return a list of pending organization invitations
//User ID passed in URL
func PendingInvitationRequests(w http.ResponseWriter, r *http.Request) {
	authUser, e := rw_helpers.GetAuthenticatedUser(w, r)
	if e != nil {
		return
	}

	userID, e := rw_helpers.ExtractIDFromURL(w, r)
	if e != nil {
		return
	}

	db := database.NewGORMInstance()
	defer db.Close()

	if rw_helpers.UserExists(w, db, userID) != nil {
		return
	}

	if rw_helpers.CanViewPendingInvitations(w, authUser, userID) == false {
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
