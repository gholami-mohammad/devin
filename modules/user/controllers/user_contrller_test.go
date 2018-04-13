package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
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

func deleteTestUser(id int) {
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
