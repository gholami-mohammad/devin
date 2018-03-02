package migrations

import "gogit/database"

// Migrate the database to a new version
func (Migration) MigrateTaskReminderReceiversTable() (e error) {
	db := database.NewPGInstance()
	defer db.Close()
	_, e = db.Exec(`CREATE TABLE IF NOT EXISTS public.task_reminder_receivers (
    id bigserial NOT NULL,
    reminder_id bigint NOT NULL,
    user_id bigint NOT NULL,
    notification_types jsonb,
    created_by_id bigint NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT task_reminder_receivers_pkey PRIMARY KEY (id),
    CONSTRAINT task_reminder_receivers_created_by_id_users_id FOREIGN KEY (created_by_id)
        REFERENCES public.users (id) MATCH SIMPLE
        ON DELETE CASCADE
        ON UPDATE CASCADE

    )`)

	return
}

// Rollback the database to previous version
func (Migration) RollbackTaskReminderReceiversTable() (e error) {
	db := database.NewPGInstance()
	defer db.Close()
	_, e = db.Exec("DROP TABLE IF EXISTS public.task_reminder_receivers;")

	return
}
