package repository

import (
	"devin/models"
	"errors"

	"github.com/jinzhu/gorm"
)

type OrganizationSearchable struct {
	UserID *uint64
}

func LoadOrganizationsFilter(db *gorm.DB, searchable OrganizationSearchable) (orgs []models.User, e error) {
	db = db.
		Preload("Owner").
		Preload("OrganizationUserMapping").
		Preload("LocalizationLanguage").
		Preload("CalendarSystem").
		Preload("OfficePhoneCountryCode").
		Preload("HomePhoneCountryCode").
		Preload("CellPhoneCountryCode").
		Preload("FaxCountryCode").
		Preload("Country").
		Preload("Province").
		Preload("City").
		Model(&orgs).
		Where(`user_type=2`)
	if searchable.UserID != nil && *searchable.UserID != 0 {
		db = db.Where("owner_id=? OR id IN (?)", *searchable.UserID, db.Table(models.UserOrganization{}.TableName()).
			Select("organization_id").
			Where("user_id=?", *searchable.UserID).
			QueryExpr())
	}

	e = db.Find(&orgs).Error

	return
}

func AddUserToOrganziation(db *gorm.DB, userID, creatorID, orgID uint64) (e error) {
	var cnt uint
	db.Model(&models.UserOrganization{}).
		Where("user_id=? AND organization_id=?", userID, orgID).
		Count(&cnt)

	if cnt > 0 {
		return errors.New("The user is a member of the organization.")
	}

	obj := models.UserOrganization{}
	obj.CanAddUserToOrganization = false
	obj.CanCreateProject = false
	obj.IsAdminOfOrganization = false
	obj.CreatedByID = creatorID
	obj.OrganizationID = &orgID
	obj.UserID = &userID

	return db.Model(&models.UserOrganization{}).Create(&obj).Error
}
