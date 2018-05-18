package repository

import (
	"github.com/jinzhu/gorm"

	"devin/models"
)

// GetPendingInvitaionsByUserID load pending invitations of user
func GetPendingInvitaionsByUserID(db *gorm.DB, userID uint64) (pendingInvitations []models.UserOrganizationInvitation, e error) {
	e = db.Preload("Organization").
		Where("user_id=? AND accepted IS NULL", userID).
		Find(&pendingInvitations).
		Error
	return
}

//GetPendingInvitaionsByID load single invitaion by ID
func GetPendingInvitaionsByID(db *gorm.DB, ID uint64) (pendingInvitation models.UserOrganizationInvitation, e error) {
	e = db.Preload("Organization").
		Where("id=? AND accepted IS NULL", ID).
		Find(&pendingInvitation).
		Error
	return
}

// SetAcceptanceStatusOfInvitaion set true for accepted or false for rejected to `accepted` column
func SetAcceptanceStatusOfInvitaion(db *gorm.DB, ID uint64, status bool) (e error) {
	return db.Model(&models.UserOrganizationInvitation{}).
		Where("id=? AND accepted IS NULL").
		UpdateColumn("accepted", status).
		Error
}
