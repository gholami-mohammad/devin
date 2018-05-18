package repository

import (
	"fmt"
	"strings"
	"testing"

	"devin/database"
	"devin/models"
)

var testID uint64 = 11000

func getTestID() uint64 {
	testID += 1
	return testID
}

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

func TestAddUserToOrganziation(t *testing.T) {
	db := database.NewGORMInstance()
	defer db.Close()
	id1 := getTestID()
	id2 := getTestID()
	id3 := getTestID()
	getValidUser(id1, false)
	defer deleteTestUser(id1)
	getValidOrganization(id2, id1)
	defer deleteTestOrganization(id2)
	getValidUser(id3, false)
	defer deleteTestUser(id3)
	t.Run("OK", func(t *testing.T) {
		e := AddUserToOrganziation(db, id3, id3, id2)
		if e != nil {
			t.Fatal(e)
		}
	})

	t.Run("User exist", func(t *testing.T) {
		AddUserToOrganziation(db, id3, id3, id2)      // add user
		e := AddUserToOrganziation(db, id3, id3, id2) // try to add again
		if e == nil {
			t.Fatal("Duplicated!")
		}

		if !strings.Contains(e.Error(), "member") {
			t.Fatal("Error message not match", e.Error())
		}
	})
}
