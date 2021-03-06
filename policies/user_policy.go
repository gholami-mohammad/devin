package policies

import (
	"devin/models"
)

//CanEditUser check permission of authenticatedUser for editing of requestedUser
func CanEditUser(authenticatedUser models.User, requestedUser models.User) bool {
	switch requestedUser.UserType {
	case 1:
		if authenticatedUser.ID == requestedUser.ID || authenticatedUser.IsRootUser {
			return true
		}
		return false
	case 2:
		if authenticatedUser.ID == *requestedUser.OwnerID || authenticatedUser.IsRootUser {
			return true
		}
		return false
	default:
		if authenticatedUser.ID == requestedUser.ID || authenticatedUser.IsRootUser {
			return true
		}
		return false

	}
}

//CanViewProfile check permission of authenticatedUser for viewing all details of requestedUser
func CanViewProfile(authenticatedUser models.User, requestedUser models.User) bool {
	switch requestedUser.UserType {
	case 1:
		if authenticatedUser.ID == requestedUser.ID || authenticatedUser.IsRootUser {
			return true
		}
		return false
	case 2:

		if authenticatedUser.ID == *requestedUser.OwnerID || authenticatedUser.IsRootUser {
			return true
		}

		for _, v := range requestedUser.OrganizationUserMapping {

			if *v.OrganizationID == requestedUser.ID && *v.UserID == authenticatedUser.ID && v.IsAdminOfOrganization == true {
				return true
			}
		}
		return false
	}
	return false
}
