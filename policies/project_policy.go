package policies

import (
	"github.com/jinzhu/gorm"

	"devin/models"
)

// CanSaveProject check permission of user to save project data
func CanSaveProject(db *gorm.DB, authUser models.User, ownerOrganizationID *uint64, project models.Project) bool {
	// The root user can do anyting
	if authUser.IsRootUser == true {
		return true
	}

	if ownerOrganizationID == nil {
		// Authenticated user is trying to create a personnal project
		if authUser.ID == project.OwnerUserID {
			return true
		}

		return false
	}

	var userOrg models.UserOrganization
	db.Model(&models.UserOrganization{}).Where("user_id=? AND organization_id=?", authUser.ID, *ownerOrganizationID).First(&userOrg)

	// This user is not a member of organization
	if userOrg.ID == 0 {
		return false
	}

	// Admin of organization allowed to do anything on the organization
	if userOrg.IsAdminOfOrganization {
		return true
	}

	// User is trying to create a project
	if project.ID > 0 {
		if userOrg.CanCreateProject == true {
			return true
		}

		return false
	}

	// User is editing an exiting project
	if userOrg.CanUpateProject == true {
		return true
	}

	return false

}
