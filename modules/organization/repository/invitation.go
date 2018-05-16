package repository

import (
	"github.com/jinzhu/gorm"

	"devin/models"
)

func GetPendingInvitaionsByUserID(db *gorm.DB, userID uint64) (pendingInvitations []models.UserOrganizationInvitation, e error) {
	e = db.Preload("Organization").
		Where("user_id=? AND accepted IS NULL", userID).
		Find(&pendingInvitations).
		Error
	return
}
