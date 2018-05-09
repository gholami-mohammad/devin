package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"

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

	reqModel.Username = strings.ToLower(reqModel.Username)

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
	e = db.Preload("Country").Preload("Province").Preload("City").Where("id=?", authUser.ID).First(&user).Error
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

// Whois load profile data of given userID
func Whois(w http.ResponseWriter, r *http.Request) {
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

	userID, ok := mux.Vars(r)["id"]
	if !ok {
		err := helpers.ErrorResponse{
			Message:   "Invalid User ID.",
			ErrorCode: http.StatusUnprocessableEntity,
		}
		helpers.NewErrorResponse(w, &err)
		return
	}

	var user models.User
	user.ID, e = strconv.ParseUint(userID, 10, 64)
	if e != nil {
		err := helpers.ErrorResponse{
			Message:   "Invalid User ID. Just integer values accepted",
			ErrorCode: http.StatusUnprocessableEntity,
		}
		helpers.NewErrorResponse(w, &err)
		return
	}

	db := database.NewGORMInstance()
	defer db.Close()

	//Loading required data from DB
	e = db.
		Preload("UserOrganizationMapping").
		Preload("OrganizationUserMapping").
		Preload("Country").
		Preload("Province").
		Preload("City").
		Where("id=?", userID).
		First(&user).
		Error
	if e != nil {
		err := helpers.ErrorResponse{
			ErrorCode: http.StatusUnauthorized,
			Message:   "Auhtentication failed.",
		}
		log.Println("Auhtentication failed,", e)
		helpers.NewErrorResponse(w, &err)

		return
	}

	user.SetFullName()

	if !policies.CanViewProfile(authUser, user) {
		err := helpers.ErrorResponse{
			ErrorCode: http.StatusForbidden,
			Message:   "This action is not allowed for you.",
		}
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
	db.Preload("Provinces").Preload("Provinces.Cities").Find(&countries)
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

// UpdatePassword handle updating of user's password
func UpdatePassword(w http.ResponseWriter, r *http.Request) {
	db := database.NewGORMInstance()
	defer db.Close()
	user, err := handleProfileSharedErrors(r, db)
	if err != nil {
		helpers.NewErrorResponse(w, err)
		return
	}
	var reqModel struct {
		Password             string
		PasswordVerification string
	}
	defer r.Body.Close()
	e := json.NewDecoder(r.Body).Decode(&reqModel)
	if e != nil {
		err := helpers.ErrorResponse{
			ErrorCode: http.StatusInternalServerError,
			Message:   "Invalid request body",
		}
		log.Println("Error on password data decoding,", e)
		helpers.NewErrorResponse(w, &err)

		return
	}

	//Check min length
	if len(reqModel.Password) < 6 {
		err := helpers.ErrorResponse{
			ErrorCode: http.StatusUnprocessableEntity,
			Message:   "Invalid password",
		}
		err.Errors = make(map[string][]string)
		err.Errors["Password"] = []string{"New password must has at least 6 characters."}
		helpers.NewErrorResponse(w, &err)

		return
	}
	//Check matching
	if !strings.EqualFold(reqModel.Password, reqModel.PasswordVerification) {
		err := helpers.ErrorResponse{
			ErrorCode: http.StatusUnprocessableEntity,
			Message:   "Invalid password",
		}
		err.Errors = make(map[string][]string)
		err.Errors["VerifyPassword"] = []string{"Password verification does not match."}
		helpers.NewErrorResponse(w, &err)

		return
	}
	//Change password
	user.SetEncryptedPassword(reqModel.Password)
	e = db.Model(&user).Where("id=?", user.ID).Update(&user).Error
	if e != nil {
		err := helpers.ErrorResponse{
			ErrorCode: http.StatusInternalServerError,
			Message:   "Error on updating password.",
		}
		log.Println("Error on updating password,", e)
		helpers.NewErrorResponse(w, &err)

		return
	}

	helpers.NewSuccessResponse(w, "Password updated.")
	return
}

// UpdateAvatar handle updating of user/organization avatar/logo
func UpdateAvatar(w http.ResponseWriter, r *http.Request) {
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

	userID, ok := mux.Vars(r)["id"]
	if !ok {
		err := helpers.ErrorResponse{
			Message:   "Invalid User ID.",
			ErrorCode: http.StatusUnprocessableEntity,
		}
		helpers.NewErrorResponse(w, &err)
		return
	}
	var user models.User
	user.ID, e = strconv.ParseUint(userID, 10, 64)
	if e != nil {
		err := helpers.ErrorResponse{
			Message:   "Invalid User ID. Just integer values accepted",
			ErrorCode: http.StatusUnprocessableEntity,
		}
		helpers.NewErrorResponse(w, &err)
		return
	}
	db := database.NewGORMInstance()
	defer db.Close()

	//Loading required data from DB
	e = db.Where("id=?", userID).First(&user).Error
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
		helpers.NewErrorResponse(w, &err)
		return
	}

	r.ParseMultipartForm(1024)
	file, fileHeader, e := r.FormFile("AvatarFile")
	if e != nil {
		err := helpers.ErrorResponse{
			ErrorCode: http.StatusUnprocessableEntity,
			Message:   "Error on updating avatar.",
		}
		err.Errors = make(map[string][]string)
		err.Errors["AvatarFile"] = []string{"Can't read file"}
		log.Println("Error on updating avatar,", e)

		helpers.NewErrorResponse(w, &err)

		return
	}
	defer file.Close()
	mimeType := fileHeader.Header.Get("Content-Type")

	if !isImageMimeType(mimeType) && !isImageFilename(fileHeader.Filename) {
		err := helpers.ErrorResponse{
			ErrorCode: http.StatusUnprocessableEntity,
			Message:   "Error on updating avatar.",
		}
		err.Errors = make(map[string][]string)
		err.Errors["AvatarFile"] = []string{"Invalid image type. Supported formats are jpg, jpeg and png"}
		helpers.NewErrorResponse(w, &err)

		return
	}

	bts, e := ioutil.ReadAll(file)
	if e != nil {
		err := helpers.ErrorResponse{
			ErrorCode: http.StatusUnprocessableEntity,
			Message:   "Error on updating avatar.",
		}
		err.Errors = make(map[string][]string)
		err.Errors["AvatarFile"] = []string{"Can't update avatar"}
		log.Println("Error on updating avatar,", e)
		helpers.NewErrorResponse(w, &err)

		return
	}
	os.MkdirAll("public", 0755)
	parts := strings.Split(fileHeader.Filename, ".")
	fileName := fmt.Sprintf("public/%v%v.%v", helpers.RandomString(10), helpers.RandomString(10), parts[len(parts)-1])
	e = ioutil.WriteFile(fileName, bts, 0666)
	if e != nil {
		err := helpers.ErrorResponse{
			ErrorCode: http.StatusUnprocessableEntity,
			Message:   "Error on updating avatar.",
		}
		err.Errors = make(map[string][]string)
		err.Errors["AvatarFile"] = []string{"Can't update avatar, IO Error"}
		log.Println("Error on updating avatar,", e)
		helpers.NewErrorResponse(w, &err)

		return
	}
	if user.Avatar != nil {
		oldFile := user.Avatar
		os.Remove(*oldFile)
	}
	user.Avatar = &fileName

	db.Model(&user).Save(&user)

	var response struct {
		Avatar string
	}
	response.Avatar = fileName

	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&response)
}

// isImageMimeType check mimeType string to be an image.
func isImageMimeType(mimeType string) bool {
	mimeType = strings.ToLower(mimeType)
	if strings.EqualFold(mimeType, "image/jpg") || strings.EqualFold(mimeType, "image/jpeg") || strings.EqualFold(mimeType, "image/png") {
		return true
	}

	return false
}

func isImageFilename(filename string) bool {
	filename = strings.ToLower(filename)
	is, e := regexp.MatchString(".+\\.(jpg$)|(png$)|(jpeg$)", filename)
	if e != nil {
		return false
	}

	return is
}
