package controllers

import (
	"bytes"
	"devin/database"
	"devin/helpers"
	"encoding/json"
	"time"

	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"
)

func getResetPasswordToken(userID uint64, expireDuration time.Duration) string {

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
		token := getResetPasswordToken(id, 1*time.Hour)
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
		token := getResetPasswordToken(id, -1*time.Hour)
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

func TestResetPassword(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		id := getTestID()
		getValidUser(id, false)
		defer deleteTestUser(id)

		token := getResetPasswordToken(id, 1*time.Hour)
		reqModel := passwordResetReqModel{}
		reqModel.Token = token
		reqModel.Password = "new_valid_pass"
		reqModel.PasswordVerify = "new_valid_pass"

		reqBts, _ := json.Marshal(&reqModel)

		route := mux.NewRouter()
		route.HandleFunc("/reset", ResetPassword)

		req, e := http.NewRequest(http.MethodPost, "/reset", bytes.NewReader(reqBts))
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
		if !strings.Contains(string(bts), "updated") {
			t.Fatal("Invalid response message", string(bts))
		}

	})

	t.Run("Nil request body", func(t *testing.T) {

		route := mux.NewRouter()
		route.HandleFunc("/reset", ResetPassword)

		req, e := http.NewRequest(http.MethodPost, "/reset", nil)
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
		if !strings.Contains(string(bts), "empty") {
			t.Fatal("Invalid response message", string(bts))
		}

	})

	t.Run("Bad request body", func(t *testing.T) {

		route := mux.NewRouter()
		route.HandleFunc("/reset", ResetPassword)

		req, e := http.NewRequest(http.MethodPost, "/reset", strings.NewReader("{Invalid}"))
		if e != nil {
			t.Fatal(e)
		}

		rr := httptest.NewRecorder()
		route.ServeHTTP(rr, req)

		res := rr.Result()
		if res.StatusCode != http.StatusBadRequest {
			t.Fatal("Status code not matched. Response is", res.StatusCode)
		}

		bts, _ := ioutil.ReadAll(res.Body)
		if !strings.Contains(string(bts), "Invalid request") {
			t.Fatal("Invalid response message", string(bts))
		}

	})
}

func TestIsValidResetPassword(t *testing.T) {
	type testModel struct {
		passwordResetReqModel
		RequestedValue bool
	}

	testTable := []testModel{}

	item1 := testModel{}
	item1.Password = "     "
	item1.RequestedValue = false
	testTable = append(testTable, item1)

	item2 := testModel{}
	item2.Password = "123"
	item2.RequestedValue = false
	testTable = append(testTable, item2)

	item3 := testModel{}
	item3.Password = "123123"
	item3.PasswordVerify = "456456"
	item3.RequestedValue = false
	testTable = append(testTable, item3)

	item4 := testModel{}
	item4.Password = "123123"
	item4.PasswordVerify = "123123"
	item4.RequestedValue = true
	testTable = append(testTable, item4)

	for _, v := range testTable {
		rr := httptest.NewRecorder()
		result := isValidResetPassword(rr, v.passwordResetReqModel)
		if result != v.RequestedValue {
			t.Fatal(result)
		}
	}
}
