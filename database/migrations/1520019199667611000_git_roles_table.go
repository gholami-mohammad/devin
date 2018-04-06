package migrations

import "devin/database"

// Migrate the database to a new version
func (Migration) MigrateGitRolesTable() (e error) {
	db := database.NewGORMInstance()
	defer db.Close()
	e = db.Exec(`CREATE TABLE IF NOT EXISTS public.git_roles (
    id serial NOT NULL,
    name varchar(100) NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    deleted_at timestamp with time zone,

    CONSTRAINT git_roles_pkey PRIMARY KEY (id)
    );`).Error

	return
}

// Rollback the database to previous version
func (Migration) RollbackGitRolesTable() (e error) {
	db := database.NewGORMInstance()
	defer db.Close()
	e = db.Exec("DROP TABLE IF EXISTS public.git_roles CASCADE;").Error

	return
}
