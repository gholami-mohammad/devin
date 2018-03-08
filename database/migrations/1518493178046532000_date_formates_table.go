package migrations

import "gogit/database"

// Migrate the database to a new version
func (Migration) MigrateDateFormatesTable() (e error) {
	db := database.NewPGInstance()
	defer db.Close()
	_, e = db.Exec(`CREATE TABLE IF NOT EXISTS public.date_formates (
    id serial NOT NULL,
    name varchar (100),
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    deleted_at timestamp with time zone,

    CONSTRAINT date_formates_pkey PRIMARY KEY(id)
    )`)

	return
}

// Rollback the database to previous version
func (Migration) RollbackDateFormatesTable() (e error) {
	db := database.NewPGInstance()
	defer db.Close()
	_, e = db.Exec("DROP TABLE IF EXISTS public.date_formates;")

	return
}
