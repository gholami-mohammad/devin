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

//CanViewProfile check permission of authenticatedUser for viewing all details of requestedUser
func CanViewProfile(authenticatedUser models.User, requestedUser models.User) bool {
	if authenticatedUser.ID == requestedUser.ID || authenticatedUser.IsRootUser {
		return true
	}
	return false
}

// CanCreateOrganization check permission of authenticatedUser for creating new organization
func CanCreateOrganization(authenticatedUser models.User, requestedOrganization models.User) bool {
	if authenticatedUser.ID == *requestedOrganization.OwnerID || authenticatedUser.IsRootUser {
		return true
	}
	return false
}
