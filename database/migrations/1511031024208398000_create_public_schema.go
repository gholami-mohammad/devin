package migrations

import "gogit/database"

// Migrate the database to a new version
func (Migration) MigrateCreatePublicSchema() (e error) {
	db := database.NewPGInstance()
	defer db.Close()
	_, e = db.Exec("CREATE SCHEMA IF NOT EXISTS public;")

	return
}

// Rollback the database to previous version
func (Migration) RollbackCreatePublicSchema() (e error) {
	db := database.NewPGInstance()
	defer db.Close()
	_, e = db.Exec("DROP SCHEMA IF EXISTS public CASCADE;")

	return
}
