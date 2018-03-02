package migrations

import "gogit/database"

// Migrate the database to a new version
func (Migration) MigrateTaskPrerequisitesTable() (e error) {
	db := database.NewPGInstance()
	defer db.Close()
	_, e = db.Exec(`CREATE TABLE IF NOT EXISTS public.task_prerequisites (
    id bigserial NOT NULL,
    task_id bigint NOT NULL,
    prerequisite_id bigint NOT NULL,
    created_by_id bigint NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT task_prerequisites_pkey PRIMARY KEY (id),
    CONSTRAINT task_prerequisites_created_by_id_users_id FOREIGN KEY (created_by_id)
        REFERENCES public.users (id) MATCH SIMPLE
        ON DELETE CASCADE
        ON UPDATE CASCADE
    )`)

	return
}

// Rollback the database to previous version
func (Migration) RollbackTaskPrerequisitesTable() (e error) {
	db := database.NewPGInstance()
	defer db.Close()
	_, e = db.Exec("DROP TABLE IF EXISTS public.task_prerequisites;")

	return
}
