package migrations

import "reflect"
import "gogit/database"

func init() {
	Migrations["CreatePublicSchema"] = reflect.TypeOf(CreatePublicSchema{})
}

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
