package policies

import (
	"devin/database"
	"devin/models"
)

// CanCreateOrganization check permission of authenticatedUser for creating new organization
func CanCreateOrganization(authenticatedUser models.User, requestedOrganization models.User) bool {
	if authenticatedUser.ID == *requestedOrganization.OwnerID || authenticatedUser.IsRootUser {
		return true
	}
	return false
}

// CanInviteUserToOrganization check permission of user
// for inviting new user to given organization
func CanInviteUserToOrganization(authenticatedUser models.User, requestedOrganization models.User) bool {
	if authenticatedUser.ID == *requestedOrganization.OwnerID || authenticatedUser.IsRootUser {
		return true
	}

	db := database.NewGORMInstance()
	defer db.Close()
	var orgUser models.UserOrganization
	db.Model(&orgUser).Where("user_id=? and organization_id=?", authenticatedUser.ID, requestedOrganization.ID).First(&orgUser)
	if orgUser.ID == 0 {
		return false
	}

	if orgUser.IsAdminOfOrganization == true || orgUser.CanAddUserToOrganization == true {
		return true
	}

	return false
}
