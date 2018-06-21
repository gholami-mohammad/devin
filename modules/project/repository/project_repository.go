package repository

import (
	"devin/models"
	"sync"

	"github.com/jinzhu/gorm"
)

// SearchProjects search on projects by given filters
// @param db A new instance of database
// @param authenticatedUser Logged in user
// @param searchModel search filters
func SearchProjects(db *gorm.DB, authenticatedUser models.User, searchModel models.ProjectSearch) (data []models.Project, total uint64, e error) {
	db = db.Debug().Model(&models.Project{})
	db = searchModel.GetWhereClause(db)
	if authenticatedUser.ID == 0 {
		// search as anonymouse
		// TODO
		return
	} else if authenticatedUser.ID == *searchModel.UserID {
		// search on his (authenticatedUser) projects
		db = allMyProjects(db, authenticatedUser.ID)
	} else if authenticatedUser.IsRootUser == true {
		// seach all projects
		// TODO
		return
	} else if searchModel.UserID == nil || *searchModel.UserID == 0 {
		// seach all public projects + searcher user's projects
		// TODO
		return
	} else {
		// authenticatedUser is searching on searchModel.UserID's projects
		// just public projects of this user can be shown
		// TODO
		return
	}

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		db.Count(&total)
	}()

	go func() {
		defer wg.Done()
		if searchModel.PerPage > 0 {
			db = db.Limit(searchModel.PerPage)
		}
		db = db.Offset((searchModel.CurrentPage - 1) * searchModel.PerPage).Find(&data)
	}()

	wg.Wait()

	return
}

// allMyProjects limit search on the given authUserID,
// Logged in user searching on his projects
func allMyProjects(db *gorm.DB, authUserID uint64) *gorm.DB {
	db = db.Where(`owner_user_id=? OR 
		id IN (SELECT project_id FROM project_users WHERE user_id=?)`, authUserID, authUserID)
	return db
}
