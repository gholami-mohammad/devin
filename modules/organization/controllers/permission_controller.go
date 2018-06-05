package controllers

import (
	"encoding/json"
	"net/http"

	"devin/database"
	"devin/helpers"
	"devin/modules/rw_helpers"
)

//UpdateUserPermissionsOnOrganization Update permissions of given user on the given organization
// @Method: POST
// @Route: /api/organization/{organization_id:[0-9]+}/user/{user_id:[0-9]+}/update_permissions
func UpdateUserPermissionsOnOrganization(w http.ResponseWriter, r *http.Request) {
	if rw_helpers.IsJSONRequest(w, r) == false {
		return
	}

	authUser, e := rw_helpers.GetAuthenticatedUser(w, r)
	if e != nil {
		return
	}

	organizationID, e := rw_helpers.ExtractOrganizationID(w, r, "organization_id")
	if e != nil {
		return
	}

	userID, e := rw_helpers.ExtractUserIDFromURL(w, r, "user_id")
	if e != nil {
		return
	}

	db := database.NewGORMInstance()
	defer db.Close()

	organization, e := rw_helpers.FetchOrganizationFromDB(w, db, organizationID)
	if e != nil {
		return
	}

	e = rw_helpers.IsUserExists(w, db, userID)
	if e != nil {
		return
	}

	if rw_helpers.CanUpdateUserOrganizationPermissions(w, db, authUser, organization) == false {
		return
	}

	if helpers.IsRequestBodyNil(w, r) {
		return
	}

	var reqModel rw_helpers.OrganizationPermissionUpdatableData
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

	e = rw_helpers.UpdateOrganizationPermissions(w, db, reqModel)
	if e != nil {
		return
	}

	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&reqModel)
}
