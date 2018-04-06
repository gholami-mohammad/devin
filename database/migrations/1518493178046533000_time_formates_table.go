package migrations

import "devin/database"

// Migrate the database to a new version
func (Migration) MigrateTimeFormatesTable() (e error) {
	db := database.NewGORMInstance()
	defer db.Close()
	e = db.Exec(`CREATE TABLE IF NOT EXISTS public.time_formates (
    id serial NOT NULL,
    name varchar (100),
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    deleted_at timestamp with time zone,

    CONSTRAINT time_formates_pkey PRIMARY KEY(id)
    )`).Error

	return
}

// Rollback the database to previous version
func (Migration) RollbackTimeFormatesTable() (e error) {
	db := database.NewGORMInstance()
	defer db.Close()
	e = db.Exec("DROP TABLE IF EXISTS public.time_formates;").Error

	return
}
