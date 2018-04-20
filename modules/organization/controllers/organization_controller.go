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

func Save(w http.ResponseWriter, r *http.Request) {

	// Check content type
	if !helpers.HasJSONRequest(r) {
		err := helpers.ErrorResponse{
			Message:   "Invalid content type.",
			ErrorCode: http.StatusUnsupportedMediaType,
		}
		helpers.NewErrorResponse(w, &err)
		return
	}

	if r.Body == nil {
		err := helpers.ErrorResponse{
			ErrorCode: http.StatusInternalServerError,
			Message:   "Request body cant be empty",
		}
		helpers.NewErrorResponse(w, &err)
		return
	}

	userID, ok := mux.Vars(r)["id"]
	if !ok {
		err := helpers.ErrorResponse{
			Message:   "Invalid User ID.",
			ErrorCode: http.StatusUnprocessableEntity,
		}
		helpers.NewErrorResponse(w, &err)
		return
	}

	ownerID, e := strconv.ParseUint(userID, 10, 64)
	if e != nil {
		err := helpers.ErrorResponse{
			Message:   "Invalid User ID. Just integer values accepted",
			ErrorCode: http.StatusUnprocessableEntity,
		}
		helpers.NewErrorResponse(w, &err)
		return
	}

	authUser, _, e := models.User{}.ExtractUserFromRequestContext(r)
	if e != nil {
		err := helpers.ErrorResponse{
			ErrorCode: http.StatusUnauthorized,
			Message:   "Auhtentication failed.",
		}
		helpers.NewErrorResponse(w, &err)
		return
	}
	var reqModel models.User

	e = json.NewDecoder(r.Body).Decode(&reqModel)
	if e != nil {
		err := helpers.ErrorResponse{
			ErrorCode: http.StatusInternalServerError,
			Message:   "Invalid request body",
		}
		helpers.NewErrorResponse(w, &err)

		return
	}

	reqModel.OwnerID = &ownerID
	reqModel.FirstName = reqModel.FullName
	reqModel.Username = strings.ToLower(reqModel.Username)
	reqModel.Email = strings.ToLower(reqModel.Email)
	reqModel.UserType = 2 //2 for organizations

	if !policies.CanCreateOrganization(authUser, reqModel) {
		err := helpers.ErrorResponse{
			ErrorCode: http.StatusForbidden,
			Message:   "This action is not allowed for you.",
		}
		helpers.NewErrorResponse(w, &err)
		return
	}

	// username validator
	isValidUsername := helpers.Validator{}.IsValidUsernameFormat(reqModel.Username)
	if isValidUsername == false {

		err := helpers.ErrorResponse{
			ErrorCode: http.StatusUnprocessableEntity,
			Message:   "Fail to save",
		}
		err.Errors = make(map[string][]string)
		err.Errors["Username"] = []string{"Invalid username"}
		helpers.NewErrorResponse(w, &err)

		return
	}

	// email validator
	isValidEmail := helpers.Validator{}.IsValidEmailFormat(reqModel.Email)
	if isValidEmail == false {

		err := helpers.ErrorResponse{
			ErrorCode: http.StatusUnprocessableEntity,
			Message:   "Fail to save",
		}
		err.Errors = make(map[string][]string)
		err.Errors["Email"] = []string{"Invalid email address"}
		helpers.NewErrorResponse(w, &err)

		return
	}

	db := database.NewGORMInstance()
	defer db.Close()

	// Check for duplication of email
	is, _ := reqModel.IsUniqueValue(db, "email", reqModel.Email, reqModel.ID)
	if is == false {

		err := helpers.ErrorResponse{
			Message:   "Invalid Email address.",
			ErrorCode: http.StatusUnprocessableEntity,
		}
		err.Errors = make(map[string][]string)
		err.Errors["Email"] = []string{"This email is already registered."}

		helpers.NewErrorResponse(w, &err)
		return

	}
	// Check for duplication of username
	is, _ = reqModel.IsUniqueValue(db, "username", reqModel.Username, reqModel.ID)
	if is == false {
		err := helpers.ErrorResponse{
			Message:   "Invalid username.",
			ErrorCode: http.StatusUnprocessableEntity,
		}
		err.Errors = make(map[string][]string)
		err.Errors["Username"] = []string{"This username is already registered."}

		helpers.NewErrorResponse(w, &err)
		return

	}

	e = db.Model(&reqModel).Save(&reqModel).Error
	if e != nil {
		err := helpers.ErrorResponse{
			ErrorCode: http.StatusInternalServerError,
			Message:   "Fail to in save in DB.",
		}
		helpers.NewErrorResponse(w, &err)
		return
	}

	json.NewEncoder(w).Encode(&reqModel)
	return
}
