package migrations

import "gogit/database"

// Migrate the database to a new version
func (Migration) MigrateTaskRemindersTable() (e error) {
	db := database.NewPGInstance()
	defer db.Close()
	_, e = db.Exec(`CREATE TABLE IF NOT EXISTS public.task_reminders (
    id bigserial NOT NULL,
    task_id bigint NOT NULL,
    title varchar(255),
    reminde_on timestamp with time zone,
    created_by_id bigint NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT task_reminders_pkey PRIMARY KEY (id),
    CONSTRAINT task_reminders_task_id FOREIGN KEY (task_id)
        REFERENCES public.tasks (id) MATCH SIMPLE
        ON DELETE CASCADE
        ON UPDATE CASCADE,
    CONSTRAINT task_reminders_created_by_id_users_id FOREIGN KEY (created_by_id)
        REFERENCES public.users (id) MATCH SIMPLE
        ON DELETE CASCADE
        ON UPDATE CASCADE
    )`)

	return
}

// Rollback the database to previous version
func (Migration) RollbackTaskRemindersTable() (e error) {
	db := database.NewPGInstance()
	defer db.Close()
	_, e = db.Exec("DROP TABLE IF EXISTS public.task_reminders CASCADE;")

	return
}
