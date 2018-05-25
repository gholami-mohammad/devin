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
