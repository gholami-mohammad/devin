package migrations

import "devin/database"

// Migrate the database to a new version
func (Migration) MigrateTaskPrioritiesTable() (e error) {
	db := database.NewGORMInstance()
	defer db.Close()
	e = db.Exec(`CREATE TABLE IF NOT EXISTS public.task_priorities(
    id serial NOT NULL,
    title varchar(255) NOT NULL,
    color varchar(7),
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    deleted_at timestamp with time zone,

    CONSTRAINT task_priorities_pkey PRIMARY KEY (id)
    );`).Error

	return
}

// Rollback the database to previous version
func (Migration) RollbackTaskPrioritiesTable() (e error) {
	db := database.NewGORMInstance()
	defer db.Close()
	e = db.Exec("DROP TABLE IF EXISTS public.task_priorities CASCADE;").Error

	return
}
