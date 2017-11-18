package migrations

import "gogit/database"

// CreateProjectManagementSchema Migration Struct
type CreateProjectManagementSchema struct{}

// Migrate the database to a new version
func (CreateProjectManagementSchema) Migrate() {
	db := database.NewPGInstance()
	defer db.Close()
	db.Exec("CREATE SCHEMA project_management;")

}

// Rollback the database to previous version
func (CreateProjectManagementSchema) Rollback() {
	db := database.NewPGInstance()
	defer db.Close()
	db.Exec("DROP SCHEMA project_management CASCADE;")

}
