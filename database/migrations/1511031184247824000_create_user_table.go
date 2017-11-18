package migrations

import "log"
import "gogit/database"
import "gogit/models"

// CreateUserTable Migration Struct
type CreateUserTable struct{}

// Migrate the database to a new version
func (CreateUserTable) Migrate() {
	db := database.NewPGInstance()
	defer db.Close()
	_, e := db.Model(&models.User{}).Exec(`
    CREATE TABLE IF NOT EXISTS ?TableName (
    id bigserial NOT NULL,
    username varchar(120) NOT NULL,
    email varchar(255) NOT NULL,
    fname varchar(255),
    lname varchar(255),
    CONSTRAINT public_users_pkey PRIMARY KEY (id),
    CONSTRAINT public_users_username_key UNIQUE (username),
    CONSTRAINT public_users_email_key UNIQUE (email)
    )
    WITH (OIDS = FALSE)
    TABLESPACE pg_default;`)
	if e != nil {
		log.Println(e)
	}

}

// Rollback the database to previous version
func (CreateUserTable) Rollback() {
	db := database.NewPGInstance()
	defer db.Close()
	_, e := db.Model(&models.User{}).Exec("DROP TABLE IF EXISTS ?TableName CASCADE;")
	if e != nil {
		log.Println(e)
	}

}
