package rw_helpers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/jinzhu/gorm"

	"devin/helpers"
	"devin/models"
	"devin/policies"
)

// DecodeProjectSearchFilters get request body, decode from json to ProjectSearch
// Handle request body erros
func DecodeProjectSearchFilters(w http.ResponseWriter, r *http.Request) (searchModel models.ProjectSearch, e error) {
	if helpers.IsRequestBodyNil(w, r) == true {
		e = errors.New("request body can't be empty")
		return
	}
	e = json.NewDecoder(r.Body).Decode(&searchModel)
	if e != nil {
		err := helpers.ErrorResponse{Message: "Invalid search filters", ErrorCode: http.StatusUnprocessableEntity}
		err.Errors = make(map[string][]string)
		err.Errors["dev"] = []string{e.Error()}
		helpers.NewErrorResponse(w, &err)

		return
	}

	return
}

// CanSaveProject check permission of authenticated user for
// insert or update a project inside an organization
func CanSaveProject(w http.ResponseWriter, db *gorm.DB, authUser models.User, ownerOrganizationID *uint64, projectReqModel models.Project) bool {
	var e error
	var project models.Project

	if projectReqModel.ID != 0 {
		// Edit mode
		project, e = GetProjectByID(w, db, projectReqModel.ID)
		if e != nil {
			return false
		}

		//Check permission
		if policies.CanSaveProject(db, authUser, ownerOrganizationID, project) == false {
			err := helpers.ErrorResponse{
				ErrorCode: http.StatusForbidden,
				Message:   "You don't have permission to update this project!",
			}
			helpers.NewErrorResponse(w, &err)

			return false
		}
	} else {
		//Check permission
		if policies.CanSaveProject(db, authUser, ownerOrganizationID, projectReqModel) == false {
			err := helpers.ErrorResponse{
				ErrorCode: http.StatusForbidden,
				Message:   "You don't have permission to create project for this ogranization!",
			}
			helpers.NewErrorResponse(w, &err)

			return false
		}
	}

	return true
}

// DecodeProjectRequestModel check request body data and try to decode it to a project object
func DecodeProjectRequestModel(w http.ResponseWriter, r *http.Request) (project models.Project, e error) {
	if helpers.IsRequestBodyNil(w, r) {
		e = errors.New("Request body is nil!")
		return
	}
	e = json.NewDecoder(r.Body).Decode(&project)

	if e != nil {
		err := helpers.ErrorResponse{}
		err.ErrorCode = http.StatusBadRequest
		err.Message = "Invalid request!"
		helpers.NewErrorResponse(w, &err)

		return
	}

	return
}

// GetProjectByID try to load project from DB. If no item found, returns an error.
// This function handle http response errors
func GetProjectByID(w http.ResponseWriter, db *gorm.DB, id uint64) (project models.Project, e error) {
	db.Model(&project).Where("id=?", id).First(&project)

	if project.ID == 0 {
		err := helpers.ErrorResponse{}
		err.ErrorCode = http.StatusNotFound
		err.Message = "No matching project found!"
		helpers.NewErrorResponse(w, &err)

		e = errors.New(err.Message)
		return
	}
	return
}

// CheckOwnerOrganizationIDOfProject will check requested data of projet
// If ownerOrganizationID is nil, it return with no error
// Otherwise check selected organization existance in DB
func CheckOwnerOrganizationIDOfProject(w http.ResponseWriter, db *gorm.DB, ownerOrganizationID *uint64) (e error) {
	if ownerOrganizationID == nil || *ownerOrganizationID == 0 {
		return
	}
	// Load organization form DB
	_, e = FetchOrganizationFromDB(w, db, *ownerOrganizationID)
	if e != nil {
		return
	}
	return
}
