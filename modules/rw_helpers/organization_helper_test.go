package rw_helpers

import (
	"fmt"
	"net/http/httptest"
	"testing"

	"devin/database"
	"devin/models"
)

func getValidUser(id uint64, isRoot bool) (user models.User, claim models.Claim, tokenString string) {
	db := database.NewGORMInstance()
	defer db.Close()
	db.Exec(`delete from users where id=?;`, id)
	e := db.Exec(`insert into users (id, username, email, is_root_user, user_type) values (?, ?, ?, ?, 1)`, id, fmt.Sprintf("mgh%v", id), fmt.Sprintf("m6devin%v@gmail.com", id), isRoot).Error
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

func TestCanViewOrganizationsOfUser(t *testing.T) {
	type testItem struct {
		AuthUserID      uint64
		IsRoot          bool
		UserID          uint64
		RequestedResult bool
	}
	var testTable []testItem
	testTable = append(testTable, testItem{
		AuthUserID:      10,
		IsRoot:          false,
		UserID:          10,
		RequestedResult: true,
	})
	testTable = append(testTable, testItem{
		AuthUserID:      11,
		IsRoot:          false,
		UserID:          12,
		RequestedResult: false,
	})
	testTable = append(testTable, testItem{
		AuthUserID:      13,
		IsRoot:          true,
		UserID:          13,
		RequestedResult: true,
	})
	testTable = append(testTable, testItem{
		AuthUserID:      14,
		IsRoot:          true,
		UserID:          15,
		RequestedResult: true,
	})
	for _, x := range testTable {
		authUser, _, _ := getValidUser(x.AuthUserID, x.IsRoot)
		defer deleteTestUser(x.AuthUserID)
		result := CanViewOrganizationsOfUser(httptest.NewRecorder(), authUser, x.UserID)
		if result != x.RequestedResult {
			t.Fatal(x)
		}
	}
}
