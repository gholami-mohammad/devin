package migrations

import "devin/database"

// Migrate the database to a new version
func (Migration) MigrateTaskListTable() (e error) {
	db := database.NewPGInstance()
	defer db.Close()
	_, e = db.Exec(`CREATE TABLE IF NOT EXISTS public.task_lists(
    id bigserial NOT NULL,
    name varchar(255) NOT NULL,
    description text,
    milestone_id bigint,
    created_by_id bigint NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    deleted_at timestamp with time zone,

    CONSTRAINT task_lists_pkey PRIMARY KEY (id),
    CONSTRAINT task_lists_milestone_id_milestones_id FOREIGN KEY (milestone_id)
        REFERENCES public.milestones (id)
        ON DELETE CASCADE
        ON UPDATE CASCADE,
    CONSTRAINT task_lists_created_by_id_users_id FOREIGN KEY (created_by_id)
        REFERENCES public.users (id)
        ON DELETE CASCADE
        ON UPDATE CASCADE
    )`)

	return
}

// Rollback the database to previous version
func (Migration) RollbackTaskListTable() (e error) {
	db := database.NewPGInstance()
	defer db.Close()
	_, e = db.Exec("DROP TABLE IF EXISTS public.task_lists CASCADE;")

	return
}
