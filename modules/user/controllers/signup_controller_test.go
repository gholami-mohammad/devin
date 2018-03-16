package controllers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"devin/database"
	"devin/helpers"
	"devin/models"
)

func TestSignup(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(Signup))
	defer server.Close()
	t.Run("bad content-type", func(t *testing.T) {
		res, e := http.Post(server.URL, "text/plain", nil)
		if e != nil {
			t.Fatal(e)
		}

		defer res.Body.Close()
		if res.StatusCode != http.StatusUnsupportedMediaType {
			t.Fatal("Status code not matched. Response is", res.StatusCode)
		}
	})

	t.Run("Bad request body", func(t *testing.T) {
		res, e := http.Post(server.URL, "application/json", strings.NewReader("a bad request"))
		if e != nil {
			t.Fatal(e)
		}

		defer res.Body.Close()
		if res.StatusCode != http.StatusBadRequest {
			t.Fatal("Status code not matched. Response is", res.StatusCode)
		}
	})

	t.Run("Invalid email address", func(t *testing.T) {
		res, e := http.Post(server.URL, "application/json", strings.NewReader(`{"email": "badEMail", "username": "mgh"}`))
		if e != nil {
			t.Fatal(e)
		}

		defer res.Body.Close()
		if res.StatusCode != http.StatusUnprocessableEntity {
			t.Fatal("Status code not matched. Response is", res.StatusCode)
		}

		var err helpers.ErrorResponse
		json.NewDecoder(res.Body).Decode(&err)
		t.Log(err.Errors["email"])
		if len(err.Errors["email"]) == 0 {
			t.Fatal("No email error found")
		}
	})

	t.Run("Invalid username", func(t *testing.T) {
		res, e := http.Post(server.URL, "application/json", strings.NewReader(`{"email": "m6devin@gmail.com", "username": "MMMHJSD8234#$%^"}`))
		if e != nil {
			t.Fatal(e)
		}

		defer res.Body.Close()
		if res.StatusCode != http.StatusUnprocessableEntity {
			t.Fatal("Status code not matched. Response is", res.StatusCode)
		}
		var err helpers.ErrorResponse
		json.NewDecoder(res.Body).Decode(&err)
		t.Log(err.Errors["username"])
		if len(err.Errors["username"]) == 0 {
			t.Fatal("No username error found")
		}
	})

	t.Run("Invalid password length", func(t *testing.T) {
		res, e := http.Post(server.URL, "application/json", strings.NewReader(`{"email": "m6devin@gmail.com", "username": "mgh", "password": "123"}`))
		if e != nil {
			t.Fatal(e)
		}

		defer res.Body.Close()
		if res.StatusCode != http.StatusUnprocessableEntity {
			t.Fatal("Status code not matched. Response is", res.StatusCode)
		}
		var err helpers.ErrorResponse
		json.NewDecoder(res.Body).Decode(&err)
		t.Log(err.Errors["password"])
		if len(err.Errors["password"]) == 0 {
			t.Fatal("No password error found")
		}
	})

	t.Run("Duplicate email", func(t *testing.T) {
		db := database.NewPGInstance()
		defer db.Close()
		db.Exec("insert into public.users (username,email) values (?,?)", "duplicate", "duplicate@gmail.com")
		defer db.Exec("delete from public.users where username='duplicate'")

		res, e := http.Post(server.URL, "application/json", strings.NewReader(`{"email": "duplicate@gmail.com", "username": "noduplicate", "password": "123123"}`))
		if e != nil {
			t.Fatal(e)
		}

		defer res.Body.Close()
		if res.StatusCode != http.StatusUnprocessableEntity {
			t.Fatal("Status code not matched. Response is", res.StatusCode)
		}
		var err helpers.ErrorResponse
		json.NewDecoder(res.Body).Decode(&err)
		t.Log(err)
		t.Log(err.Errors["email"])
		if len(err.Errors["email"]) == 0 {
			t.Fatal("No email error found")
		}
	})

	t.Run("Duplicate username", func(t *testing.T) {
		db := database.NewPGInstance()
		defer db.Close()
		db.Exec("insert into public.users (username,email) values (?,?)", "duplicate", "duplicate@gmail.com")
		defer db.Exec("delete from public.users where username='duplicate'")
		res, e := http.Post(server.URL, "application/json", strings.NewReader(`{"email": "noduplicate@gmail.com", "username": "duplicate", "password": "123123"}`))
		if e != nil {
			t.Fatal(e)
		}

		defer res.Body.Close()
		if res.StatusCode != http.StatusUnprocessableEntity {
			t.Fatal("Status code not matched. Response is", res.StatusCode)
		}
		var err helpers.ErrorResponse
		json.NewDecoder(res.Body).Decode(&err)
		t.Log(err.Errors["username"])
		if len(err.Errors["username"]) == 0 {
			t.Fatal("No username error found")
		}
	})

	t.Run("OK data", func(t *testing.T) {
		db := database.NewPGInstance()
		defer db.Close()
		defer db.Exec("delete from public.users where username='success'")
		res, e := http.Post(server.URL, "application/json", strings.NewReader(`{"email": "success@gmail.com", "username": "success", "password": "123123"}`))
		if e != nil {
			t.Fatal(e)
		}

		defer res.Body.Close()
		if res.StatusCode != http.StatusOK {
			t.Fatal("Status code not matched. Response is", res.StatusCode)
		}
		var user models.User
		json.NewDecoder(res.Body).Decode(&user)
		t.Log("Registration succssed. User ID is ", user.ID)
	})
}
