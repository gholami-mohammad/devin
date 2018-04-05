package migrations

import "devin/database"

// Migrate the database to a new version
func (Migration) MigrateCalendarSystemsTable() (e error) {
	db := database.NewPGInstance()
	defer db.Close()
	_, e = db.Exec(`CREATE TABLE IF NOT EXISTS public.calendar_systems (
    id serial NOT NULL,
    name varchar(100),
    component_name varchar(100),
    filter_name varchar(100),
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    deleted_at timestamp with time zone,

    CONSTRAINT calendar_systems_pkey PRIMARY KEY (id)
    )`)

	return
}

// Rollback the database to previous version
func (Migration) RollbackCalendarSystemsTable() (e error) {
	db := database.NewPGInstance()
	defer db.Close()
	_, e = db.Exec("DROP TABLE IF EXISTS public.calendar_systems;")

	return
}
