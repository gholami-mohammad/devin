package migrations

import "gogit/database"

// Migrate the database to a new version
func (Migration) MigrateCreateProjectManagementSchema() (e error) {
	db := database.NewPGInstance()
	defer db.Close()
	_, e = db.Exec("CREATE SCHEMA IF NOT EXISTS pm;")
	return
}

// Rollback the database to previous version
func (Migration) RollbackCreateProjectManagementSchema() (e error) {
	db := database.NewPGInstance()
	defer db.Close()
	_, e = db.Exec("DROP SCHEMA IF EXISTS pm CASCADE;")
	return
}
