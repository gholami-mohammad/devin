package migrations

import "gogit/database"

// Migrate the database to a new version
func (Migration) MigrateTaskListUsersTable() (e error) {
	db := database.NewPGInstance()
	defer db.Close()
	_, e = db.Exec(`CREATE TABLE IF NOT EXISTS public.task_list_users(
    id bigint NOT NULL,
    user_id bigint NOT NULL,
    task_list_id bigint NOT NULL,
    created_by_id bigint NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT task_list_users_pkey PRIMARY KEY(id),
    CONSTRAINT task_list_users_user_id_users_id FOREIGN KEY(user_id)
        REFERENCES users (id)
        ON DELETE CASCADE
        ON UPDATE CASCADE,
    CONSTRAINT task_list_users_task_list_id_task_lists_id FOREIGN KEY (task_list_id)
        REFERENCES public.task_lists (id)
        ON DELETE CASCADE
        ON UPDATE CASCADE,
    CONSTRAINT task_list_users_created_by_id_users_id FOREIGN KEY (created_by_id)
        REFERENCES public.users (id)
        ON DELETE CASCADE
        ON UPDATE CASCADE

    )`)

	return
}

// Rollback the database to previous version
func (Migration) RollbackTaskListUsersTable() (e error) {
	db := database.NewPGInstance()
	defer db.Close()
	_, e = db.Exec("DROP TABLE IF EXISTS public.task_list_users;")

	return
}
