package migrations

import "devin/database"

// Migrate the database to a new version
func (Migration) MigrateRepositoryContributersTable() (e error) {
	db := database.NewGORMInstance()
	defer db.Close()
	e = db.Exec(`CREATE TABLE IF NOT EXISTS public.repository_contributers(
    id bigserial NOT NULL,
    repository_id bigint NOT NULL,
    user_id bigint NOT NULL,
    role_id integer NOT NULL,
    created_by_id bigint NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    deleted_at timestamp with time zone,

    CONSTRAINT repository_contributers_pkey PRIMARY KEY(id),
    CONSTRAINT repository_contributers_repository_id_repositories_id FOREIGN KEY (repository_id)
        REFERENCES public.repositories (id) MATCH SIMPLE
        ON DELETE CASCADE
        ON UPDATE CASCADE,
    CONSTRAINT repository_contributers_user_id_users_id FOREIGN KEY (user_id)
        REFERENCES public.users (id) MATCH SIMPLE
        ON DELETE CASCADE
        ON UPDATE CASCADE,
    CONSTRAINT repository_contributers_role_id_git_roles_id FOREIGN KEY (role_id)
        REFERENCES public.git_roles (id) MATCH SIMPLE
        ON DELETE CASCADE
        ON UPDATE CASCADE,
    CONSTRAINT repository_contributers_created_by_id_users_id FOREIGN KEY (created_by_id)
        REFERENCES public.users (id) MATCH SIMPLE
        ON DELETE CASCADE
        ON UPDATE CASCADE
    )`).Error

	return
}

// Rollback the database to previous version
func (Migration) RollbackRepositoryContributersTable() (e error) {
	db := database.NewGORMInstance()
	defer db.Close()
	e = db.Exec("DROP TABLE IF EXISTS public.repository_contributers;").Error

	return
}
