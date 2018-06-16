package models

import (
	"github.com/jinzhu/gorm"

	"devin/helpers"
)

// ProjectSearch is model to performe search on Project model
type ProjectSearch struct {
	ID    *uint64
	Name  *string
	Title *string

	// search on projects of an organization
	OrganizationID *uint64

	StatusIDs []uint

	// If UserID != nil, then search for projects with OwnerUserID and all project members
	// to find projects which this user is its owner or its member
	UserID *uint64

	// =-=-=-=-=-=-=-=-=-=
	// Pagination options
	// =-=-=-=-=-=-=-=-=-=

	// Offset of found records to be return
	Offset uint64

	// Item count to be return
	// If you want to get all matching items, set limit to 0
	Limit uint64

	// Custom order sql
	Order string
}

// GetWhereClause generate where clause using given filters
func (search *ProjectSearch) GetWhereClause(db *gorm.DB) *gorm.DB {
	if helpers.IsNilUint64(search.ID) == false {
		db.Where("id = ? ", search.ID)
	}
	if helpers.IsNilOrEmptyString(search.Name) == false {
		db = db.Where("name LIKE ", "%"+*search.Name+"%")
	}
	if helpers.IsNilOrEmptyString(search.Title) == false {
		db = db.Where("title LIKE ", "%"+*search.Title+"%")
	}

	if helpers.IsNilUint64(search.OrganizationID) == false {
		db = db.Where("owner_organization_id = ?", search.OrganizationID)
	}

	if len(search.StatusIDs) > 0 {
		var statusesInterface []interface{}
		for _, v := range search.StatusIDs {
			statusesInterface = append(statusesInterface, v)
		}
		db = db.Where("staus_id IN ?", statusesInterface...)
	}

	return db
}
