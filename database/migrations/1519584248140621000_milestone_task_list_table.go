package migrations

import "devin/database"

// Migrate the database to a new version
func (Migration) MigrateMilestoneTaskListTable() (e error) {
	db := database.NewGORMInstance()
	defer db.Close()
	e = db.Exec(`CREATE TABLE IF NOT EXISTS public.milestone_task_lists(
    id bigserial NOT NULL,
    milestone_id bigint NOT NULL,
    task_list_id bigint NOT NULL,
    created_by_id bigint NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT milestone_task_lists_pkey PRIMARY KEY (id),
    CONSTRAINT milestone_task_lists_milestone_id_milestones_id FOREIGN KEY (milestone_id)
        REFERENCES public.milestones (id)
        ON DELETE CASCADE
        ON UPDATE CASCADE,
    CONSTRAINT milestone_task_lists_task_list_id_task_lists_id FOREIGN KEY (task_list_id)
        REFERENCES public.task_lists (id)
        ON DELETE CASCADE
        ON UPDATE CASCADE,
    CONSTRAINT milestone_task_lists_created_by_id_users_id FOREIGN KEY (created_by_id)
        REFERENCES public.users (id)
        ON DELETE CASCADE
        ON UPDATE CASCADE
    )`).Error

	return
}

// Rollback the database to previous version
func (Migration) RollbackMilestoneTaskListTable() (e error) {
	db := database.NewGORMInstance()
	defer db.Close()
	e = db.Exec("DROP TABLE IF EXISTS public.milestone_task_lists CASCADE;").Error

	return
}
