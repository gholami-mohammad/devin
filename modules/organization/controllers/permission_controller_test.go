package controllers

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"

	"devin/middlewares"
)

func TestUpdateUserPermissionsOnOrganization(t *testing.T) {

	t.Run("OK", func(t *testing.T) {
		id1 := getTestID()
		id2 := getTestID()

		_, _, token := getValidUser(id1, false)
		defer deleteTestUser(id1)
		getValidOrganization(id2, id1)
		defer deleteTestOrganization(id2)

		route := mux.NewRouter()
		route.Use(middlewares.Authenticate)
		route.HandleFunc("/{organization_id}/{user_id}", UpdateUserPermissionsOnOrganization)

		rr := httptest.NewRecorder()
		req, e := http.NewRequest(http.MethodPost, fmt.Sprintf("/%v/%v", id2, id1), strings.NewReader(`{"IsAdminOfOrganization": true, "CanCreateProject": true}`))
		if e != nil {
			t.Fatal(e)
		}
		req.Header.Add("Authorization", token)
		req.Header.Add("Content-Type", "application/json")

		route.ServeHTTP(rr, req)

		res := rr.Result()
		if res.StatusCode != http.StatusOK {
			t.Fatal("Status code not matched. Response is", res.StatusCode)
		}

	})

	t.Run("Bad Content Type", func(t *testing.T) {
		id1 := getTestID()
		id2 := getTestID()

		_, _, token := getValidUser(id1, false)
		defer deleteTestUser(id1)
		getValidOrganization(id2, id1)
		defer deleteTestOrganization(id2)

		route := mux.NewRouter()
		route.Use(middlewares.Authenticate)
		route.HandleFunc("/{organization_id}/{user_id}", UpdateUserPermissionsOnOrganization)

		rr := httptest.NewRecorder()
		req, e := http.NewRequest(http.MethodPost, fmt.Sprintf("/%v/%v", id2, id1), strings.NewReader(`{"IsAdminOfOrganization": true, "CanCreateProject": true}`))
		if e != nil {
			t.Fatal(e)
		}
		req.Header.Add("Authorization", token)

		route.ServeHTTP(rr, req)

		res := rr.Result()
		if res.StatusCode != http.StatusUnsupportedMediaType {
			t.Fatal("Status code not matched. Response is", res.StatusCode)
		}

	})

	t.Run("Authentication faild", func(t *testing.T) {
		id1 := getTestID()
		id2 := getTestID()

		getValidUser(id1, false)
		defer deleteTestUser(id1)
		getValidOrganization(id2, id1)
		defer deleteTestOrganization(id2)

		route := mux.NewRouter()
		route.HandleFunc("/{organization_id}/{user_id}", UpdateUserPermissionsOnOrganization)

		rr := httptest.NewRecorder()
		req, e := http.NewRequest(http.MethodPost, fmt.Sprintf("/%v/%v", id2, id1), strings.NewReader(`{"IsAdminOfOrganization": true, "CanCreateProject": true}`))
		if e != nil {
			t.Fatal(e)
		}
		req.Header.Add("Content-Type", "application/json")

		route.ServeHTTP(rr, req)

		res := rr.Result()
		if res.StatusCode != http.StatusUnauthorized {
			t.Fatal("Status code not matched. Response is", res.StatusCode)
		}

	})

	t.Run("Bad Organization ID", func(t *testing.T) {
		id1 := getTestID()
		id2 := getTestID()

		_, _, token := getValidUser(id1, false)
		defer deleteTestUser(id1)
		getValidOrganization(id2, id1)
		defer deleteTestOrganization(id2)

		route := mux.NewRouter()
		route.Use(middlewares.Authenticate)
		route.HandleFunc("/{organization_id}/{user_id}", UpdateUserPermissionsOnOrganization)

		rr := httptest.NewRecorder()
		req, e := http.NewRequest(http.MethodPost, fmt.Sprintf("/%v/%v", "BAD_ID", id1), strings.NewReader(`{"IsAdminOfOrganization": true, "CanCreateProject": true}`))
		if e != nil {
			t.Fatal(e)
		}
		req.Header.Add("Authorization", token)
		req.Header.Add("Content-Type", "application/json")

		route.ServeHTTP(rr, req)

		res := rr.Result()
		if res.StatusCode != http.StatusUnprocessableEntity {
			t.Fatal("Status code not matched. Response is", res.StatusCode)
		}

		bts, _ := ioutil.ReadAll(res.Body)
		if !strings.Contains(string(bts), "Invalid Organization ID") {
			t.Fatal("Invalid response message", string(bts))
		}

	})

	t.Run("Bad User ID", func(t *testing.T) {
		id1 := getTestID()
		id2 := getTestID()

		_, _, token := getValidUser(id1, false)
		defer deleteTestUser(id1)
		getValidOrganization(id2, id1)
		defer deleteTestOrganization(id2)

		route := mux.NewRouter()
		route.Use(middlewares.Authenticate)
		route.HandleFunc("/{organization_id}/{user_id}", UpdateUserPermissionsOnOrganization)

		rr := httptest.NewRecorder()
		req, e := http.NewRequest(http.MethodPost, fmt.Sprintf("/%v/%v", id2, "BAD_ID"), strings.NewReader(`{"IsAdminOfOrganization": true, "CanCreateProject": true}`))
		if e != nil {
			t.Fatal(e)
		}
		req.Header.Add("Authorization", token)
		req.Header.Add("Content-Type", "application/json")

		route.ServeHTTP(rr, req)

		res := rr.Result()
		if res.StatusCode != http.StatusUnprocessableEntity {
			t.Fatal("Status code not matched. Response is", res.StatusCode)
		}

		bts, _ := ioutil.ReadAll(res.Body)
		if !strings.Contains(string(bts), "Invalid User ID") {
			t.Fatal("Invalid response message", string(bts))
		}

	})

	t.Run("Organization Not found", func(t *testing.T) {
		id1 := getTestID()
		id2 := getTestID()

		_, _, token := getValidUser(id1, false)
		defer deleteTestUser(id1)

		route := mux.NewRouter()
		route.Use(middlewares.Authenticate)
		route.HandleFunc("/{organization_id}/{user_id}", UpdateUserPermissionsOnOrganization)

		rr := httptest.NewRecorder()
		req, e := http.NewRequest(http.MethodPost, fmt.Sprintf("/%v/%v", id2, id1), strings.NewReader(`{"IsAdminOfOrganization": true, "CanCreateProject": true}`))
		if e != nil {
			t.Fatal(e)
		}
		req.Header.Add("Authorization", token)
		req.Header.Add("Content-Type", "application/json")

		route.ServeHTTP(rr, req)

		res := rr.Result()
		if res.StatusCode != http.StatusNotFound {
			t.Fatal("Status code not matched. Response is", res.StatusCode)
		}

	})

	t.Run("Permission Denied", func(t *testing.T) {
		id1 := getTestID()
		id2 := getTestID()
		id3 := getTestID()

		getValidUser(id1, false)
		defer deleteTestUser(id1)
		_, _, token := getValidUser(id3, false)
		defer deleteTestUser(id3)
		getValidOrganization(id2, id1)
		defer deleteTestOrganization(id2)

		route := mux.NewRouter()
		route.Use(middlewares.Authenticate)
		route.HandleFunc("/{organization_id}/{user_id}", UpdateUserPermissionsOnOrganization)

		rr := httptest.NewRecorder()
		req, e := http.NewRequest(http.MethodPost, fmt.Sprintf("/%v/%v", id2, id3), strings.NewReader(`{"IsAdminOfOrganization": true, "CanCreateProject": true}`))
		if e != nil {
			t.Fatal(e)
		}
		req.Header.Add("Authorization", token)
		req.Header.Add("Content-Type", "application/json")

		route.ServeHTTP(rr, req)

		res := rr.Result()
		if res.StatusCode != http.StatusForbidden {
			t.Fatal("Status code not matched. Response is", res.StatusCode)
		}

	})

	t.Run("Nil request body", func(t *testing.T) {
		id1 := getTestID()
		id2 := getTestID()

		_, _, token := getValidUser(id1, false)
		defer deleteTestUser(id1)
		getValidOrganization(id2, id1)
		defer deleteTestOrganization(id2)

		route := mux.NewRouter()
		route.Use(middlewares.Authenticate)
		route.HandleFunc("/{organization_id}/{user_id}", UpdateUserPermissionsOnOrganization)

		rr := httptest.NewRecorder()
		req, e := http.NewRequest(http.MethodPost, fmt.Sprintf("/%v/%v", id2, id1), nil)
		if e != nil {
			t.Fatal(e)
		}
		req.Header.Add("Authorization", token)
		req.Header.Add("Content-Type", "application/json")

		route.ServeHTTP(rr, req)

		res := rr.Result()
		if res.StatusCode != http.StatusInternalServerError {
			t.Fatal("Status code not matched. Response is", res.StatusCode)
		}

		bts, _ := ioutil.ReadAll(res.Body)
		if !strings.Contains(string(bts), "empty") {
			t.Fatal("Invalid response message", string(bts))
		}

	})

	t.Run("Bad Request Body", func(t *testing.T) {
		id1 := getTestID()
		id2 := getTestID()

		_, _, token := getValidUser(id1, false)
		defer deleteTestUser(id1)
		getValidOrganization(id2, id1)
		defer deleteTestOrganization(id2)

		route := mux.NewRouter()
		route.Use(middlewares.Authenticate)
		route.HandleFunc("/{organization_id}/{user_id}", UpdateUserPermissionsOnOrganization)

		rr := httptest.NewRecorder()
		req, e := http.NewRequest(http.MethodPost, fmt.Sprintf("/%v/%v", id2, id1), strings.NewReader(``))
		if e != nil {
			t.Fatal(e)
		}
		req.Header.Add("Authorization", token)
		req.Header.Add("Content-Type", "application/json")

		route.ServeHTTP(rr, req)

		res := rr.Result()
		if res.StatusCode != http.StatusUnprocessableEntity {
			t.Fatal("Status code not matched. Response is", res.StatusCode)
		}

		bts, _ := ioutil.ReadAll(res.Body)
		if !strings.Contains(string(bts), "Invalid request") {
			t.Fatal("Invalid response message", string(bts))
		}

	})

	t.Run("No user", func(t *testing.T) {
		id1 := getTestID()
		id2 := getTestID()
		id3 := getTestID()

		_, _, token := getValidUser(id1, false)
		defer deleteTestUser(id1)

		getValidOrganization(id2, id1)
		defer deleteTestOrganization(id2)

		route := mux.NewRouter()
		route.Use(middlewares.Authenticate)
		route.HandleFunc("/{organization_id}/{user_id}", UpdateUserPermissionsOnOrganization)

		rr := httptest.NewRecorder()
		req, e := http.NewRequest(http.MethodPost, fmt.Sprintf("/%v/%v", id2, id3), strings.NewReader(`{"IsAdminOfOrganization": true, "CanCreateProject": true}`))
		if e != nil {
			t.Fatal(e)
		}
		req.Header.Add("Authorization", token)
		req.Header.Add("Content-Type", "application/json")

		route.ServeHTTP(rr, req)

		res := rr.Result()
		if res.StatusCode != http.StatusNotFound {
			t.Fatal("Status code not matched. Response is", res.StatusCode)
		}

		bts, _ := ioutil.ReadAll(res.Body)
		if !strings.Contains(string(bts), "Not Found") {
			t.Fatal("Invalid response message", string(bts))
		}

	})
}
