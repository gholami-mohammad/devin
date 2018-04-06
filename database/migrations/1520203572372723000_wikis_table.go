package migrations

import "devin/database"

// Migrate the database to a new version
func (Migration) MigrateWikisTable() (e error) {
	db := database.NewGORMInstance()
	defer db.Close()
	e = db.Exec(`CREATE TABLE IF NOT EXISTS public.wikis (
    id bigserial NOT NULL,
    name varchar(255) NOT NULL,
    project_id bigint,
    repository_id bigint,
    created_by_id bigint NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    deleted_at timestamp with time zone,

    CONSTRAINT wikis_pkey PRIMARY KEY (id),
    CONSTRAINT wikis_project_id_projects_id FOREIGN KEY (project_id)
        REFERENCES public.projects (id) MATCH SIMPLE
        ON DELETE CASCADE
        ON UPDATE CASCADE,
    CONSTRAINT wikis_repository_id_repositories_id FOREIGN KEY (repository_id)
        REFERENCES public.repositories (id) MATCH SIMPLE
        ON DELETE CASCADE
        ON UPDATE CASCADE,
    CONSTRAINT wikis_created_by_id_users_id FOREIGN KEY (created_by_id)
        REFERENCES public.users (id) MATCH SIMPLE
        ON DELETE CASCADE
        ON UPDATE CASCADE
    )`).Error

	return
}

// Rollback the database to previous version
func (Migration) RollbackWikisTable() (e error) {
	db := database.NewGORMInstance()
	defer db.Close()
	e = db.Exec("DROP TABLE IF EXISTS public.wikis CASCADE;").Error

	return
}
