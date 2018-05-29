package controllers

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"
)

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
