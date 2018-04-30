package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"

	"devin/database"
	"devin/helpers"
	"devin/models"
	"devin/policies"
)

// InviteUser send invitaion request to given user
func InviteUser(w http.ResponseWriter, r *http.Request) {
	// Check content type
	if !helpers.HasJSONRequest(r) {
		err := helpers.ErrorResponse{
			Message:   "Invalid content type.",
			ErrorCode: http.StatusUnsupportedMediaType,
		}
		helpers.NewErrorResponse(w, &err)
		return
	}

	// Load organization ID from URL
	orgID, ok := mux.Vars(r)["id"]
	if !ok {
		err := helpers.ErrorResponse{
			Message:   "Invalid Organization ID.",
			ErrorCode: http.StatusUnprocessableEntity,
		}
		helpers.NewErrorResponse(w, &err)
		return
	}

	ID, e := strconv.ParseUint(orgID, 10, 64)
	if e != nil {
		err := helpers.ErrorResponse{
			Message:   "Invalid Organization ID. Just integer values accepted",
			ErrorCode: http.StatusUnprocessableEntity,
		}
		helpers.NewErrorResponse(w, &err)
		return
	}

	db := database.NewGORMInstance()
	defer db.Close()

	var organization models.User
	//Check DB for existance of organization
	db.Model(&organization).Where("id=? AND user_type=2", ID).First(&organization)
	if organization.ID == 0 {
		err := helpers.ErrorResponse{
			Message:   "Organization not found",
			ErrorCode: http.StatusNotFound,
		}
		helpers.NewErrorResponse(w, &err)
		return
	}

	// Get authenticated user
	authUser, _, e := models.User{}.ExtractUserFromRequestContext(r)
	if e != nil {
		err := helpers.ErrorResponse{
			ErrorCode: http.StatusUnauthorized,
			Message:   "Auhtentication failed.",
		}
		helpers.NewErrorResponse(w, &err)
		return
	}

	var reqModel models.UserOrganizationInvitation

	// Check request boby
	if r.Body == nil {
		err := helpers.ErrorResponse{
			ErrorCode: http.StatusInternalServerError,
			Message:   "Request body cant be empty",
		}
		helpers.NewErrorResponse(w, &err)
		return
	}

	// Decode request json data to model
	e = json.NewDecoder(r.Body).Decode(&reqModel)
	if e != nil {
		err := helpers.ErrorResponse{
			ErrorCode: http.StatusInternalServerError,
			Message:   "Invalid request body",
		}
		helpers.NewErrorResponse(w, &err)

		return
	}

	//Check permission of user to invite others
	if policies.CanInviteUserToOrganization(authUser, organization) == false {
		err := helpers.ErrorResponse{
			ErrorCode: http.StatusForbidden,
			Message:   "This request is not permitted for you.",
		}
		helpers.NewErrorResponse(w, &err)

		return
	}

	// Check email address of null data
	if reqModel.Email == nil || strings.EqualFold(*reqModel.Email, "") {
		err := helpers.ErrorResponse{
			ErrorCode: http.StatusUnprocessableEntity,
			Message:   "Invalid request data",
		}
		err.Errors = make(map[string][]string)
		err.Errors["Email"] = []string{"Invitation's email address can't be empty"}
		helpers.NewErrorResponse(w, &err)
		return
	}
	// Check for valid email address
	if !new(helpers.Validator).IsValidEmailFormat(*reqModel.Email) {
		err := helpers.ErrorResponse{
			ErrorCode: http.StatusUnprocessableEntity,
			Message:   "Invalid request data",
		}
		err.Errors = make(map[string][]string)
		err.Errors["Email"] = []string{"Invalid email address"}
		helpers.NewErrorResponse(w, &err)
		return
	}
	// Check for invited user registration: already registered or not
	// Send request email
	reqModel.OrganizationID = ID

	// Save request data in DB
}
