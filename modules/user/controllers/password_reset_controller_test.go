package controllers

import (
	"devin/database"
	"devin/helpers"
	"time"

	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"
)

func getResetPasswordTokne(userID uint64, expireDuration time.Duration) string {

	db := database.NewGORMInstance()
	defer db.Close()
	token := helpers.RandomString(64)
	db.Exec(`insert into password_resets
           (user_id, token, expires_at, used_for_reset)
    values (   ?   ,   ?  ,     ?     ,    ?          ) `, userID, token, time.Now().Add(expireDuration), false)

	return token
}

func TestRequestPasswordReset(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		id1 := getTestID()
		user, _, _ := getValidUser(id1, false)
		defer deleteTestUser(id1)
		payload := fmt.Sprintf(`{"Email":"%v"}`, user.Email)
		route := mux.NewRouter()
		route.HandleFunc("/request", RequestPasswordReset)
		req, e := http.NewRequest(http.MethodPost, "/request", strings.NewReader(payload))
		if e != nil {
			t.Fatal(e)
		}
		rr := httptest.NewRecorder()
		route.ServeHTTP(rr, req)

		res := rr.Result()
		if res.StatusCode != http.StatusOK {
			t.Fatal("Status code not matched. Response is", res.StatusCode)
		}

		bts, _ := ioutil.ReadAll(res.Body)
		if !strings.Contains(string(bts), "Password reset link sent to your email") {
			t.Fatal("Invalid response message")
		}
	})

	t.Run("Nil Email", func(t *testing.T) {
		route := mux.NewRouter()
		route.HandleFunc("/request", RequestPasswordReset)
		req, e := http.NewRequest(http.MethodPost, "/request", nil)
		if e != nil {
			t.Fatal(e)
		}
		rr := httptest.NewRecorder()
		route.ServeHTTP(rr, req)

		res := rr.Result()
		if res.StatusCode != http.StatusInternalServerError {
			t.Fatal("Status code not matched. Response is", res.StatusCode)
		}

		bts, _ := ioutil.ReadAll(res.Body)
		if !strings.Contains(string(bts), "Request body cant be empty") {
			t.Fatal("Invalid response message", string(bts))
		}
	})

	t.Run("Inalid email", func(t *testing.T) {
		payload := fmt.Sprintf(`{"Email":"bademail"}`)
		route := mux.NewRouter()
		route.HandleFunc("/request", RequestPasswordReset)
		req, e := http.NewRequest(http.MethodPost, "/request", strings.NewReader(payload))
		if e != nil {
			t.Fatal(e)
		}
		rr := httptest.NewRecorder()
		route.ServeHTTP(rr, req)

		res := rr.Result()
		if res.StatusCode != http.StatusUnprocessableEntity {
			t.Fatal("Status code not matched. Response is", res.StatusCode)
		}

		bts, _ := ioutil.ReadAll(res.Body)
		if !strings.Contains(string(bts), "Invalid email address") {
			t.Fatal("Invalid response message")
		}
	})

	t.Run("No User found by email", func(t *testing.T) {
		payload := fmt.Sprintf(`{"Email":"notfound@nothing.com"}`)
		route := mux.NewRouter()
		route.HandleFunc("/request", RequestPasswordReset)
		req, e := http.NewRequest(http.MethodPost, "/request", strings.NewReader(payload))
		if e != nil {
			t.Fatal(e)
		}
		rr := httptest.NewRecorder()
		route.ServeHTTP(rr, req)

		res := rr.Result()
		if res.StatusCode != http.StatusNotFound {
			t.Fatal("Status code not matched. Response is", res.StatusCode)
		}

		bts, _ := ioutil.ReadAll(res.Body)
		if !strings.Contains(string(bts), "User not found") {
			t.Fatal("Invalid response message", string(bts))
		}
	})

}

func TestValidatePasswordResetLink(t *testing.T) {

	t.Run("OK", func(t *testing.T) {
		id := getTestID()
		getValidUser(id, false)
		defer deleteTestUser(id)
		token := getResetPasswordTokne(id, 1*time.Hour)
		route := mux.NewRouter()
		route.HandleFunc("/validate", ValidatePasswordResetLink)

		req, e := http.NewRequest(http.MethodGet, "/validate?token="+token, nil)
		if e != nil {
			t.Fatal(e)
		}

		rr := httptest.NewRecorder()
		route.ServeHTTP(rr, req)

		res := rr.Result()
		if res.StatusCode != http.StatusOK {
			t.Fatal("Status code not matched. Response is", res.StatusCode)
		}

	})

	t.Run("No Token", func(t *testing.T) {

		route := mux.NewRouter()
		route.HandleFunc("/validate", ValidatePasswordResetLink)

		req, e := http.NewRequest(http.MethodGet, "/validate?token=", nil)
		if e != nil {
			t.Fatal(e)
		}

		rr := httptest.NewRecorder()
		route.ServeHTTP(rr, req)

		res := rr.Result()
		if res.StatusCode != http.StatusUnprocessableEntity {
			t.Fatal("Status code not matched. Response is", res.StatusCode)
		}

		bts, _ := ioutil.ReadAll(res.Body)
		if !strings.Contains(string(bts), "Invalid token") {
			t.Fatal("Invalid response message", string(bts))
		}

	})

	t.Run("Invalid token", func(t *testing.T) {

		route := mux.NewRouter()
		route.HandleFunc("/validate", ValidatePasswordResetLink)

		req, e := http.NewRequest(http.MethodGet, "/validate?token=invalid", nil)
		if e != nil {
			t.Fatal(e)
		}

		rr := httptest.NewRecorder()
		route.ServeHTTP(rr, req)

		res := rr.Result()
		if res.StatusCode != http.StatusUnprocessableEntity {
			t.Fatal("Status code not matched. Response is", res.StatusCode)
		}

		bts, _ := ioutil.ReadAll(res.Body)
		if !strings.Contains(string(bts), "Invalid") {
			t.Fatal("Invalid response message", string(bts))
		}

	})

	t.Run("Expired token", func(t *testing.T) {
		id := getTestID()
		getValidUser(id, false)
		defer deleteTestUser(id)
		token := getResetPasswordTokne(id, -1*time.Hour)
		route := mux.NewRouter()
		route.HandleFunc("/validate", ValidatePasswordResetLink)

		req, e := http.NewRequest(http.MethodGet, "/validate?token="+token, nil)
		if e != nil {
			t.Fatal(e)
		}

		rr := httptest.NewRecorder()
		route.ServeHTTP(rr, req)

		res := rr.Result()
		if res.StatusCode != http.StatusUnprocessableEntity {
			t.Fatal("Status code not matched. Response is", res.StatusCode)
		}

		bts, _ := ioutil.ReadAll(res.Body)
		if !strings.Contains(string(bts), "Expired") {
			t.Fatal("Invalid response message", string(bts))
		}

	})

}
