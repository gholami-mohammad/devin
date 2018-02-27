package migrations

import "gogit/database"

// Migrate the database to a new version
func (Migration) MigrateTaskBoardsTable() (e error) {
	db := database.NewPGInstance()
	defer db.Close()
	_, e = db.Exec(`CREATE TABLE IF NOT EXISTS public.task_boards(
    id bigserial NOT NULL,
    name varchar(255) NOT NULL,
    project_id bigint,
    color varchar(25),
    created_by_id bigint NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    deleted_at timestamp with time zone,

    CONSTRAINT task_boards_pkey PRIMARY KEY (id),
    CONSTRAINT task_boards_project_id_projects_id FOREIGN KEY (project_id)
        REFERENCES public.projects (id)
        ON DELETE CASCADE
        ON UPDATE CASCADE,
    CONSTRAINT task_boards_created_by_id_users_id FOREIGN KEY (created_by_id)
        REFERENCES public.users (id)
        ON DELETE CASCADE
        ON UPDATE CASCADE
    )`)

	return
}

// Rollback the database to previous version
func (Migration) RollbackTaskBoardsTable() (e error) {
	db := database.NewPGInstance()
	defer db.Close()
	_, e = db.Exec("DROP TABLE IF EISTS public.task_boards CASCADE;")

	return
}
