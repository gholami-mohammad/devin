package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"

	"devin/database"
	"devin/helpers"
	"devin/models"
	"devin/policies"
)

func handleProfileSharedErrors(r *http.Request, db *gorm.DB) (user models.User, err *helpers.ErrorResponse) {

	// Check content type
	if !helpers.HasJSONRequest(r) {
		err = &helpers.ErrorResponse{
			Message:   "Invalid content type.",
			ErrorCode: http.StatusUnsupportedMediaType,
		}
		return
	}

	userID, ok := mux.Vars(r)["id"]
	if !ok {
		err = &helpers.ErrorResponse{
			Message:   "Invalid User ID.",
			ErrorCode: http.StatusUnprocessableEntity,
		}

		return
	}
	var e error
	user.ID, e = strconv.ParseUint(userID, 10, 64)
	if e != nil {
		err = &helpers.ErrorResponse{
			Message:   "Invalid User ID. Just integer values accepted",
			ErrorCode: http.StatusUnprocessableEntity,
		}
		return
	}
	//Load current user data from DB
	e = db.Where("id=?", user.ID).First(&user).Error
	if e != nil {
		err = &helpers.ErrorResponse{
			ErrorCode: http.StatusInternalServerError,
			Message:   "Error on loading user data.",
		}

		return
	}

	authUser, _, e := user.ExtractUserFromRequestContext(r)
	if e != nil {
		err = &helpers.ErrorResponse{
			ErrorCode: http.StatusUnauthorized,
			Message:   "Auhtentication failed.",
		}

		return
	}

	if !policies.CanEditUser(authUser, user) {
		err = &helpers.ErrorResponse{
			ErrorCode: http.StatusForbidden,
			Message:   "This action is not allowed for you.",
		}

		return
	}

	if r.Body == nil {
		err = &helpers.ErrorResponse{
			ErrorCode: http.StatusInternalServerError,
			Message:   "Request body cant be empty",
		}

		return
	}

	return

}

// UpdateProfile handle user profile updates.
// Authorizations handlers loaded in middleware
// This function get json request of user model and update associated model
// If requested user_id == logged in user_id => user is trying to update his profile,
// Else, user is trying to update someone else => need authorization check
func UpdateProfile(w http.ResponseWriter, r *http.Request) {
	db := database.NewGORMInstance()
	defer db.Close()
	user, err := handleProfileSharedErrors(r, db)
	if err != nil {
		helpers.NewErrorResponse(w, err)
		return
	}

	var profile models.PublicProfile
	e := json.NewDecoder(r.Body).Decode(&profile)
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
	e = db.Model(&user).Where("id=?", user.ID).Update(&profile).Error
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

// UpdateUsername will update username of giver userID
func UpdateUsername(w http.ResponseWriter, r *http.Request) {
	db := database.NewGORMInstance()
	defer db.Close()
	user, err := handleProfileSharedErrors(r, db)
	if err != nil {
		helpers.NewErrorResponse(w, err)
		return
	}

	var reqModel struct {
		Username string
	}

	e := json.NewDecoder(r.Body).Decode(&reqModel)
	if e != nil {
		err := helpers.ErrorResponse{
			ErrorCode: http.StatusInternalServerError,
			Message:   "Invalid request body",
		}
		log.Println("Error on profile data decoding,", e)
		helpers.NewErrorResponse(w, &err)

		return
	}

	//Check for valid characters
	isValidUsername := helpers.Validator{}.IsValidUsernameFormat(reqModel.Username)
	if isValidUsername == false {
		messages := make(map[string][]string)
		messages["Username"] = []string{"This username has invalid characters."}
		err := helpers.ErrorResponse{
			Message:   "Invalid username.",
			ErrorCode: http.StatusUnprocessableEntity,
			Errors:    messages,
		}
		helpers.NewErrorResponse(w, &err)
		return

	}

	//Check dupication
	isUnique, _ := user.IsUniqueValue(db, "username", reqModel.Username, user.ID)
	if isUnique == false {
		messages := make(map[string][]string)
		messages["Username"] = []string{"This username is already taken."}
		err := helpers.ErrorResponse{
			Message:   "Invalid username.",
			ErrorCode: http.StatusUnprocessableEntity,
			Errors:    messages,
		}
		helpers.NewErrorResponse(w, &err)
		return

	}

	//Update
	user.Username = reqModel.Username
	e = db.Model(&user).Update(&user).Error
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

// Whoami load profile data of current logged in user
func Whoami(w http.ResponseWriter, r *http.Request) {
	authUser, _, e := models.User{}.ExtractUserFromRequestContext(r)
	if e != nil {
		err := helpers.ErrorResponse{
			ErrorCode: http.StatusUnauthorized,
			Message:   "Auhtentication failed.",
		}
		log.Println("Auhtentication failed,", e)
		helpers.NewErrorResponse(w, &err)

		return
	}

	db := database.NewGORMInstance()
	defer db.Close()
	var user models.User

	//Loading required data from DB
	e = db.Where("id=?", authUser.ID).First(&user).Error
	if e != nil {
		err := helpers.ErrorResponse{
			ErrorCode: http.StatusUnauthorized,
			Message:   "Auhtentication failed.",
		}
		log.Println("Auhtentication failed,", e)
		helpers.NewErrorResponse(w, &err)

		return
	}

	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&user)
}

// ProfileBasicInfo return array of basic informations
// required to render profile edit form
func ProfileBasicInfo(w http.ResponseWriter, r *http.Request) {
	info := make(map[string]interface{})
	db := database.NewGORMInstance()
	defer db.Close()

	var countries []models.Country
	db.Find(&countries)
	var dateFormats []models.DateFormat
	db.Find(&dateFormats)
	var timeFormats []models.TimeFormat
	db.Find(&timeFormats)
	var calendarSystems []models.CalendarSystem
	db.Find(&calendarSystems)

	info["LocalizationLanguages"] = countries
	info["DateFormats"] = dateFormats
	info["TimeFormats"] = timeFormats
	info["CalendarSystems"] = calendarSystems
	info["OfficePhoneCountryCodes"] = countries
	info["HomePhoneCountryCodes"] = countries
	info["CellPhoneCountryCodes"] = countries
	info["FaxCountryCodes"] = countries
	info["Countries"] = countries

	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&info)
}
