package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"devin/database"
	"devin/helpers"
	"devin/models"
	"devin/policies"
)

// UpdateProfile handle user profile updates.
// Authorizations handlers loaded in middleware
// This function get json request of user model and update associated model
// If requested user_id == logged in user_id => user is trying to update his profile,
// Else, user is trying to update someone else => need authorization check
func UpdateProfile(w http.ResponseWriter, r *http.Request) {

	// Check content type
	if !helpers.HasJSONRequest(r) {
		err := helpers.ErrorResponse{
			Message:   "Invalid content type.",
			ErrorCode: http.StatusUnsupportedMediaType,
		}
		helpers.NewErrorResponse(w, &err)
		return
	}

	var user models.User
	userID, ok := mux.Vars(r)["id"]
	if !ok {
		err := helpers.ErrorResponse{
			Message:   "Invalid User ID.",
			ErrorCode: http.StatusUnprocessableEntity,
		}
		helpers.NewErrorResponse(w, &err)
		return
	}
	var e error
	user.ID, e = strconv.ParseUint(userID, 10, 64)
	if e != nil {
		err := helpers.ErrorResponse{
			Message:   "Invalid User ID. Just integer values accepted",
			ErrorCode: http.StatusUnprocessableEntity,
		}
		helpers.NewErrorResponse(w, &err)
		return
	}
	//Load current user data from DB
	db := database.NewPGInstance()
	defer db.Close()
	e = db.Model(&user).Where("id=?", user.ID).First()
	if e != nil {
		err := helpers.ErrorResponse{
			ErrorCode: http.StatusInternalServerError,
			Message:   "Error on loading user data.",
		}
		log.Println("Error on loading user data,", e)
		helpers.NewErrorResponse(w, &err)

		return
	}

	authUser, _, e := user.ExtractUserFromRequestContext(r)
	if e != nil {
		err := helpers.ErrorResponse{
			ErrorCode: http.StatusUnauthorized,
			Message:   "Auhtentication failed.",
		}
		log.Println("Auhtentication failed,", e)
		helpers.NewErrorResponse(w, &err)

		return
	}

	if !policies.CanEditUser(authUser, user) {
		err := helpers.ErrorResponse{
			ErrorCode: http.StatusForbidden,
			Message:   "This action is not allowed for you.",
		}
		log.Println("This action is not allowed for you.,")
		helpers.NewErrorResponse(w, &err)

		return
	}

	if r.Body == nil {
		err := helpers.ErrorResponse{
			ErrorCode: http.StatusInternalServerError,
			Message:   "Request body cant be empty",
		}
		log.Println("Request body cant be empty,", e)
		helpers.NewErrorResponse(w, &err)

		return
	}

	var profile models.PublicProfile
	e = json.NewDecoder(r.Body).Decode(&profile)
	if e != nil {
		err := helpers.ErrorResponse{
			ErrorCode: http.StatusInternalServerError,
			Message:   "Invalid request body",
		}
		log.Println("Error on profile data decoding,", e)
		helpers.NewErrorResponse(w, &err)

		return
	}

	user.PublicProfile = profile
	_, e = db.Model(&user).Where("id=?", user.ID).Update(&profile)
	if e != nil {
		err := helpers.ErrorResponse{
			ErrorCode: http.StatusInternalServerError,
			Message:   "Error on updating user data.",
		}
		log.Println("Error on updating user data,", e)
		helpers.NewErrorResponse(w, &err)

		return
	}
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&user)
}
