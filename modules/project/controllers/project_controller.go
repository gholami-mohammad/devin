package controllers

import (
	"encoding/json"
	"net/http"

	"devin/database"
	"devin/modules/rw_helpers"
)

// ProjectController handle functionalities of Project
type ProjectController struct{}

// Save handle inserting and updating of project
// Request body is json encoded of project model
// If no ID present in the request model, it will insert as new project
// otherwise the given project will be updated
// @Route: /api/project/save
// @Content-Type: application/json
func (ProjectController) Save(w http.ResponseWriter, r *http.Request) {
	//get authenticated user
	authUser, e := rw_helpers.GetAuthenticatedUser(w, r)
	if e != nil {
		return
	}

	// Decode request body to project
	projectReqModel, e := rw_helpers.DecodeProjectRequestModel(w, r)
	if e != nil {
		return
	}

	db := database.NewGORMInstance()
	defer db.Close()

	e = rw_helpers.CheckOwnerOrganizationIDOfProject(w, db, projectReqModel.OwnerOrganizationID)
	if e != nil {
		return
	}

	//Check permission
	if rw_helpers.CanSaveProject(w, db, authUser, projectReqModel.OwnerOrganizationID, projectReqModel) == false {
		return

	}

	if projectReqModel.ID != 0 {
		db.Model(&projectReqModel).Where("id=?", projectReqModel.ID).Update(&projectReqModel)
	} else {
		db.Model(&projectReqModel).Create(&projectReqModel)
	}

	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&projectReqModel)
}
