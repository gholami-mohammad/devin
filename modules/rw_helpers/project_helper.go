package rw_helpers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/jinzhu/gorm"

	"devin/helpers"
	"devin/models"
	"devin/policies"
)

// DecodeProjectSearchFilters get request body, decode from json to ProjectSearch
// Handle request body erros
func DecodeProjectSearchFilters(w http.ResponseWriter, r *http.Request) (searchModel models.ProjectSearch, e error) {
	q := r.URL.Query().Get("q")
	if strings.EqualFold(q, "") {
		q = `{}`
	}
	e = json.Unmarshal([]byte(q), &searchModel)
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

// ValidateProjectRequestModel will check request data for creating or updating of a project
func ValidateProjectRequestModel(w http.ResponseWriter, db *gorm.DB, reqModel models.Project) (err error) {
	resErr := helpers.ErrorResponse{}
	resErr.Errors = make(map[string][]string)

	// name field is required
	if strings.EqualFold(reqModel.Name, "") {
		resErr.Errors["name"] = append(resErr.Errors["name"], "The given name already taken!")
	}

	// check uniqueness of name
	if isProjectNameUnique(db, reqModel.Name, reqModel.ID) == false {
		resErr.Errors["name"] = append(resErr.Errors["name"], "Project name can't be empty!")
	}

	if len(resErr.Errors) == 0 {
		return nil
	}
	resErr.ErrorCode = http.StatusUnprocessableEntity
	resErr.Message = "Invalid data!"
	helpers.NewErrorResponse(w, &resErr)

	return errors.New(resErr.Message)
}

// isProjectNameUnique check uniqueness of project name
// If you pass 0 to the projectID, it will check all records
// otherwise check all records insted of given ID
func isProjectNameUnique(db *gorm.DB, projectName string, projectID uint64) bool {
	var project models.Project
	if projectID <= 0 {
		db.Where("name = ?", projectName).First(&project)
	} else {
		db.Where("name = ?", projectName).Where("id != ?", projectID).First(&project)
	}

	if project.ID == 0 {
		// Project name IS unique
		return true
	}

	// Project name IS NOT unique
	return false
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
