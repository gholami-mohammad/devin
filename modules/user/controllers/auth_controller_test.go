package controllers

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"golang.org/x/crypto/bcrypt"

	"devin/database"
)

func TestSignin(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(Signin))
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

	t.Run("Empty Email", func(t *testing.T) {
		res, e := http.Post(server.URL, "application/json", strings.NewReader(`{"email":null, "password": "123"}`))
		if e != nil {
			t.Fatal(e)
		}

		defer res.Body.Close()
		if res.StatusCode != http.StatusUnprocessableEntity {
			t.Fatal("Status code not matched. Response is", res.StatusCode)
		}
	})

	t.Run("Empty Password", func(t *testing.T) {
		res, e := http.Post(server.URL, "application/json", strings.NewReader(`{"email": "mgh", "password": ""}`))
		if e != nil {
			t.Fatal(e)
		}

		defer res.Body.Close()
		if res.StatusCode != http.StatusUnprocessableEntity {
			t.Fatal("Status code not matched. Response is", res.StatusCode)
		}
	})

	t.Run("No account found", func(t *testing.T) {
		res, e := http.Post(server.URL, "application/json", strings.NewReader(`{"email": "noone", "password": "pswd"}`))
		if e != nil {
			t.Fatal(e)
		}

		defer res.Body.Close()
		if res.StatusCode != http.StatusUnauthorized {
			t.Fatal("Status code not matched. Response is", res.StatusCode)
		}
	})

	t.Run("Email not verified", func(t *testing.T) {
		db := database.NewPGInstance()
		defer db.Close()
		bts, _ := bcrypt.GenerateFromPassword([]byte("pswd"), bcrypt.DefaultCost)
		db.Exec("insert into public.users (username,email,password,email_verified) values (?,?,?,?)", "success", "success@gmail.com", string(bts), false)
		defer db.Exec("delete from public.users where username='success'")

		res, e := http.Post(server.URL, "application/json", strings.NewReader(`{"email": "success@gmail.com", "password": "pswd"}`))
		if e != nil {
			t.Fatal(e)
		}

		defer res.Body.Close()
		if res.StatusCode != http.StatusUnauthorized {
			t.Fatal("Status code not matched. Response is", res.StatusCode)
		}
	})

	t.Run("wrong password", func(t *testing.T) {
		db := database.NewPGInstance()
		defer db.Close()
		bts, _ := bcrypt.GenerateFromPassword([]byte("pswd"), bcrypt.DefaultCost)
		db.Exec("insert into public.users (username,email,password, email_verified) values (?,?,?,?)", "success", "success@gmail.com", string(bts), true)
		defer db.Exec("delete from public.users where username='success'")

		res, e := http.Post(server.URL, "application/json", strings.NewReader(`{"email": "success@gmail.com", "password": "bad_pass"}`))
		if e != nil {
			t.Fatal(e)
		}

		defer res.Body.Close()
		bts, _ = ioutil.ReadAll(res.Body)
		t.Log(string(bts))
		if res.StatusCode != http.StatusUnauthorized {
			t.Fatal("Status code not matched. Response is", res.StatusCode)
		}
	})

	t.Run("Success login with email", func(t *testing.T) {
		db := database.NewPGInstance()
		defer db.Close()
		bts, _ := bcrypt.GenerateFromPassword([]byte("pswd"), bcrypt.DefaultCost)
		db.Exec("insert into public.users (username,email,password, email_verified) values (?,?,?,?)", "success_email", "success_email@gmail.com", string(bts), true)
		defer db.Exec("delete from public.users where username='success_email'")

		res, e := http.Post(server.URL, "application/json", strings.NewReader(`{"email": "success_email@gmail.com", "password": "pswd"}`))
		if e != nil {
			t.Fatal(e)
		}

		defer res.Body.Close()
		// bts, _ = ioutil.ReadAll(res.Body)
		// t.Log(string(bts))
		if res.StatusCode != http.StatusOK {
			t.Fatal("Status code not matched. Response is", res.StatusCode)
		}
	})

	t.Run("Success login with username", func(t *testing.T) {
		db := database.NewPGInstance()
		defer db.Close()
		bts, _ := bcrypt.GenerateFromPassword([]byte("pswd"), bcrypt.DefaultCost)
		db.Exec("insert into public.users (username,email,password, email_verified) values (?,?,?,?)", "success_username", "success_username@gmail.com", string(bts), true)
		defer db.Exec("delete from public.users where username='success_username'")

		res, e := http.Post(server.URL, "application/json", strings.NewReader(`{"email": "success_username@gmail.com", "password": "pswd"}`))
		if e != nil {
			t.Fatal(e)
		}

		defer res.Body.Close()
		// bts, _ = ioutil.ReadAll(res.Body)
		// t.Log(string(bts))
		if res.StatusCode != http.StatusOK {
			t.Fatal("Status code not matched. Response is", res.StatusCode)
		}
	})
}
