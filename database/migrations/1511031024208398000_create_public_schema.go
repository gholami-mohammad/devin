package migrations

import "gogit/database"

// CreatePublicSchema Migration Struct
type CreatePublicSchema struct{}

// Migrate the database to a new version
func (CreatePublicSchema) Migrate() {
	db := database.NewPGInstance()
	defer db.Close()
	db.Exec("CREATE SCHEMA public;")

}

// Rollback the database to previous version
func (CreatePublicSchema) Rollback() {
	db := database.NewPGInstance()
	defer db.Close()
	db.Exec("DROP SCHEMA public CASCADE;")
}
