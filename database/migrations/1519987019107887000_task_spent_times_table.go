package migrations

import "gogit/database"

// Migrate the database to a new version
func (Migration) MigrateTaskSpentTimesTable() (e error) {
	db := database.NewPGInstance()
	defer db.Close()
	_, e = db.Exec(`CREATE TABLE IF NOT EXISTS public.task_spent_times (
    id bigserial NOT NULL,
    spent_by_id bigint NOT NULL,
    task_id bigint,
    start_date timestamp with time zone,
    end_date timestamp with time zone,
    is_billabel bool DEFAULT false,
    is_billed bool DEFAULT false,
    description text,
    created_by_id bigint NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    deleted_at timestamp with time zone,

    CONSTRAINT task_spent_times_pkey PRIMARY KEY (id),
    CONSTRAINT task_spent_times_spent_by_id_tasks_id FOREIGN KEY (spent_by_id)
        REFERENCES public.users (id) MATCH SIMPLE
        ON DELETE CASCADE
        ON UPDATE CASCADE,
    CONSTRAINT task_spent_times_task_id_tasks_id FOREIGN KEY (task_id)
        REFERENCES public.tasks (id) MATCH SIMPLE
        ON DELETE CASCADE
        ON UPDATE CASCADE,
    CONSTRAINT task_spent_times_created_by_id_tasks_id FOREIGN KEY (created_by_id)
        REFERENCES public.users (id) MATCH SIMPLE
        ON DELETE CASCADE
        ON UPDATE CASCADE
    )`)

	return
}

// Rollback the database to previous version
func (Migration) RollbackTaskSpentTimesTable() (e error) {
	db := database.NewPGInstance()
	defer db.Close()
	_, e = db.Exec("DROP TABLE IF EXISTS public.task_spent_times;")

	return
}
