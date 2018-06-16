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
func SearchProjects(db *gorm.DB, authenticatedUser models.User, searchModel models.ProjectSearch) (pagination models.Pagination, e error) {
	db = db.Model(&models.Project{})
	db = searchModel.GetWhereClause(db)

	if authenticatedUser.ID == 0 {
		// search as anonymouse
	} else if authenticatedUser.IsRootUser == true {
		// seach all projects
	} else if searchModel.UserID == nil || *searchModel.UserID == 0 {
		// seach all public projects + searcher user's projects
	} else if authenticatedUser.ID == *searchModel.UserID {
		// search on authenticatedUser projects
	} else {
		// authenticatedUser is searching on searchModel.UserID's projects
		// just public projects of this user can be shown
	}

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		db.Count(&pagination.Total)
	}()

	go func() {
		defer wg.Done()
		if searchModel.Limit > 0 {
			db = db.Limit(searchModel.Limit)
		}

		db = db.Offset(searchModel.Offset).Find(&pagination.Data)
	}()

	wg.Wait()

	return
}
