package controllers

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"devin/database"
	"devin/helpers"
	"devin/models"

	"github.com/gorilla/mux"
)

func getUnverifiedUser() (user models.User) {
	db := database.NewGORMInstance()
	defer db.Close()
	db.Exec("delete from users where username='unverified';")
	user = models.User{
		Email:    "unverified@example.com",
		Username: "unverified",
	}
	user.SetNewEmailVerificationToken()

	user.EmailVerified = false
	db.Model(&user).Create(&user)

	return
}

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
		res, e := http.Post(server.URL, "application/json", strings.NewReader(`{"Email": "badEMail", "Username": "mgh"}`))
		if e != nil {
			t.Fatal(e)
		}

		defer res.Body.Close()
		if res.StatusCode != http.StatusUnprocessableEntity {
			t.Fatal("Status code not matched. Response is", res.StatusCode)
		}

		var err helpers.ErrorResponse
		json.NewDecoder(res.Body).Decode(&err)
		t.Log(err.Errors["Email"])
		if len(err.Errors["Email"]) == 0 {
			t.Fatal("No email error found")
		}
	})

	t.Run("Invalid username", func(t *testing.T) {
		res, e := http.Post(server.URL, "application/json", strings.NewReader(`{"Email": "m6devin@gmail.com", "Username": "MMMHJSD8234#$%^"}`))
		if e != nil {
			t.Fatal(e)
		}

		defer res.Body.Close()
		if res.StatusCode != http.StatusUnprocessableEntity {
			t.Fatal("Status code not matched. Response is", res.StatusCode)
		}
		var err helpers.ErrorResponse
		json.NewDecoder(res.Body).Decode(&err)
		t.Log(err.Errors["Username"])
		if len(err.Errors["Username"]) == 0 {
			t.Fatal("No username error found")
		}
	})

	t.Run("Invalid password length", func(t *testing.T) {
		res, e := http.Post(server.URL, "application/json", strings.NewReader(`{"Email": "m6devin@gmail.com", "Username": "mgh", "Password": "123"}`))
		if e != nil {
			t.Fatal(e)
		}

		defer res.Body.Close()
		if res.StatusCode != http.StatusUnprocessableEntity {
			t.Fatal("Status code not matched. Response is", res.StatusCode)
		}
		var err helpers.ErrorResponse
		json.NewDecoder(res.Body).Decode(&err)
		t.Log(err.Errors["Password"])
		if len(err.Errors["Password"]) == 0 {
			t.Fatal("No password error found")
		}
	})

	t.Run("Duplicate email", func(t *testing.T) {
		db := database.NewGORMInstance()
		defer db.Close()
		db.Exec("insert into public.users (username,email) values (?,?)", "duplicate", "duplicate@gmail.com")
		defer db.Exec("delete from public.users where username='duplicate'")

		res, e := http.Post(server.URL, "application/json", strings.NewReader(`{"Email": "duplicate@gmail.com", "Username": "noduplicate", "Password": "123123"}`))
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
		t.Log(err.Errors["Email"])
		if len(err.Errors["Email"]) == 0 {
			t.Fatal("No email error found")
		}
	})

	t.Run("Duplicate username", func(t *testing.T) {
		db := database.NewGORMInstance()
		defer db.Close()
		db.Exec("insert into public.users (username,email) values (?,?)", "duplicate", "duplicate@gmail.com")
		defer db.Exec("delete from public.users where username='duplicate'")
		res, e := http.Post(server.URL, "application/json", strings.NewReader(`{"Email": "noduplicate@gmail.com", "Username": "duplicate", "Password": "123123"}`))
		if e != nil {
			t.Fatal(e)
		}

		defer res.Body.Close()
		if res.StatusCode != http.StatusUnprocessableEntity {
			t.Fatal("Status code not matched. Response is", res.StatusCode)
		}
		var err helpers.ErrorResponse
		json.NewDecoder(res.Body).Decode(&err)
		t.Log(err.Errors["Username"])
		if len(err.Errors["Username"]) == 0 {
			t.Fatal("No username error found")
		}
	})

	t.Run("OK data", func(t *testing.T) {
		db := database.NewGORMInstance()
		defer db.Close()
		defer db.Exec("delete from public.users where username='success'")
		res, e := http.Post(server.URL, "application/json", strings.NewReader(`{"Email": "success@gmail.com", "Username": "success", "Password": "123123"}`))
		if e != nil {
			t.Fatal(e)
		}

		defer res.Body.Close()
		if res.StatusCode != http.StatusOK {
			t.Fatal("Status code not matched. Response is", res.StatusCode)
		}
		var user models.User
		json.NewDecoder(res.Body).Decode(&user)
		if user.ID == 0 {
			t.Fatal("Registration failed")
		}
	})
}

func TestVerifySignup(t *testing.T) {
	user := getUnverifiedUser()
	defer deleteTestUser(user.ID)
	path := "/api/signup/verify"
	route := mux.NewRouter()
	route.HandleFunc(path, VerifySignup)

	t.Run("OK", func(t *testing.T) {
		rr := httptest.NewRecorder()
		req, e := http.NewRequest(http.MethodGet, path+"?token="+*user.EmailVerificationToken, nil)
		if e != nil {
			t.Fatal(e)
		}
		route.ServeHTTP(rr, req)

		res := rr.Result()
		if res.StatusCode != http.StatusOK {
			t.Fatal("Status code not matched. Response is", res.StatusCode)
		}
		bts, _ := ioutil.ReadAll(res.Body)
		defer res.Body.Close()

		if !bytes.Contains(bts, []byte("activated")) {
			t.Fatal("Invalid response message", string(bts))
		}
	})

	t.Run("No token passed", func(t *testing.T) {
		rr := httptest.NewRecorder()
		req, e := http.NewRequest(http.MethodGet, path, nil)
		if e != nil {
			t.Fatal(e)
		}
		route.ServeHTTP(rr, req)

		res := rr.Result()
		if res.StatusCode != http.StatusUnprocessableEntity {
			t.Fatal("Status code not matched. Response is", res.StatusCode)
		}
		bts, _ := ioutil.ReadAll(res.Body)
		defer res.Body.Close()

		if !bytes.Contains(bts, []byte("Invalid token")) {
			t.Fatal("Invalid response message", string(bts))
		}
	})

	t.Run("Token not found", func(t *testing.T) {
		rr := httptest.NewRecorder()
		req, e := http.NewRequest(http.MethodGet, path+"?token=notfound", nil)
		if e != nil {
			t.Fatal(e)
		}
		route.ServeHTTP(rr, req)

		res := rr.Result()
		if res.StatusCode != http.StatusNotFound {
			t.Fatal("Status code not matched. Response is", res.StatusCode)
		}
		bts, _ := ioutil.ReadAll(res.Body)
		defer res.Body.Close()

		if !bytes.Contains(bts, []byte("not found")) {
			t.Fatal("Invalid response message", string(bts))
		}
	})

}
