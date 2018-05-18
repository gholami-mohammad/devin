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
	"devin/middlewares"
	"devin/models"
)

var testID uint64 = 10000

func getTestID() uint64 {
	testID += 1
	return testID
}

func TestInviteUser(t *testing.T) {
	id1 := getTestID()
	id2 := getTestID()
	id3 := getTestID()
	_, _, tokenString := getValidUser(id1, true)
	defer deleteTestUser(id1)

	getValidUser(id3, false)
	defer deleteTestUser(id3)

	getValidOrganization(id2, id1)
	defer deleteTestOrganization(id2)

	path := "/api/organization/{id}/invite_user"
	route := mux.NewRouter()
	route.Use(middlewares.Authenticate)
	route.HandleFunc(path, InviteUser)

	t.Run("Bad Content Type", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodPost, strings.Replace(path, "{id}", fmt.Sprintf("%v", id2), 1), nil)
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

	t.Run("No Orgaanization ID", func(t *testing.T) {
		path := "/api/organization/invite_user"
		route := mux.NewRouter()
		route.Use(middlewares.Authenticate)
		route.HandleFunc(path, InviteUser)

		req, _ := http.NewRequest(http.MethodPost, path, strings.NewReader("{}"))
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
		if !strings.Contains(string(bts), "Invalid Organization ID") {
			t.Fatal("Invalid response message")
		}
	})

	t.Run("Invalid Organization ID data type", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodPost, strings.Replace(path, "{id}", "ID", 1), strings.NewReader(`{}`))
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
		if !strings.Contains(string(bts), "integer values accepted") {
			t.Fatal("Invalid response message")
		}
	})

	t.Run("Organization Not Found", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodPost, strings.Replace(path, "{id}", "1000099900000", 1), strings.NewReader(`{}`))
		req.Header.Add("Authorization", tokenString)
		req.Header.Add("Content-Type", "application/json")
		rr := httptest.NewRecorder()
		route.ServeHTTP(rr, req)

		res := rr.Result()
		if res.StatusCode != http.StatusNotFound {
			t.Fatal("Status code not matched. Response is", res.StatusCode)
		}
		defer res.Body.Close()
		bts, _ := ioutil.ReadAll(res.Body)
		if !strings.Contains(string(bts), "not found") {
			t.Fatal("Invalid response message")
		}
	})

	t.Run("Authentication Failed", func(t *testing.T) {
		route := mux.NewRouter()
		route.HandleFunc(path, InviteUser)

		req, _ := http.NewRequest(http.MethodPost, strings.Replace(path, "{id}", fmt.Sprintf("%v", id2), 1), strings.NewReader(`{}`))
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

	t.Run("Empty Request body", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodPost, strings.Replace(path, "{id}", fmt.Sprintf("%v", id2), 1), nil)
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

	t.Run("Bad Request body", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodPost, strings.Replace(path, "{id}", fmt.Sprintf("%v", id2), 1), strings.NewReader("Bad Content"))
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

	t.Run("Permission Denied", func(t *testing.T) {
		getValidUser(2, true)
		_, _, tokenString3 := getValidUser(3, false)
		defer deleteTestUser(3)
		getValidOrganization(102, 2)
		defer deleteTestUser(2)
		defer deleteTestOrganization(102)

		req, _ := http.NewRequest(http.MethodPost, strings.Replace(path, "{id}", "102", 1), strings.NewReader(`{"Identifier":"mgh2"}`))
		req.Header.Add("Authorization", tokenString3)
		req.Header.Add("Content-Type", "application/json")
		rr := httptest.NewRecorder()
		route.ServeHTTP(rr, req)

		res := rr.Result()

		if res.StatusCode != http.StatusForbidden {
			t.Fatal("Status code not matched. Response is", res.StatusCode)
		}
		defer res.Body.Close()
		bts, _ := ioutil.ReadAll(res.Body)
		if !strings.Contains(string(bts), "not permitted") {
			t.Fatal("Invalid response message")
		}
	})

	t.Run("Empty Identifier", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodPost, strings.Replace(path, "{id}", fmt.Sprintf("%v", id2), 1), strings.NewReader(`{}`))
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
		if !strings.Contains(string(bts), "Username or email is required") {
			t.Fatal("Invalid response message")
		}
	})

	t.Run("invited user not found", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodPost, strings.Replace(path, "{id}", fmt.Sprintf("%v", id2), 1), strings.NewReader(`{"Identifier":"notfound_user"}`))
		req.Header.Add("Authorization", tokenString)
		req.Header.Add("Content-Type", "application/json")
		rr := httptest.NewRecorder()
		route.ServeHTTP(rr, req)

		res := rr.Result()

		if res.StatusCode != http.StatusNotFound {
			t.Fatal("Status code not matched. Response is", res.StatusCode)
		}
		defer res.Body.Close()
		bts, _ := ioutil.ReadAll(res.Body)
		if !strings.Contains(string(bts), "Not found") {
			t.Fatal("Invalid response message")
		}
	})

	t.Run("User is member of organization", func(t *testing.T) {
		var userID uint64 = getTestID()
		var orgID uint64 = getTestID()
		_, _, ts := getValidUser(userID, true)
		getValidOrganization(orgID, userID)
		defer deleteTestUser(userID)
		defer deleteTestOrganization(orgID)

		obj := models.UserOrganization{}
		obj.CreatedByID = userID
		obj.OrganizationID = &orgID
		obj.UserID = &userID
		obj.IsAdminOfOrganization = true
		db := database.NewGORMInstance()
		defer db.Close()

		db.Create(&obj)

		req, _ := http.NewRequest(http.MethodPost, strings.Replace(path, "{id}", fmt.Sprintf("%v", orgID), 1), strings.NewReader(fmt.Sprintf(`{"Identifier": "mgh%v"}`, userID)))
		req.Header.Add("Authorization", ts)
		req.Header.Add("Content-Type", "application/json")
		rr := httptest.NewRecorder()
		route.ServeHTTP(rr, req)

		res := rr.Result()

		if res.StatusCode != http.StatusUnprocessableEntity {
			t.Fatal("Status code not matched. Response is", res.StatusCode)
		}
		defer res.Body.Close()
		bts, _ := ioutil.ReadAll(res.Body)
		if !strings.Contains(string(bts), "User exists") {
			t.Fatal("Invalid response message", string(bts))
		}
	})

	t.Run("User already invited", func(t *testing.T) {
		var userID uint64 = getTestID() //80
		var orgID uint64 = getTestID()  //81
		_, _, ts := getValidUser(userID, true)
		getValidOrganization(orgID, userID)
		defer deleteTestUser(userID)
		defer deleteTestOrganization(orgID)

		obj := models.UserOrganizationInvitation{}
		obj.CreatedByID = userID
		obj.OrganizationID = orgID
		obj.UserID = &userID
		db := database.NewGORMInstance()
		defer db.Close()

		db.Create(&obj)

		req, _ := http.NewRequest(http.MethodPost, strings.Replace(path, "{id}", fmt.Sprintf("%v", orgID), 1), strings.NewReader(fmt.Sprintf(`{"Identifier": "mgh%v"}`, userID)))
		req.Header.Add("Authorization", ts)
		req.Header.Add("Content-Type", "application/json")
		rr := httptest.NewRecorder()
		route.ServeHTTP(rr, req)

		res := rr.Result()

		if res.StatusCode != http.StatusUnprocessableEntity {
			t.Fatal("Status code not matched. Response is", res.StatusCode)
		}
		defer res.Body.Close()
		bts, _ := ioutil.ReadAll(res.Body)
		if !strings.Contains(string(bts), "User already invited") {
			t.Fatal("Invalid response message")
		}
	})

	t.Run("OK", func(t *testing.T) {
		id1 := getTestID()
		id2 := getTestID()
		_, _, ts := getValidUser(id1, true)
		getValidOrganization(id2, id1)
		defer deleteTestUser(id1)
		defer deleteTestOrganization(id2)
		req, _ := http.NewRequest(http.MethodPost, strings.Replace(path, "{id}", fmt.Sprintf("%v", id2), 1), strings.NewReader(fmt.Sprintf(`{"Identifier": "mgh%v"}`, id1)))
		req.Header.Add("Authorization", ts)
		req.Header.Add("Content-Type", "application/json")
		rr := httptest.NewRecorder()
		route.ServeHTTP(rr, req)

		res := rr.Result()

		defer res.Body.Close()
		if res.StatusCode != http.StatusOK {
			bts, _ := ioutil.ReadAll(res.Body)
			t.Log(string(bts))
			t.Fatal("Status code not matched. Response is", res.StatusCode)
		}
	})
}

func createPendingInvitation(userID, creatorID, orgID uint64) (id uint64) {
	db := database.NewGORMInstance()
	defer db.Close()

	var IDStruct struct {
		ID uint64
	}
	db.Raw(`insert into user_organization_invitations
        (user_id, organization_id, created_by_id)  values
        (?, ?, ?)
        returning id`, userID, orgID, creatorID).Scan(&IDStruct)
	return IDStruct.ID
}

func TestPendingInvitationRequests(t *testing.T) {
	id1 := getTestID()
	id2 := getTestID()
	id3 := getTestID()
	getValidUser(id1, true)
	defer deleteTestUser(id1)
	_, _, token := getValidUser(id2, true)
	defer deleteTestUser(id2)
	getValidOrganization(id3, id1)
	defer deleteTestOrganization(id3)
	createPendingInvitation(id2, id1, id3)

	path := "/api/user/{id}/pending_invitations"

	t.Run("OK", func(t *testing.T) {
		route := mux.NewRouter()
		route.Use(middlewares.Authenticate)
		route.HandleFunc(path, PendingInvitationRequests)

		req, e := http.NewRequest(http.MethodGet, strings.Replace(path, "{id}", fmt.Sprintf("%v", id2), 1), nil)
		if e != nil {
			t.Fatal(e)
		}

		req.Header.Add("Authorization", token)
		rr := httptest.NewRecorder()
		route.ServeHTTP(rr, req)

		res := rr.Result()

		defer res.Body.Close()
		bts, _ := ioutil.ReadAll(res.Body)
		var items []models.UserOrganizationInvitation
		e = json.Unmarshal(bts, &items)
		if e != nil {
			t.Fatal(e)
		}

		if len(items) < 1 {
			t.Fatal("Incorrect response count")
		}
	})

	t.Run("Authentication faild", func(t *testing.T) {
		route := mux.NewRouter()
		route.HandleFunc(path, PendingInvitationRequests)

		req, e := http.NewRequest(http.MethodGet, strings.Replace(path, "{id}", fmt.Sprintf("%v", id2), 1), nil)
		if e != nil {
			t.Fatal(e)
		}

		rr := httptest.NewRecorder()
		route.ServeHTTP(rr, req)

		res := rr.Result()

		defer res.Body.Close()
		bts, _ := ioutil.ReadAll(res.Body)
		if !bytes.Contains(bts, []byte("Auhtentication failed")) {
			t.Fatal("Invalid response message")
		}
	})

	t.Run("Invalid ID in URL", func(t *testing.T) {
		route := mux.NewRouter()
		route.HandleFunc("/", PendingInvitationRequests)

		req, e := http.NewRequest(http.MethodGet, "/", nil)
		if e != nil {
			t.Fatal(e)
		}
		req.Header.Add("Authorization", token)

		rr := httptest.NewRecorder()
		route.ServeHTTP(rr, req)

		res := rr.Result()

		defer res.Body.Close()
		bts, _ := ioutil.ReadAll(res.Body)
		if !bytes.Contains(bts, []byte("Invalid ID")) {

			t.Fatal("Invalid response message", string(bts))
		}
	})

	t.Run("Requested user not exist", func(t *testing.T) {
		id4 := getTestID()
		deleteTestUser(id4)
		route := mux.NewRouter()
		route.Use(middlewares.Authenticate)
		route.HandleFunc(path, PendingInvitationRequests)

		req, e := http.NewRequest(http.MethodGet, strings.Replace(path, "{id}", fmt.Sprintf("%v", id4), 1), nil)
		if e != nil {
			t.Fatal(e)
		}

		req.Header.Add("Authorization", token)
		rr := httptest.NewRecorder()
		route.ServeHTTP(rr, req)

		res := rr.Result()

		defer res.Body.Close()
		bts, _ := ioutil.ReadAll(res.Body)
		if !bytes.Contains(bts, []byte("User not found")) {

			t.Fatal("Invalid response message", string(bts))
		}
	})

	t.Run("Permission denied", func(t *testing.T) {
		id5 := getTestID()
		id6 := getTestID()
		_, _, tokenId5 := getValidUser(id5, false)
		getValidUser(id6, false)
		defer deleteTestUser(id5)
		defer deleteTestUser(id6)
		route := mux.NewRouter()
		route.Use(middlewares.Authenticate)
		route.HandleFunc(path, PendingInvitationRequests)

		req, e := http.NewRequest(http.MethodGet, strings.Replace(path, "{id}", fmt.Sprintf("%v", id6), 1), nil)
		if e != nil {
			t.Fatal(e)
		}

		req.Header.Add("Authorization", tokenId5)
		rr := httptest.NewRecorder()
		route.ServeHTTP(rr, req)

		res := rr.Result()

		defer res.Body.Close()
		bts, _ := ioutil.ReadAll(res.Body)
		if !bytes.Contains(bts, []byte("not permitted")) {

			t.Fatal("Invalid response message", string(bts))
		}
	})

}

func TestAcceptOrRejectInvitation(t *testing.T) {
	path := "/{id}/{acceptance_status}"
	id1 := getTestID()
	id2 := getTestID()
	id3 := getTestID()
	getValidUser(id1, false)
	defer deleteTestUser(id1)
	_, _, token := getValidUser(id3, false)
	defer deleteTestUser(id3)
	getValidOrganization(id2, id1)
	defer deleteTestOrganization(id2)

	t.Run("OK_Accept", func(t *testing.T) {
		invID := createPendingInvitation(id3, id1, id2)
		route := mux.NewRouter()
		route.Use(middlewares.Authenticate)
		route.HandleFunc(path, AcceptOrRejectInvitation)

		req, e := http.NewRequest(http.MethodGet, fmt.Sprintf("/%v/%v", invID, "accept"), nil)
		if e != nil {
			t.Fatal(e)
		}
		req.Header.Add("Authorization", token)
		rr := httptest.NewRecorder()
		route.ServeHTTP(rr, req)

		res := rr.Result()
		bts, e := ioutil.ReadAll(res.Body)
		if e != nil {
			t.Fatal(e)
		}
		defer res.Body.Close()

		if res.StatusCode != http.StatusOK {
			t.Fatal("Status code not matched. Response is", res.StatusCode)
		}

		if !bytes.Contains(bts, []byte("updated")) {
			t.Fatal("Invalid response message", string(bts))
		}
	})

	t.Run("OK_Reject", func(t *testing.T) {
		invID := createPendingInvitation(id3, id1, id2)
		route := mux.NewRouter()
		route.Use(middlewares.Authenticate)
		route.HandleFunc(path, AcceptOrRejectInvitation)

		req, e := http.NewRequest(http.MethodGet, fmt.Sprintf("/%v/%v", invID, "reject"), nil)
		if e != nil {
			t.Fatal(e)
		}
		req.Header.Add("Authorization", token)
		rr := httptest.NewRecorder()
		route.ServeHTTP(rr, req)

		res := rr.Result()
		bts, e := ioutil.ReadAll(res.Body)
		if e != nil {
			t.Fatal(e)
		}
		defer res.Body.Close()

		if res.StatusCode != http.StatusOK {
			t.Fatal("Status code not matched. Response is", res.StatusCode)
		}

		if !bytes.Contains(bts, []byte("updated")) {
			t.Fatal("Invalid response message", string(bts))
		}
	})

	t.Run("Authentication failed", func(t *testing.T) {
		invID := createPendingInvitation(id3, id1, id2)
		route := mux.NewRouter()
		route.HandleFunc(path, AcceptOrRejectInvitation)

		req, e := http.NewRequest(http.MethodGet, fmt.Sprintf("/%v/%v", invID, "reject"), nil)
		if e != nil {
			t.Fatal(e)
		}
		rr := httptest.NewRecorder()
		route.ServeHTTP(rr, req)

		res := rr.Result()
		bts, e := ioutil.ReadAll(res.Body)
		if e != nil {
			t.Fatal(e)
		}
		defer res.Body.Close()

		if res.StatusCode != http.StatusUnauthorized {
			t.Fatal("Status code not matched. Response is", res.StatusCode)
		}

		if !bytes.Contains(bts, []byte("Auhtentication failed")) {
			t.Fatal("Invalid response message", string(bts))
		}
	})

	t.Run("No Acceptance Status", func(t *testing.T) {
		invID := createPendingInvitation(id3, id1, id2)
		route := mux.NewRouter()
		route.HandleFunc("/{id}", AcceptOrRejectInvitation)
		route.Use(middlewares.Authenticate)

		req, e := http.NewRequest(http.MethodGet, fmt.Sprintf("/%v", invID), nil)
		if e != nil {
			t.Fatal(e)
		}
		req.Header.Add("Authorization", token)

		rr := httptest.NewRecorder()
		route.ServeHTTP(rr, req)

		res := rr.Result()
		bts, e := ioutil.ReadAll(res.Body)
		if e != nil {
			t.Fatal(e)
		}
		defer res.Body.Close()

		if res.StatusCode != http.StatusUnprocessableEntity {
			t.Fatal("Status code not matched. Response is", res.StatusCode)
		}

		if !bytes.Contains(bts, []byte("Invalid acceptance status")) {
			t.Fatal("Invalid response message", string(bts))
		}
	})

	t.Run("Bad Acceptance Status", func(t *testing.T) {
		invID := createPendingInvitation(id3, id1, id2)
		route := mux.NewRouter()
		route.HandleFunc(path, AcceptOrRejectInvitation)
		route.Use(middlewares.Authenticate)

		req, e := http.NewRequest(http.MethodGet, fmt.Sprintf("/%v/%v", invID, "invalid"), nil)
		if e != nil {
			t.Fatal(e)
		}
		req.Header.Add("Authorization", token)

		rr := httptest.NewRecorder()
		route.ServeHTTP(rr, req)

		res := rr.Result()
		bts, e := ioutil.ReadAll(res.Body)
		if e != nil {
			t.Fatal(e)
		}
		defer res.Body.Close()

		if res.StatusCode != http.StatusUnprocessableEntity {
			t.Fatal("Status code not matched. Response is", res.StatusCode)
		}

		if !bytes.Contains(bts, []byte("Invalid acceptance status")) {
			t.Fatal("Invalid response message", string(bts))
		}
	})

	t.Run("No ID passed in url", func(t *testing.T) {
		route := mux.NewRouter()
		route.HandleFunc("/{acceptance_status}", AcceptOrRejectInvitation)
		route.Use(middlewares.Authenticate)

		req, e := http.NewRequest(http.MethodGet, "/accept", nil)
		if e != nil {
			t.Fatal(e)
		}
		req.Header.Add("Authorization", token)

		rr := httptest.NewRecorder()
		route.ServeHTTP(rr, req)

		res := rr.Result()
		bts, e := ioutil.ReadAll(res.Body)
		if e != nil {
			t.Fatal(e)
		}
		defer res.Body.Close()

		if res.StatusCode != http.StatusUnprocessableEntity {
			t.Fatal("Status code not matched. Response is", res.StatusCode)
		}

		if !bytes.Contains(bts, []byte("Invalid ID")) {
			t.Fatal("Invalid response message", string(bts))
		}
	})

	t.Run("Invitation not found", func(t *testing.T) {

		route := mux.NewRouter()
		route.HandleFunc(path, AcceptOrRejectInvitation)
		route.Use(middlewares.Authenticate)

		req, e := http.NewRequest(http.MethodGet, fmt.Sprintf("/%v/%v", "0", "accept"), nil)
		if e != nil {
			t.Fatal(e)
		}
		req.Header.Add("Authorization", token)

		rr := httptest.NewRecorder()
		route.ServeHTTP(rr, req)

		res := rr.Result()
		bts, e := ioutil.ReadAll(res.Body)
		if e != nil {
			t.Fatal(e)
		}
		defer res.Body.Close()

		if res.StatusCode != http.StatusNotFound {
			t.Fatal("Status code not matched. Response is", res.StatusCode)
		}

		if !bytes.Contains(bts, []byte("No match found")) {
			t.Fatal("Invalid response message", string(bts))
		}
	})

	t.Run("Permission denied", func(t *testing.T) {
		invID := createPendingInvitation(id1, id1, id2)
		route := mux.NewRouter()
		route.HandleFunc(path, AcceptOrRejectInvitation)
		route.Use(middlewares.Authenticate)

		req, e := http.NewRequest(http.MethodGet, fmt.Sprintf("/%v/%v", invID, "accept"), nil)
		if e != nil {
			t.Fatal(e)
		}
		req.Header.Add("Authorization", token)

		rr := httptest.NewRecorder()
		route.ServeHTTP(rr, req)

		res := rr.Result()
		bts, e := ioutil.ReadAll(res.Body)
		if e != nil {
			t.Fatal(e)
		}
		defer res.Body.Close()

		if res.StatusCode != http.StatusForbidden {
			t.Fatal("Status code not matched. Response is", res.StatusCode)
		}

		if !bytes.Contains(bts, []byte("Operation not permitted")) {
			t.Fatal("Invalid response message", string(bts))
		}
	})
}
