package controllers

import (
	"net/http"

	"devin/database"
	"devin/helpers"
	"devin/modules/organization/repository"
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

	OrganizationID, e := extractOrganizationID(w, r, "id")
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
