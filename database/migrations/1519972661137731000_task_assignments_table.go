package migrations

import "devin/database"

// Migrate the database to a new version
func (Migration) MigrateTaskAssignmentsTable() (e error) {
	db := database.NewGORMInstance()
	defer db.Close()
	e = db.Exec(`CREATE TABLE IF NOT EXISTS public.task_assignments (
    id bigserial NOT NULL,
    user_id bigint NOT NULL,
    task_id bigint NOT NULL,
    created_by_id bigint NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT task_assignments_pkey PRIMARY KEY (id),
    CONSTRAINT task_assignments_user_id_users_id FOREIGN KEY (user_id)
        REFERENCES public.users (id) MATCH SIMPLE
        ON DELETE CASCADE
        ON UPDATE CASCADE,
    CONSTRAINT task_assignments_task_id_tasks_id FOREIGN KEY (task_id)
        REFERENCES public.tasks (id) MATCH SIMPLE
        ON DELETE CASCADE
        ON UPDATE CASCADE,
    CONSTRAINT task_assignments_created_by_id_users_id FOREIGN KEY (created_by_id)
        REFERENCES public.users (id) MATCH SIMPLE
        ON DELETE CASCADE
        ON UPDATE CASCADE
    )`).Error

	return
}

// Rollback the database to previous version
func (Migration) RollbackTaskAssignmentsTable() (e error) {
	db := database.NewGORMInstance()
	defer db.Close()
	e = db.Exec("DROP TABLE IF EXISTS public.task_assignments;").Error

	return
}
