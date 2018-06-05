package controllers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"

	"devin/database"
	"devin/helpers"
	"devin/modules/organization/repository"
	"devin/modules/rw_helpers"
)

func init() {
	log.SetFlags(log.Lshortfile)
}

//UserOrganizationsIndex return a list of organizations that the given user is a member of that.
func UserOrganizationsIndex(w http.ResponseWriter, r *http.Request) {

	authUser, e := rw_helpers.GetAuthenticatedUser(w, r)
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

		if rw_helpers.CanViewOrganizationsOfUser(w, authUser, *searchObject.UserID) == false {
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

	if rw_helpers.IsJSONRequest(w, r) == false {
		return
	}

	if helpers.IsRequestBodyNil(w, r) == true {
		return
	}

	ownerID, e := rw_helpers.ExtractIDFromURL(w, r)
	if e != nil {
		return
	}

	authUser, e := rw_helpers.GetAuthenticatedUser(w, r)
	if e != nil {
		return
	}

	reqModel, e := rw_helpers.DecodeOrganizationRequestModel(w, r)

	reqModel.OwnerID = &ownerID
	reqModel.FirstName = reqModel.FullName
	reqModel.Username = strings.ToLower(reqModel.Username)
	reqModel.Email = strings.ToLower(reqModel.Email)
	reqModel.UserType = 2 //2 for organizations

	if rw_helpers.CanCreateOrganization(w, authUser, reqModel) == false {
		return
	}

	if rw_helpers.IsOrganizationUsernameValid(w, reqModel) == false {
		return
	}

	if rw_helpers.IsValidEmail(w, reqModel) == false {
		return
	}

	db := database.NewGORMInstance()
	defer db.Close()

	if rw_helpers.IsUniqueEmail(w, db, reqModel) == false {
		return
	}

	if rw_helpers.IsUniqueUsername(w, db, reqModel) == false {
		return
	}

	reqModel, e = rw_helpers.SaveOrganization(w, db, reqModel)
	if e != nil {
		return
	}

	json.NewEncoder(w).Encode(&reqModel)
	return
}
