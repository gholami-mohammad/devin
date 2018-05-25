package controllers

import (
	"encoding/json"
	"net/http"

	"devin/database"
	"devin/helpers"
	"devin/models"
)

//UpdateUserPermissionsOnOrganization Update permissions of given user on the given organization
// @Method: POST
// @Route: /api/organization/{organization_id:[0-9]+}/user/{user_id:[0-9]+}/update_permissions
func UpdateUserPermissionsOnOrganization(w http.ResponseWriter, r *http.Request) {
	if isJsonRequest(w, r) == false {
		return
	}

	authUser, e := getAuthenticatedUser(w, r)
	if e != nil {
		return
	}

	organizationID, e := extractOrganizationID(w, r, "organization_id")
	if e != nil {
		return
	}

	userID, e := extractUserIDFromURL(w, r, "user_id")
	if e != nil {
		return
	}

	db := database.NewGORMInstance()
	defer db.Close()

	organization, e := fetchOrganizationFromDB(w, db, organizationID)
	if e != nil {
		return
	}

	if canUpdateUserOrganizationPermissions(w, db, authUser, organization) == false {
		return
	}

	if isRequestBodyNil(w, r) {
		return
	}

	type updatableData struct {
		UserID                   uint64
		OrganizationID           uint64
		IsAdminOfOrganization    bool
		CanCreateProject         bool
		CanAddUserToOrganization bool
	}

	var reqModel updatableData
	e = json.NewDecoder(r.Body).Decode(&reqModel)
	if e != nil {
		err := helpers.ErrorResponse{
			Message:   "Invalid request!",
			ErrorCode: http.StatusUnprocessableEntity,
		}
		helpers.NewErrorResponse(w, &err)
		return
	}
	defer r.Body.Close()

	reqModel.UserID = userID
	reqModel.OrganizationID = organizationID

	db.Model(&models.UserOrganization{}).
		Where("user_id=? AND organization_id=?", reqModel.UserID, reqModel.OrganizationID).
		UpdateColumns(&reqModel)

	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&reqModel)
}
