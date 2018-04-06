package middlewares

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"

	"devin/database"
	"devin/models"
)

func TestAuthenticate(t *testing.T) {
	route := mux.NewRouter()

	route.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})
	route.Use(Authenticate)

	server := httptest.NewServer(route)
	defer server.Close()

	t.Run("No_Header_No_Cookie", func(t *testing.T) {
		req, e := http.NewRequest("GET", server.URL, nil)
		if e != nil {
			t.Fatal(e)
		}

		client := http.Client{}
		res, e := client.Do(req)
		if e != nil {
			t.Fatal(e)
		}
		defer res.Body.Close()
		bts, _ := ioutil.ReadAll(res.Body)
		t.Log(string(bts))

		if res.StatusCode != http.StatusUnauthorized {
			t.Fatal("Wrong status code")
		}

		if !strings.Contains(string(bts), "Authentication Token not found") {
			t.Fatal("Wrong response message")
		}
	})

	t.Run("Expired_Token", func(t *testing.T) {
		user := createValidUser()
		db := database.NewGORMInstance()
		defer db.Exec("delete from public.users where username='success_token'")
		defer db.Close()

		claim := user.GenerateNewTokenClaimWithCustomLifetime(-10 * time.Minute)
		tokenString, _ := user.GenerateNewTokenString(claim)
		req, e := http.NewRequest("GET", server.URL, nil)
		if e != nil {
			t.Fatal(e)
		}
		req.Header.Add("Authorization", tokenString)

		client := http.Client{}
		res, e := client.Do(req)
		if e != nil {
			t.Fatal(e)
		}
		defer res.Body.Close()
		bts, _ := ioutil.ReadAll(res.Body)
		t.Log(string(bts))

		if !strings.Contains(string(bts), "Token Expired") {
			t.Fatal("Response dose not match")
		}
	})

	t.Run("Bad_Token", func(t *testing.T) {
		user := createValidUser()
		db := database.NewGORMInstance()
		defer db.Exec("delete from public.users where username='success_token'")
		defer db.Close()

		claim := user.GenerateNewTokenClaim()
		tokenString, _ := user.GenerateNewTokenString(claim)
		tokenString += "23478bjhsdf"
		req, e := http.NewRequest("GET", server.URL, nil)
		if e != nil {
			t.Fatal(e)
		}
		req.Header.Add("Authorization", tokenString)

		client := http.Client{}
		res, e := client.Do(req)
		if e != nil {
			t.Fatal(e)
		}
		defer res.Body.Close()
		bts, _ := ioutil.ReadAll(res.Body)
		t.Log(string(bts))

		if !strings.Contains(string(bts), "Error while Parsing Token") {
			t.Fatal("Response dose not match")
		}
	})

	t.Run("Bad_Payload", func(t *testing.T) {
		user := createValidUser()
		db := database.NewGORMInstance()
		defer db.Exec("delete from public.users where username='success_token'")
		defer db.Close()

		claim := user.GenerateNewTokenClaim()
		claim.Payload += "6d"
		tokenString, _ := user.GenerateNewTokenString(claim)
		req, e := http.NewRequest("GET", server.URL, nil)
		if e != nil {
			t.Fatal(e)
		}
		req.Header.Add("Authorization", tokenString)

		client := http.Client{}
		res, e := client.Do(req)
		if e != nil {
			t.Fatal(e)
		}
		defer res.Body.Close()
		bts, _ := ioutil.ReadAll(res.Body)
		t.Log(string(bts))

		if !strings.Contains(string(bts), "Auhtentication failed (Bad payload)") {
			t.Fatal("Response dose not match")
		}
	})

	t.Run("Not found user", func(t *testing.T) {
		user := models.User{
			ID:       123456,
			Username: "MGH_Notfound",
			Email:    "NOTEXISTS@example.com",
		}

		claim := user.GenerateNewTokenClaim()
		tokenString, _ := user.GenerateNewTokenString(claim)

		req, e := http.NewRequest("GET", server.URL, nil)
		if e != nil {
			t.Fatal(e)
		}
		req.Header.Add("Authorization", tokenString)

		client := http.Client{}
		res, e := client.Do(req)
		if e != nil {
			t.Fatal(e)
		}
		defer res.Body.Close()
		bts, _ := ioutil.ReadAll(res.Body)
		t.Log(string(bts))

		if !strings.Contains(string(bts), "Auhtentication failed (User not found)") {
			t.Fatal("Response dose not match")
		}
	})

	t.Run("OK, less than 25% lifetime", func(t *testing.T) {
		user := createValidUser()
		db := database.NewGORMInstance()
		defer db.Exec("delete from public.users where username='success_token'")
		defer db.Close()

		claim := user.GenerateNewTokenClaimWithCustomLifetime(user.TokenLifetime() / 10)
		tokenString, _ := user.GenerateNewTokenString(claim)
		req, e := http.NewRequest("GET", server.URL, nil)
		if e != nil {
			t.Fatal(e)
		}
		req.Header.Add("Authorization", tokenString)

		client := http.Client{}
		res, e := client.Do(req)
		if e != nil {
			t.Fatal(e)
		}
		defer res.Body.Close()
		bts, _ := ioutil.ReadAll(res.Body)
		t.Log(string(bts))

		if !strings.EqualFold(string(bts), "OK") {
			t.Fatal("Response dose not match")
		}

	})

	t.Run("OK", func(t *testing.T) {
		user := createValidUser()
		db := database.NewGORMInstance()
		defer db.Exec("delete from public.users where username='success_token'")
		defer db.Close()

		claim := user.GenerateNewTokenClaim()
		tokenString, _ := user.GenerateNewTokenString(claim)
		req, e := http.NewRequest("GET", server.URL, nil)
		if e != nil {
			t.Fatal(e)
		}
		req.Header.Add("Authorization", tokenString)

		client := http.Client{}
		res, e := client.Do(req)
		if e != nil {
			t.Fatal(e)
		}
		defer res.Body.Close()
		bts, _ := ioutil.ReadAll(res.Body)
		t.Log(string(bts))

		if !strings.EqualFold(string(bts), "OK") {
			t.Fatal("Response dose not match")
		}

	})
}

func createValidUser() models.User {
	db := database.NewGORMInstance()
	defer db.Close()
	bts, _ := bcrypt.GenerateFromPassword([]byte("pswd"), bcrypt.DefaultCost)
	db.Exec("insert into public.users (username,email,password, email_verified) values (?,?,?,?)", "success_token", "success_token@gmail.com", string(bts), true)

	var user models.User
	db.Where("username='success_token'").First(&user)

	return user
}
