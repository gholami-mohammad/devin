package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"

	"devin/database"
	"devin/middlewares"
	"devin/models"
)

func getValidUser(id uint64, isRoot bool) (user models.User, claim models.Claim, tokenString string) {
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

func getValidOrganization(id uint64, ownerID uint64) models.User {
	db := database.NewGORMInstance()
	defer db.Close()
	db.Exec(`delete from users where id=?;`, id)
	e := db.Exec(`insert into users (id, username, email, user_type, owner_id) values (?, ?, ?, 2, ?)`, id, fmt.Sprintf("org%v", id), fmt.Sprintf("org%v@gmail.com", id), ownerID).Error
	if e != nil {
		panic(e.Error())
	}

	var org models.User
	db.Where("id=?", id).First(&org)

	return org
}

func deleteTestOrganization(id uint64) {
	deleteTestUser(id)
}

func addUserToOrganization(userID, orgID uint64) {
	db := database.NewGORMInstance()
	defer db.Close()
	e := db.Exec(`insert into user_organization (user_id, organization_id, created_by_id) values (?, ?, ?)`, userID, orgID, userID).Error
	if e != nil {
		panic(e)
	}
}

func TestSave(t *testing.T) {
	_, _, tokenString := getValidUser(1, true)
	defer deleteTestUser(1)

	path := "/api/user/{id}/organization/save"
	route := mux.NewRouter()
	route.Use(middlewares.Authenticate)
	route.HandleFunc(path, Save)

	t.Run("Invalid content type", func(t *testing.T) {
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

	t.Run("Empty Request Body", func(t *testing.T) {
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

	t.Run("Authentication Failed", func(t *testing.T) {
		route := mux.NewRouter()
		route.HandleFunc(path, Save)

		req, _ := http.NewRequest(http.MethodPost, strings.Replace(path, "{id}", "1", 1), strings.NewReader(`{}`))
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

	t.Run("Invalid Request Body", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodPost, strings.Replace(path, "{id}", "1", 1), strings.NewReader(`Invalid Body`))
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

	t.Run("Unauthorized action", func(t *testing.T) {
		_, _, tokenString := getValidUser(100, false)
		defer deleteTestUser(100)
		req, _ := http.NewRequest(http.MethodPost, strings.Replace(path, "{id}", "1", 1), strings.NewReader(`{"Fullname": "Fake name"}`))
		req.Header.Add("Authorization", tokenString)
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

	t.Run("No User ID", func(t *testing.T) {
		noIDPath := strings.Replace(path, "/{id}", "", 1)
		route := mux.NewRouter()
		route.HandleFunc(noIDPath, Save)

		req, _ := http.NewRequest(http.MethodPost, noIDPath, strings.NewReader(`{"Fullname": "Fake name"}`))
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

	t.Run("Bad User ID - No integer value", func(t *testing.T) {
		route := mux.NewRouter()
		route.HandleFunc(path, Save)

		req, _ := http.NewRequest(http.MethodPost, strings.Replace(path, "{id}", "BAD_ID", 1), strings.NewReader(`{"Fullname": "Fake name"}`))
		req.Header.Add("Content-Type", "application/json")
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

	t.Run("Bad Username", func(t *testing.T) {
		route := mux.NewRouter()
		route.HandleFunc(path, Save)

		req, _ := http.NewRequest(http.MethodPost, strings.Replace(path, "{id}", "1", 1), strings.NewReader(`{"Username":"empty or invalid username", "Fullname": "Fake name"}`))
		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("Authorization", tokenString)
		rr := httptest.NewRecorder()
		route.ServeHTTP(rr, req)

		res := rr.Result()
		if res.StatusCode != http.StatusUnprocessableEntity {
			t.Fatal("Status code not matched. Response is", res.StatusCode)
		}
		defer res.Body.Close()
		bts, _ := ioutil.ReadAll(res.Body)
		if !strings.Contains(string(bts), "Invalid username") {
			t.Fatal("Invalid response message")
		}
	})

	t.Run("Bad email", func(t *testing.T) {
		route := mux.NewRouter()
		route.HandleFunc(path, Save)

		req, _ := http.NewRequest(http.MethodPost, strings.Replace(path, "{id}", "1", 1), strings.NewReader(`{"Username":"fake_username", "Email": "empty or invalid email", "Fullname": "Fake name"}`))
		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("Authorization", tokenString)
		rr := httptest.NewRecorder()
		route.ServeHTTP(rr, req)

		res := rr.Result()
		if res.StatusCode != http.StatusUnprocessableEntity {
			t.Fatal("Status code not matched. Response is", res.StatusCode)
		}
		defer res.Body.Close()
		bts, _ := ioutil.ReadAll(res.Body)
		if !strings.Contains(string(bts), "Invalid email") {
			t.Fatal("Invalid response message")
		}
	})

	t.Run("Duplicate email", func(t *testing.T) {
		getValidUser(30, true)
		defer deleteTestUser(30)
		req, _ := http.NewRequest(http.MethodPost, strings.Replace(path, "{id}", "1", 1), strings.NewReader(`{"Username": "fake_username", "Email": "m6devin30@gmail.com", "Fullname": "Fake name"}`))
		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("Authorization", tokenString)
		rr := httptest.NewRecorder()
		route.ServeHTTP(rr, req)

		res := rr.Result()
		if res.StatusCode != http.StatusUnprocessableEntity {
			t.Fatal("Status code not matched. Response is", res.StatusCode)
		}
		defer res.Body.Close()
		bts, _ := ioutil.ReadAll(res.Body)
		if !strings.Contains(string(bts), "This email is already registered") {
			t.Fatal("Invalid response message")
		}
	})

	t.Run("Duplicate username", func(t *testing.T) {
		getValidUser(30, true)
		defer deleteTestUser(30)
		req, _ := http.NewRequest(http.MethodPost, strings.Replace(path, "{id}", "1", 1), strings.NewReader(`{"Username": "mgh30", "Email": "fake_email@gmail.com", "Fullname": "Fake name"}`))
		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("Authorization", tokenString)
		rr := httptest.NewRecorder()
		route.ServeHTTP(rr, req)

		res := rr.Result()
		if res.StatusCode != http.StatusUnprocessableEntity {
			t.Fatal("Status code not matched. Response is", res.StatusCode)
		}
		defer res.Body.Close()
		bts, _ := ioutil.ReadAll(res.Body)
		if !strings.Contains(string(bts), "This username is already registered") {
			t.Fatal("Invalid response message")
		}
	})

	t.Run("OK", func(t *testing.T) {

		req, _ := http.NewRequest(http.MethodPost, strings.Replace(path, "{id}", "1", 1), strings.NewReader(`{"Username": "fake_username", "Email": "fake_email@example.com", "Fullname": "Fake name"}`))
		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("Authorization", tokenString)
		rr := httptest.NewRecorder()
		route.ServeHTTP(rr, req)

		res := rr.Result()
		if res.StatusCode != http.StatusOK {
			t.Fatal("Status code not matched. Response is", res.StatusCode)
		}
		defer res.Body.Close()
		bts, _ := ioutil.ReadAll(res.Body)
		if !strings.Contains(string(bts), "Fake name") {
			t.Fatal("Invalid response message")
		}
		var user models.User
		json.Unmarshal(bts, &user)
		deleteTestUser(user.ID)
	})
}

func TestUserOrganizationsIndex(t *testing.T) {
	_, _, tokenString := getValidUser(200, true)
	defer deleteTestUser(200)
	_, _, tokenString300 := getValidUser(300, true)
	defer deleteTestUser(300)

	getValidOrganization(201, 200)
	defer deleteTestOrganization(201)

	getValidOrganization(301, 300)
	defer deleteTestOrganization(300)

	addUserToOrganization(200, 301)

	t.Run("Bad URL UserID", func(t *testing.T) {
		path := "/api/user/organizations"
		route := mux.NewRouter()
		route.Use(middlewares.Authenticate)
		route.HandleFunc(path, UserOrganizationsIndex)

		req, _ := http.NewRequest(http.MethodGet, path, nil)
		req.Header.Add("Authorization", tokenString)
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

	t.Run("Bad URL UserID", func(t *testing.T) {
		path := "/api/user/{id}/organizations"
		route := mux.NewRouter()
		route.HandleFunc(path, UserOrganizationsIndex)

		req, _ := http.NewRequest(http.MethodGet, strings.Replace(path, "{id}", "200", 1), nil)

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

	t.Run("OK", func(t *testing.T) {
		path := "/api/user/{id}/organizations"
		route := mux.NewRouter()
		route.Use(middlewares.Authenticate)
		route.HandleFunc(path, UserOrganizationsIndex)

		req, _ := http.NewRequest(http.MethodGet, strings.Replace(path, "{id}", "200", 1), nil)
		req.Header.Add("Authorization", tokenString300)
		rr := httptest.NewRecorder()
		route.ServeHTTP(rr, req)

		res := rr.Result()
		if res.StatusCode != http.StatusForbidden {
			t.Fatal("Status code not matched. Response is", res.StatusCode)
		}
	})

	t.Run("OK", func(t *testing.T) {
		path := "/api/user/{id}/organizations"
		route := mux.NewRouter()
		route.Use(middlewares.Authenticate)
		route.HandleFunc(path, UserOrganizationsIndex)

		req, _ := http.NewRequest(http.MethodGet, strings.Replace(path, "{id}", "200", 1), nil)
		req.Header.Add("Authorization", tokenString)
		rr := httptest.NewRecorder()
		route.ServeHTTP(rr, req)

		res := rr.Result()
		if res.StatusCode != http.StatusOK {
			t.Fatal("Status code not matched. Response is", res.StatusCode)
		}
	})
}

func TestCanViewOrganizationsOfUser(t *testing.T) {
	type testItem struct {
		AuthUserID      uint64
		IsRoot          bool
		UserID          uint64
		RequestedResult bool
	}
	var testTable []testItem
	testTable = append(testTable, testItem{
		AuthUserID:      10,
		IsRoot:          false,
		UserID:          10,
		RequestedResult: true,
	})
	testTable = append(testTable, testItem{
		AuthUserID:      11,
		IsRoot:          false,
		UserID:          12,
		RequestedResult: false,
	})
	testTable = append(testTable, testItem{
		AuthUserID:      13,
		IsRoot:          true,
		UserID:          13,
		RequestedResult: true,
	})
	testTable = append(testTable, testItem{
		AuthUserID:      14,
		IsRoot:          true,
		UserID:          15,
		RequestedResult: true,
	})
	for _, x := range testTable {
		authUser, _, _ := getValidUser(x.AuthUserID, x.IsRoot)
		defer deleteTestUser(x.AuthUserID)
		result := canViewOrganizationsOfUser(httptest.NewRecorder(), authUser, x.UserID)
		if result != x.RequestedResult {
			t.Fatal(x)
		}
	}
}
