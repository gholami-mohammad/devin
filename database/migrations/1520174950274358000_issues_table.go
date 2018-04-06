package migrations

import "devin/database"

// Migrate the database to a new version
func (Migration) MigrateIssuesTable() (e error) {
	db := database.NewGORMInstance()
	defer db.Close()
	e = db.Exec(`CREATE TABLE IF NOT EXISTS public.issues(
    id bigserial NOT NULL,
    project_id bigint NOT NULL,
    repository_id bigint NOT NULL,
    message text,
    attachment_file_path varchar(255),
    status_id integer,
    set_as_in_progress_datetime timestamp with time zone,
    set_as_resolved_datetime timestamp with time zone,
    created_by_id bigint NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    deleted_at timestamp with time zone,

    CONSTRAINT issues_pkey PRIMARY KEY (id),
    CONSTRAINT issues_project_id_projects_id FOREIGN KEY (project_id)
        REFERENCES public.projects (id)
        ON DELETE CASCADE
        ON UPDATE CASCADE,
    CONSTRAINT issues_repository_id_repositories_id FOREIGN KEY (repository_id)
        REFERENCES public.repositories (id)
        ON DELETE CASCADE
        ON UPDATE CASCADE,
    CONSTRAINT issues_created_by_id_created_by_id FOREIGN KEY (created_by_id)
        REFERENCES public.users (id)
        ON DELETE CASCADE
        ON UPDATE CASCADE
    )`).Error

	return
}

// Rollback the database to previous version
func (Migration) RollbackIssuesTable() (e error) {
	db := database.NewGORMInstance()
	defer db.Close()
	e = db.Exec("DROP TABLE IF EXISTS public.issues CASCADE;").Error

	return
}
