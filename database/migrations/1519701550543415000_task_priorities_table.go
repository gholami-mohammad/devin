package migrations

import "gogit/database"

// Migrate the database to a new version
func (Migration) MigrateTaskPrioritiesTable() (e error) {
	db := database.NewPGInstance()
	defer db.Close()
	_, e = db.Exec(`CREATE TABLE IF NOT EXISTS public.task_priorities(
    id serial NOT NULL,
    title varchar(255) NOT NULL,
    color varchar(7),
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    deleted_at timestamp with time zone,

    CONSTRAINT task_priorities_pkey PRIMARY KEY (id)
    );`)

	return
}

// Rollback the database to previous version
func (Migration) RollbackTaskPrioritiesTable() (e error) {
	db := database.NewPGInstance()
	defer db.Close()
	_, e = db.Exec("DROP TABLE IF EXISTS public.task_priorities CASCADE;")

	return
}
