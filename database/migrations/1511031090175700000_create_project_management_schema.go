package migrations

import "log"
import "gogit/database"

// CreateProjectManagementSchema Migration Struct
type CreateProjectManagementSchema struct{}

// Migrate the database to a new version
func (CreateProjectManagementSchema) Migrate() {
	db := database.NewPGInstance()
	defer db.Close()
	_, e := db.Exec("CREATE SCHEMA project_management;")
	if e != nil {
		log.Println(e)
	}

}

// Rollback the database to previous version
func (CreateProjectManagementSchema) Rollback() {
	db := database.NewPGInstance()
	defer db.Close()
	_, e := db.Exec("DROP SCHEMA project_management CASCADE;")
	if e != nil {
		log.Println(e)
	}

}
