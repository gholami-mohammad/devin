package test_helpers

import (
	"fmt"

	"devin/database"
	"devin/models"
)

func GetValidUser(id uint64, isRoot bool) (user models.User, claim models.Claim, tokenString string) {
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

func DeleteTestUser(id uint64) {
	db := database.NewGORMInstance()
	defer db.Close()
	db.Exec(`delete from users where id=?;`, id)
}

func GetValidOrganization(id uint64, ownerID uint64) models.User {
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

func DeleteTestOrganization(id uint64) {
	DeleteTestUser(id)
}

func AddUserToOrganization(userID, orgID uint64) {
	db := database.NewGORMInstance()
	defer db.Close()
	e := db.Exec(`insert into user_organization (user_id, organization_id, created_by_id) values (?, ?, ?)`, userID, orgID, userID).Error
	if e != nil {
		panic(e)
	}
}
