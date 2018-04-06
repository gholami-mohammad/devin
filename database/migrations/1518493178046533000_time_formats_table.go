package migrations

import "devin/database"

// Migrate the database to a new version
func (Migration) MigrateTimeFormatsTable() (e error) {
	db := database.NewGORMInstance()
	defer db.Close()
	e = db.Exec(`CREATE TABLE IF NOT EXISTS public.time_formats (
    id serial NOT NULL,
    name varchar (100),
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    deleted_at timestamp with time zone,

    CONSTRAINT time_formats_pkey PRIMARY KEY(id)
    )`).Error

	return
}

// Rollback the database to previous version
func (Migration) RollbackTimeFormatsTable() (e error) {
	db := database.NewGORMInstance()
	defer db.Close()
	e = db.Exec("DROP TABLE IF EXISTS public.time_formats;").Error

	return
}
