package policies

import (
	"devin/models"
)

//CanEditUser check permission of authenticatedUser for editing of requestedUser
func CanEditUser(authenticatedUser models.User, requestedUser models.User) bool {
	if authenticatedUser.ID == requestedUser.ID || authenticatedUser.IsRootUser {
		return true
	}
	return false
}
