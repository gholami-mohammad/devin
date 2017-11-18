package migrations

import "log"
import "gogit/database"

// CreatePublicSchema Migration Struct
type CreatePublicSchema struct{}

// Migrate the database to a new version
func (CreatePublicSchema) Migrate() {
	db := database.NewPGInstance()
	defer db.Close()
	_, e := db.Exec("CREATE SCHEMA public;")
	if e != nil {
		log.Println(e)
	}

}

// Rollback the database to previous version
func (CreatePublicSchema) Rollback() {
	db := database.NewPGInstance()
	defer db.Close()
	_, e := db.Exec("DROP SCHEMA public CASCADE;")
	if e != nil {
		log.Println(e)
	}

}
