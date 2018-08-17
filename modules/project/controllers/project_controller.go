package controllers

import (
	"encoding/json"
	"net/http"

	"devin/database"
	"devin/models"
	project_repo "devin/modules/project/repository"
	"devin/modules/rw_helpers"
)

// ProjectController handle functionalities of Project
type ProjectController struct{}

// ProjectsIndex return list of projects owned by given user or this user is a member of that
func (ProjectController) ProjectsIndex(w http.ResponseWriter, r *http.Request) {
	authUser, _, _ := models.User{}.ExtractUserFromRequestContext(r)

	// Decode search data
	searchModel, e := rw_helpers.DecodeProjectSearchFilters(w, r)
	if e != nil {
		return
	}

	searchModel.PerPage = rw_helpers.GetPerPage(r)
	searchModel.CurrentPage = rw_helpers.GetCurrectpage(r)

	defer r.Body.Close()

	db := database.NewGORMInstance()
	defer db.Close()

	data, total, _ := project_repo.SearchProjects(db, authUser, searchModel)
	var pgn models.Pagination
	pgn.Make(data, total, searchModel.CurrentPage, searchModel.PerPage)

	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&pgn)

	return
}

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

// BasicInfo load basic information to render project create/update form
// @Route: /api/projects/basic_info
func (ProjectController) BasicInfo(w http.ResponseWriter, r *http.Request) {
	db := database.NewGORMInstance()
	defer db.Close()

	var basicInfo struct {
		Statuses      []models.ProjectStatus
		Organizations []models.User
	}

	db.Find(&basicInfo.Statuses)
	user, err := rw_helpers.GetAuthenticatedUser(w, r)
	if err != nil {
		return
	}

	db.Model(&models.User{}).
		Where("user_type=2").
		Where("owner_id=?", user.ID).
		Find(&basicInfo.Organizations)
	db.Raw(`select * from users u where user_type=2 and owner_id=?
	union
select u.* from users u
	where u.id in (
			select organization_id 
			from user_organization uo 
			where uo.user_id = ? 
			and (uo.is_admin_of_organization=true or uo.can_create_project=true or uo.can_update_project=true ) 
	)`, user.ID, user.ID).
		Scan(&basicInfo.Organizations)
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&basicInfo)
}
