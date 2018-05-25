package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/gorilla/mux"

	"devin/database"
	"devin/helpers"
	"devin/middlewares"
	"devin/models"
)

func getValidUser(id int, isRoot bool) (user models.User, claim models.Claim, tokenString string) {
	db := database.NewGORMInstance()
	defer db.Close()
	db.Exec(`delete from users where id=?;`, id)
	e := db.Exec(`insert into users (id, username, email, is_root_user) values (?, ?, ?, ?)`, id, fmt.Sprintf("mgh%v", id), fmt.Sprintf("m6devin%v@gmail.com", id), isRoot).Error
	if e != nil {
		panic(e.Error())
	}
	db.Where("id=?", id).First(&user)
	claim = user.GenerateNewTokenClaim()
	tokenString, _ = user.GenerateNewTokenString(claim)

	return user, claim, tokenString
}

func deleteTestUser(id uint64) {
	db := database.NewGORMInstance()
	defer db.Close()
	db.Exec(`delete from users where id=?;`, id)
}

func TestHandleProfileSharedErrors(t *testing.T) {
	_, _, tokenString := getValidUser(1, true)
	_, _, tokenStringStandardUser := getValidUser(2, false)
	defer deleteTestUser(1)
	defer deleteTestUser(2)

	route := mux.NewRouter()
	path := "/user/{id}/update"
	route.Handle(path, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		db := database.NewGORMInstance()
		defer db.Close()
		_, err := handleProfileSharedErrors(r, db)
		if err == nil {
			return
		}
		helpers.NewErrorResponse(w, err)
		return
	}))
	route.Use(middlewares.Authenticate)

	t.Run("bad content-type", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodPost, strings.Replace(path, "{id}", "1", 1), nil)
		req.Header.Add("Authorization", tokenString)
		rr := httptest.NewRecorder()
		route.ServeHTTP(rr, req)
		res := rr.Result()
		if res.StatusCode != http.StatusUnsupportedMediaType {
			t.Fatal("Status code not matched. Response is", res.StatusCode)
		}
		defer res.Body.Close()
		bts, _ := ioutil.ReadAll(res.Body)
		if !strings.Contains(string(bts), "Invalid content type") {
			t.Fatal("Invalid response message")
		}
	})

	t.Run("no id variable passed to the mux router", func(t *testing.T) {
		route := mux.NewRouter()
		path := "/user/update"
		route.Handle(path, http.HandlerFunc(UpdateProfile))
		req, _ := http.NewRequest(http.MethodPost, path, nil)

		req.Header.Add("Authorization", tokenString)
		req.Header.Add("Content-Type", "application/json")
		rr := httptest.NewRecorder()
		route.ServeHTTP(rr, req)
		res := rr.Result()
		if res.StatusCode != http.StatusUnprocessableEntity {
			t.Fatal("Status code not matched. Response is", res.StatusCode)
		}

		defer res.Body.Close()
		bts, _ := ioutil.ReadAll(res.Body)
		if !strings.Contains(string(bts), "Invalid User ID") {
			t.Fatal("Invalid response message")
		}
	})

	t.Run("invalid user_id data type", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodPost, strings.Replace(path, "{id}", "-1", 1), nil)

		req.Header.Add("Authorization", tokenString)
		req.Header.Add("Content-Type", "application/json")
		rr := httptest.NewRecorder()
		route.ServeHTTP(rr, req)
		res := rr.Result()

		if res.StatusCode != http.StatusUnprocessableEntity {
			t.Fatal("Status code not matched. Response is", res.StatusCode)
		}

		defer res.Body.Close()
		bts, _ := ioutil.ReadAll(res.Body)
		if !strings.Contains(string(bts), "Invalid User ID. Just integer values accepted") {
			t.Fatal("Invalid response message")
		}
	})

	t.Run("not exists user in DB", func(t *testing.T) {
		deleteTestUser(0)
		req, _ := http.NewRequest(http.MethodPost, strings.Replace(path, "{id}", "0", 1), nil)

		req.Header.Add("Authorization", tokenString)
		req.Header.Add("Content-Type", "application/json")
		rr := httptest.NewRecorder()
		route.ServeHTTP(rr, req)
		res := rr.Result()

		if res.StatusCode != http.StatusInternalServerError {
			t.Fatal("Status code not matched. Response is", res.StatusCode)
		}

		defer res.Body.Close()
		bts, _ := ioutil.ReadAll(res.Body)
		if !strings.Contains(string(bts), "Error on loading user data") {
			t.Fatal("Invalid response message")
		}
	})

	t.Run("Authorization token error", func(t *testing.T) {
		route := mux.NewRouter()
		route.Handle(path, http.HandlerFunc(UpdateProfile))
		req, _ := http.NewRequest(http.MethodPost, strings.Replace(path, "{id}", "1", 1), nil)

		req.Header.Add("Content-Type", "application/json")
		rr := httptest.NewRecorder()
		route.ServeHTTP(rr, req)
		res := rr.Result()

		if res.StatusCode != http.StatusUnauthorized {
			t.Fatal("Status code not matched. Response is", res.StatusCode)
		}

		defer res.Body.Close()
		bts, _ := ioutil.ReadAll(res.Body)
		if !strings.Contains(string(bts), "Auhtentication failed") {
			t.Fatal("Invalid response message")
		}
	})

	t.Run("Access denied", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodPost, strings.Replace(path, "{id}", "1", 1), nil)

		req.Header.Add("Authorization", tokenStringStandardUser)
		req.Header.Add("Content-Type", "application/json")
		rr := httptest.NewRecorder()
		route.ServeHTTP(rr, req)
		res := rr.Result()

		if res.StatusCode != http.StatusForbidden {
			t.Fatal("Status code not matched. Response is", res.StatusCode)
		}

		defer res.Body.Close()
		bts, _ := ioutil.ReadAll(res.Body)
		if !strings.Contains(string(bts), "This action is not allowed for you") {
			t.Fatal("Invalid response message")
		}
	})

	t.Run("Empty Request body", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodPost, strings.Replace(path, "{id}", "1", 1), nil)

		req.Header.Add("Authorization", tokenString)
		req.Header.Add("Content-Type", "application/json")
		rr := httptest.NewRecorder()
		route.ServeHTTP(rr, req)
		res := rr.Result()

		if res.StatusCode != http.StatusInternalServerError {
			t.Fatal("Status code not matched. Response is", res.StatusCode)
		}

		defer res.Body.Close()
		bts, _ := ioutil.ReadAll(res.Body)
		if !strings.Contains(string(bts), "Request body cant be empty") {
			t.Fatal("Invalid response message")
		}
	})

}

func TestUpdateProfile(t *testing.T) {
	user1, _, tokenString := getValidUser(1, true)
	// _, _, tokenStringStandardUser := getValidUser(2, false)
	defer deleteTestUser(1)
	// defer deleteTestUser(2)

	route := mux.NewRouter()
	path := "/user/{id}/update"
	route.Handle(path, http.HandlerFunc(UpdateProfile))
	route.Use(middlewares.Authenticate)

	t.Run("Invalid request body", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodPost, strings.Replace(path, "{id}", "1", 1), strings.NewReader("Bad Request body"))

		req.Header.Add("Authorization", tokenString)
		req.Header.Add("Content-Type", "application/json")
		rr := httptest.NewRecorder()
		route.ServeHTTP(rr, req)
		res := rr.Result()

		if res.StatusCode != http.StatusInternalServerError {
			t.Fatal("Status code not matched. Response is", res.StatusCode)
		}

		defer res.Body.Close()
		bts, _ := ioutil.ReadAll(res.Body)
		if !strings.Contains(string(bts), "Invalid request body") {
			t.Fatal("Invalid response message")
		}
	})

	t.Run("OK", func(t *testing.T) {
		sp := "Updated first name"
		user1.FirstName = &sp
		bts, _ := json.Marshal(&user1)
		req, _ := http.NewRequest(http.MethodPost, strings.Replace(path, "{id}", "1", 1), bytes.NewReader(bts))

		req.Header.Add("Authorization", tokenString)
		req.Header.Add("Content-Type", "application/json")
		rr := httptest.NewRecorder()
		route.ServeHTTP(rr, req)
		res := rr.Result()

		if res.StatusCode != http.StatusOK {
			t.Fatal("Status code not matched. Response is", res.StatusCode)
		}

		defer res.Body.Close()
		bts, _ = ioutil.ReadAll(res.Body)
		if !strings.Contains(string(bts), "Updated first name") {
			t.Fatal("Invalid response")
		}
	})

}

func TestWhoami(t *testing.T) {

	t.Run("Bad request context", func(t *testing.T) {
		rr := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/api/whoami", nil)
		handler := http.HandlerFunc(Whoami)
		handler.ServeHTTP(rr, req)

		res := rr.Result()
		if res.StatusCode != http.StatusUnauthorized {
			t.Fatal()
		}
	})

	t.Run("User not found", func(t *testing.T) {
		_, _, tokenString := getValidUser(14, true)
		deleteTestUser(14)

		rr := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/api/whoami", nil)
		req.Header.Add("Authorization", tokenString)
		handler := http.HandlerFunc(Whoami)
		handler.ServeHTTP(rr, req)

		res := rr.Result()
		if res.StatusCode != http.StatusUnauthorized {
			t.Fatal()
		}
	})

	t.Run("OK", func(t *testing.T) {
		_, _, tokenString := getValidUser(13, true)
		defer deleteTestUser(13)

		rr := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/api/whoami", nil)
		req.Header.Add("Authorization", tokenString)
		handler := http.HandlerFunc(Whoami)
		handler.ServeHTTP(rr, req)

		res := rr.Result()
		if res.StatusCode != http.StatusOK {
			t.Fatal()
		}
	})
}

func TestUpdateUsername(t *testing.T) {
	_, _, tokenString := getValidUser(1, true)
	getValidUser(2, false)
	defer deleteTestUser(1)
	defer deleteTestUser(2)

	route := mux.NewRouter()
	path := "/user/{id}/update_username"
	route.Handle(path, http.HandlerFunc(UpdateUsername))
	route.Use(middlewares.Authenticate)

	t.Run("Invalid request body", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodPost, strings.Replace(path, "{id}", "1", 1), strings.NewReader("Bad Request body"))

		req.Header.Add("Authorization", tokenString)
		req.Header.Add("Content-Type", "application/json")
		rr := httptest.NewRecorder()
		route.ServeHTTP(rr, req)
		res := rr.Result()

		if res.StatusCode != http.StatusInternalServerError {
			t.Fatal("Status code not matched. Response is", res.StatusCode)
		}

		defer res.Body.Close()
		bts, _ := ioutil.ReadAll(res.Body)
		if !strings.Contains(string(bts), "Invalid request body") {
			t.Fatal("Invalid response message")
		}
	})

	t.Run("Invalid username charachters", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodPost, strings.Replace(path, "{id}", "1", 1), strings.NewReader(`{"Username": "Bad**&&username"}`))

		req.Header.Add("Authorization", tokenString)
		req.Header.Add("Content-Type", "application/json")
		rr := httptest.NewRecorder()
		route.ServeHTTP(rr, req)
		res := rr.Result()

		if res.StatusCode != http.StatusUnprocessableEntity {
			t.Fatal("Status code not matched. Response is", res.StatusCode)
		}

		defer res.Body.Close()
		bts, _ := ioutil.ReadAll(res.Body)
		if !strings.Contains(string(bts), "invalid characters") {
			t.Fatal("Invalid response message")
		}
	})

	t.Run("Username already taken", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodPost, strings.Replace(path, "{id}", "1", 1), strings.NewReader(`{"Username": "mgh2"}`))

		req.Header.Add("Authorization", tokenString)
		req.Header.Add("Content-Type", "application/json")
		rr := httptest.NewRecorder()
		route.ServeHTTP(rr, req)
		res := rr.Result()

		if res.StatusCode != http.StatusUnprocessableEntity {
			t.Fatal("Status code not matched. Response is", res.StatusCode)
		}

		defer res.Body.Close()
		bts, _ := ioutil.ReadAll(res.Body)
		if !strings.Contains(string(bts), "already taken") {
			t.Fatal("Invalid response message")
		}
	})

	t.Run("OK", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodPost, strings.Replace(path, "{id}", "1", 1), strings.NewReader(`{"Username": "mgh1_updated"}`))

		req.Header.Add("Authorization", tokenString)
		req.Header.Add("Content-Type", "application/json")
		rr := httptest.NewRecorder()
		route.ServeHTTP(rr, req)
		res := rr.Result()

		if res.StatusCode != http.StatusOK {
			t.Fatal("Status code not matched. Response is", res.StatusCode)
		}

	})

}

func TestUpdateEmail(t *testing.T) {
	_, _, tokenString := getValidUser(1, true)
	getValidUser(2, false)
	defer deleteTestUser(1)
	defer deleteTestUser(2)

	route := mux.NewRouter()
	path := "/user/{id}/update_email"
	route.Handle(path, http.HandlerFunc(UpdateEmail))
	route.Use(middlewares.Authenticate)

	t.Run("Invalid request body", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodPost, strings.Replace(path, "{id}", "1", 1), strings.NewReader("Bad Request body"))

		req.Header.Add("Authorization", tokenString)
		req.Header.Add("Content-Type", "application/json")
		rr := httptest.NewRecorder()
		route.ServeHTTP(rr, req)
		res := rr.Result()

		if res.StatusCode != http.StatusInternalServerError {
			t.Fatal("Status code not matched. Response is", res.StatusCode)
		}

		defer res.Body.Close()
		bts, _ := ioutil.ReadAll(res.Body)
		if !strings.Contains(string(bts), "Invalid request body") {
			t.Fatal("Invalid response message")
		}
	})

	t.Run("Invalid email charachters", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodPost, strings.Replace(path, "{id}", "1", 1), strings.NewReader(`{"Email": "Bad**&&email"}`))

		req.Header.Add("Authorization", tokenString)
		req.Header.Add("Content-Type", "application/json")
		rr := httptest.NewRecorder()
		route.ServeHTTP(rr, req)
		res := rr.Result()

		if res.StatusCode != http.StatusUnprocessableEntity {
			t.Fatal("Status code not matched. Response is", res.StatusCode)
		}

		defer res.Body.Close()
		bts, _ := ioutil.ReadAll(res.Body)
		if !strings.Contains(string(bts), "invalid characters") {
			t.Fatal("Invalid response message")
		}
	})

	t.Run("Email already taken", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodPost, strings.Replace(path, "{id}", "1", 1), strings.NewReader(`{"Email": "m6devin2@gmail.com"}`))

		req.Header.Add("Authorization", tokenString)
		req.Header.Add("Content-Type", "application/json")
		rr := httptest.NewRecorder()
		route.ServeHTTP(rr, req)
		res := rr.Result()

		if res.StatusCode != http.StatusUnprocessableEntity {
			t.Fatal("Status code not matched. Response is", res.StatusCode)
		}

		defer res.Body.Close()
		bts, _ := ioutil.ReadAll(res.Body)
		if !strings.Contains(string(bts), "already taken") {
			t.Fatal("Invalid response message")
		}
	})

	t.Run("OK", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodPost, strings.Replace(path, "{id}", "1", 1), strings.NewReader(`{"Email": "m6devin1@updated.com"}`))

		req.Header.Add("Authorization", tokenString)
		req.Header.Add("Content-Type", "application/json")
		rr := httptest.NewRecorder()
		route.ServeHTTP(rr, req)
		res := rr.Result()

		if res.StatusCode != http.StatusOK {
			t.Fatal("Status code not matched. Response is", res.StatusCode)
		}

	})

}

func TestProfileBasicInfo(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(ProfileBasicInfo))
	defer server.Close()

	res, e := http.Get(server.URL)
	if e != nil {
		t.Fatal(e)
	}
	defer res.Body.Close()
	bts, _ := ioutil.ReadAll(res.Body)
	str := string(bts)

	if !strings.Contains(str, "LocalizationLanguages") || !strings.Contains(str, "DateFormats") || !strings.Contains(str, "TimeFormats") || !strings.Contains(str, "CalendarSystems") || !strings.Contains(str, "OfficePhoneCountryCodes") || !strings.Contains(str, "HomePhoneCountryCodes") || !strings.Contains(str, "CellPhoneCountryCodes") || !strings.Contains(str, "FaxCountryCodes") || !strings.Contains(str, "Countries") {
		t.Fatal("Response dose not contains all keys")
	}
}

func TestUpdatePassword(t *testing.T) {
	_, _, tokenString := getValidUser(1, true)
	defer deleteTestUser(1)

	route := mux.NewRouter()
	path := "/user/{id}/update_username"
	route.Handle(path, http.HandlerFunc(UpdatePassword))
	route.Use(middlewares.Authenticate)

	t.Run("Invalid request body", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodPost, strings.Replace(path, "{id}", "1", 1), strings.NewReader("Bad Request body"))

		req.Header.Add("Authorization", tokenString)
		req.Header.Add("Content-Type", "application/json")
		rr := httptest.NewRecorder()
		route.ServeHTTP(rr, req)
		res := rr.Result()

		if res.StatusCode != http.StatusInternalServerError {
			t.Fatal("Status code not matched. Response is", res.StatusCode)
		}

		defer res.Body.Close()
		bts, _ := ioutil.ReadAll(res.Body)
		if !strings.Contains(string(bts), "Invalid request body") {
			t.Fatal("Invalid response message")
		}
	})

	t.Run("Small password length", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodPost, strings.Replace(path, "{id}", "1", 1), strings.NewReader(`{"Password": "min", "PasswordVerification": "min"}`))

		req.Header.Add("Authorization", tokenString)
		req.Header.Add("Content-Type", "application/json")
		rr := httptest.NewRecorder()
		route.ServeHTTP(rr, req)
		res := rr.Result()

		if res.StatusCode != http.StatusUnprocessableEntity {
			t.Fatal("Status code not matched. Response is", res.StatusCode)
		}

		defer res.Body.Close()
		bts, _ := ioutil.ReadAll(res.Body)
		if !strings.Contains(string(bts), "at least 6 characters") {
			t.Fatal("Invalid response message")
		}
	})

	t.Run("Verification not match", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodPost, strings.Replace(path, "{id}", "1", 1), strings.NewReader(`{"Password": "abcdef", "PasswordVerification": "ghijke"}`))

		req.Header.Add("Authorization", tokenString)
		req.Header.Add("Content-Type", "application/json")
		rr := httptest.NewRecorder()
		route.ServeHTTP(rr, req)
		res := rr.Result()

		if res.StatusCode != http.StatusUnprocessableEntity {
			t.Fatal("Status code not matched. Response is", res.StatusCode)
		}

		defer res.Body.Close()
		bts, _ := ioutil.ReadAll(res.Body)
		if !strings.Contains(string(bts), "Password verification does not match") {
			t.Fatal("Invalid response message")
		}
	})

	t.Run("OK", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodPost, strings.Replace(path, "{id}", "1", 1), strings.NewReader(`{"Password": "abcdef", "PasswordVerification": "abcdef"}`))

		req.Header.Add("Authorization", tokenString)
		req.Header.Add("Content-Type", "application/json")
		rr := httptest.NewRecorder()
		route.ServeHTTP(rr, req)
		res := rr.Result()

		if res.StatusCode != http.StatusOK {
			t.Fatal("Status code not matched. Response is", res.StatusCode)
		}

		defer res.Body.Close()
		bts, _ := ioutil.ReadAll(res.Body)
		if !strings.Contains(string(bts), "Password updated") {
			t.Fatal("Invalid response message")
		}
	})
}

func TestUpdateAvatar(t *testing.T) {
	_, _, tokenString := getValidUser(1, true)
	defer deleteTestUser(1)

	route := mux.NewRouter()
	path := "/user/{id}/update_avatar"
	route.Handle(path, http.HandlerFunc(UpdateAvatar))
	route.Use(middlewares.Authenticate)

	t.Run("Authentication faild", func(t *testing.T) {
		route := mux.NewRouter()
		path := "/user/{id}/update_avatar"
		route.Handle(path, http.HandlerFunc(UpdateAvatar))
		req, e := http.NewRequest(http.MethodPost, strings.Replace(path, "{id}", "1", 1), nil)
		if e != nil {
			t.Fatal(e)
		}

		rr := httptest.NewRecorder()
		route.ServeHTTP(rr, req)
		res := rr.Result()
		if res.StatusCode != http.StatusUnauthorized {
			t.Fatal("Status code not matched. Response is", res.StatusCode)
		}

		defer res.Body.Close()
		bts, _ := ioutil.ReadAll(res.Body)
		if !strings.Contains(string(bts), "Auhtentication failed") {
			t.Fatal("Invalid response message")
		}

	})

	t.Run("No userID in route", func(t *testing.T) {
		route := mux.NewRouter()
		path := "/user/update_avatar"
		route.Handle(path, http.HandlerFunc(UpdateAvatar))
		req, e := http.NewRequest(http.MethodPost, path, nil)
		req.Header.Add("Authorization", tokenString)
		if e != nil {
			t.Fatal(e)
		}

		rr := httptest.NewRecorder()
		route.ServeHTTP(rr, req)
		res := rr.Result()
		if res.StatusCode != http.StatusUnprocessableEntity {
			t.Fatal("Status code not matched. Response is", res.StatusCode)
		}

		defer res.Body.Close()
		bts, _ := ioutil.ReadAll(res.Body)
		if !strings.Contains(string(bts), "Invalid User ID") {
			t.Fatal("Invalid response message")
		}
	})

	t.Run("Bad userID in route, not integer", func(t *testing.T) {
		route := mux.NewRouter()
		path := "/user/{id}/pdate_avatar"
		route.Handle(path, http.HandlerFunc(UpdateAvatar))
		req, e := http.NewRequest(http.MethodPost, strings.Replace(path, "{id}", "INVALID", 1), nil)
		req.Header.Add("Authorization", tokenString)
		if e != nil {
			t.Fatal(e)
		}

		rr := httptest.NewRecorder()
		route.ServeHTTP(rr, req)
		res := rr.Result()
		if res.StatusCode != http.StatusUnprocessableEntity {
			t.Fatal("Status code not matched. Response is", res.StatusCode)
		}

		defer res.Body.Close()
		bts, _ := ioutil.ReadAll(res.Body)
		if !strings.Contains(string(bts), "integer values accepted") {
			t.Fatal("Invalid response message")
		}
	})

	t.Run("User Not Found", func(t *testing.T) {
		_, _, token2 := getValidUser(2, true)
		deleteTestUser(2)
		route := mux.NewRouter()
		path := "/user/{id}/pdate_avatar"
		route.Handle(path, http.HandlerFunc(UpdateAvatar))
		req, e := http.NewRequest(http.MethodPost, strings.Replace(path, "{id}", "2", 1), nil)
		req.Header.Add("Authorization", token2)
		if e != nil {
			t.Fatal(e)
		}

		rr := httptest.NewRecorder()
		route.ServeHTTP(rr, req)
		res := rr.Result()
		if res.StatusCode != http.StatusUnauthorized {
			t.Fatal("Status code not matched. Response is", res.StatusCode)
		}

		defer res.Body.Close()
		bts, _ := ioutil.ReadAll(res.Body)
		if !strings.Contains(string(bts), "Auhtentication failed") {
			t.Fatal("Invalid response message")
		}
	})

	t.Run("Permission Denied", func(t *testing.T) {
		_, _, token3 := getValidUser(3, false)
		defer deleteTestUser(3)
		route := mux.NewRouter()
		path := "/user/{id}/pdate_avatar"
		route.Handle(path, http.HandlerFunc(UpdateAvatar))
		req, e := http.NewRequest(http.MethodPost, strings.Replace(path, "{id}", "1", 1), nil)
		req.Header.Add("Authorization", token3)
		if e != nil {
			t.Fatal(e)
		}

		rr := httptest.NewRecorder()
		route.ServeHTTP(rr, req)
		res := rr.Result()
		if res.StatusCode != http.StatusForbidden {
			t.Fatal("Status code not matched. Response is", res.StatusCode)
		}

		defer res.Body.Close()
		bts, _ := ioutil.ReadAll(res.Body)
		if !strings.Contains(string(bts), "not allowed") {
			t.Fatal("Invalid response message")
		}
	})

	t.Run("Error Parse Form", func(t *testing.T) {

		req, e := http.NewRequest(http.MethodPost, strings.Replace(path, "{id}", "1", 1), nil)
		req.Header.Add("Authorization", tokenString)
		if e != nil {
			t.Fatal(e)
		}

		rr := httptest.NewRecorder()
		route.ServeHTTP(rr, req)
		res := rr.Result()
		if res.StatusCode != http.StatusUnprocessableEntity {
			t.Fatal("Status code not matched. Response is", res.StatusCode)
		}

		defer res.Body.Close()
		bts, _ := ioutil.ReadAll(res.Body)
		if !strings.Contains(string(bts), "Can't read file") {
			t.Fatal("Invalid response message")
		}
	})

	t.Run("Invalid image type", func(t *testing.T) {

		file, e := os.Open("../../../test_files/text_file.txt")
		if e != nil {
			t.Fatal(e)
		}
		defer file.Close()

		fi, e := file.Stat()
		if e != nil {
			t.Fatal(e)
		}

		fileContent, e := ioutil.ReadAll(file)
		if e != nil {
			t.Fatal(e)
		}

		writer := &bytes.Buffer{}

		mw := multipart.NewWriter(writer)
		part, e := mw.CreateFormFile("AvatarFile", fi.Name())
		if e != nil {
			t.Fatal(e)
		}
		part.Write(fileContent)
		mw.WriteField("AvatarFile", fi.Name())
		req, e := http.NewRequest(http.MethodPost, strings.Replace(path, "{id}", "1", 1), writer)
		if e != nil {
			t.Fatal(e)
		}
		req.Header.Add("Content-Type", mw.FormDataContentType())
		req.Header.Add("Authorization", tokenString)
		mw.Close()

		rr := httptest.NewRecorder()
		route.ServeHTTP(rr, req)
		res := rr.Result()
		if res.StatusCode != http.StatusUnprocessableEntity {
			t.Fatal("Status code not matched. Response is", res.StatusCode)
		}

		defer res.Body.Close()
		bts, _ := ioutil.ReadAll(res.Body)
		if !strings.Contains(string(bts), "Invalid image type") {
			t.Fatal("Invalid response message")
		}
	})

	t.Run("OK", func(t *testing.T) {

		file, e := os.Open("../../../test_files/golang.png")
		if e != nil {
			t.Fatal(e)
		}
		defer file.Close()

		fi, e := file.Stat()
		if e != nil {
			t.Fatal(e)
		}

		fileContent, e := ioutil.ReadAll(file)
		if e != nil {
			t.Fatal(e)
		}

		writer := &bytes.Buffer{}
		mw := multipart.NewWriter(writer)
		part, e := mw.CreateFormFile("AvatarFile", fi.Name())
		if e != nil {
			t.Fatal(e)
		}
		part.Write(fileContent)
		mw.WriteField("AvatarFile", fi.Name())

		req, e := http.NewRequest(http.MethodPost, strings.Replace(path, "{id}", "1", 1), writer)
		if e != nil {
			t.Fatal(e)
		}
		req.Header.Add("Content-Type", mw.FormDataContentType())
		req.Header.Add("Authorization", tokenString)

		mw.Close()

		rr := httptest.NewRecorder()
		route.ServeHTTP(rr, req)
		res := rr.Result()
		if res.StatusCode != http.StatusOK {
			t.Fatal("Status code not matched. Response is", res.StatusCode)
		}

		defer res.Body.Close()
		bts, _ := ioutil.ReadAll(res.Body)
		if !strings.Contains(string(bts), "Avatar") {
			t.Log(string(bts))
			t.Fatal("Invalid response message")
		}
	})

}

func TestIsImageMimeType(t *testing.T) {
	testTable := make(map[string]bool, 4)
	testTable["image/jpeg"] = true
	testTable["image/jpg"] = true
	testTable["image/png"] = true
	testTable["image/other"] = false
	for k, v := range testTable {
		if isImageMimeType(k) != v {
			t.Fail()
		}
	}
}

func TestIsImageFilename(t *testing.T) {
	filenames := make(map[string]bool, 5)
	filenames["a.png"] = true
	filenames["a.jpeg"] = true
	filenames["a.jpg"] = true
	filenames["a.pdf"] = false
	filenames["a.some*"] = false
	for k, v := range filenames {
		if isImageFilename(k) != v {
			t.Fail()
		}
	}
}
