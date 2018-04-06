package migrations

import "devin/database"

// Migrate the database to a new version
func (Migration) MigrateRepositoriesTable() (e error) {
	db := database.NewGORMInstance()
	defer db.Close()
	e = db.Exec(`CREATE TABLE IF NOT EXISTS public.repositories (
    id bigserial NOT NULL,
    owner_id bigint NOT NULL,
    project_id bigint NOT NULL,
    name varchar(255) NOT NULL,
    title varchar(255),
    description text,
    website varchar(300),
    default_branch varchar(300) DEFAULT 'master',
    byte_size bigint,

    watches_count integer,
    stars_count integer,
    issues_count integer,
    forks_count integer,
    closed_issues_number integer,

    is_bare bool,
    is_mirror bool,
    is_private bool,
    enable_pull_request bool,
    is_forked bool,
    forked_from_repository_id bigint,
    created_by_id bigint,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    deleted_at timestamp with time zone,

    CONSTRAINT repositories_pkey PRIMARY KEY (id),
    CONSTRAINT repositories_name_key UNIQUE (name, project_id),
    CONSTRAINT repositories_owner_id_users_id FOREIGN KEY (owner_id)
        REFERENCES public.users (id)
        ON DELETE CASCADE
        ON UPDATE CASCADE,
    CONSTRAINT repositories_project_id_projects_id FOREIGN KEY (project_id)
        REFERENCES public.projects (id)
        ON DELETE CASCADE
        ON UPDATE CASCADE,
    CONSTRAINT repositories_forked_from_repository_id_users_id FOREIGN KEY (forked_from_repository_id)
        REFERENCES public.repositories (id)
        ON DELETE CASCADE
        ON UPDATE CASCADE,
    CONSTRAINT repositories_created_by_id_users_id FOREIGN KEY (created_by_id)
        REFERENCES public.users (id)
        ON DELETE CASCADE
        ON UPDATE CASCADE
    )`).Error

	return
}

// Rollback the database to previous version
func (Migration) RollbackRepositoriesTable() (e error) {
	db := database.NewGORMInstance()
	defer db.Close()
	e = db.Exec("DROP TABLE IF EXISTS public.repositories CASCADE;").Error

	return
}
