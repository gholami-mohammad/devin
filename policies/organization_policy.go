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

//CanViewOrganizationsOfUser check permission of authenticatedUser to access organizations list of userID
func CanViewOrganizationsOfUser(authenticatedUser models.User, userID uint64) bool {
	if authenticatedUser.ID == userID || authenticatedUser.IsRootUser == true {
		return true
	}

	return false
}

// CanSeePendingInvitations check permission of authenticatedUser to view pending invitation requests of userID
func CanViewPendingInvitations(authenticatedUser models.User, userID uint64) bool {
	if authenticatedUser.ID == userID || authenticatedUser.IsRootUser == true {
		return true
	}

	return false
}

//CanUserChangeAcceptanceStatus check permission for changing acceptance status of an invitation
func CanUserChangeAcceptanceStatus(user models.User, invitation models.UserOrganizationInvitation) bool {
	// Accepted or rejected invitaions can't be changed
	if invitation.Accepted != nil || invitation.UserID == nil {
		return false
	}

	if user.ID == *invitation.UserID || user.IsRootUser == true {
		return true
	}

	return false
}
