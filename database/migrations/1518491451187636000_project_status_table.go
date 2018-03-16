package migrations

import "devin/database"

// Migrate the database to a new version
func (Migration) MigrateProjectStatusTable() (e error) {
	db := database.NewPGInstance()
	defer db.Close()
	_, e = db.Exec(`CREATE TABLE IF NOT EXISTS public.project_statuses(
    id serial NOT NULL,
    status varchar(50) NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    deleted_at timestamp with time zone,

    CONSTRAINT project_statuses_pkey PRIMARY KEY(id)
    )`)

	return
}

// Rollback the database to previous version
func (Migration) RollbackProjectStatusTable() (e error) {
	db := database.NewPGInstance()
	defer db.Close()
	_, e = db.Exec("DROP TABLE IF EXISTS public.project_statuses CASCADE;")

	return
}
