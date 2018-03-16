package migrations

import "devin/database"

// Migrate the database to a new version
func (Migration) MigrateMilestonesTable() (e error) {
	db := database.NewPGInstance()
	defer db.Close()
	_, e = db.Exec(`CREATE TABLE IF NOT EXISTS public.milestones (
    id bigserial NOT NULL,
    name varchar(255) NOT NULL,
    due_date timestamp with time zone,
    description text,
    created_by_id bigint,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    deleted_at timestamp with time zone,

    CONSTRAINT milestones_pkey PRIMARY KEY (id),
    CONSTRAINT milestones_created_by_id_users_id FOREIGN KEY (created_by_id)
        REFERENCES public.users (id)
        ON DELETE CASCADE
        ON UPDATE CASCADE
    );`)

	return
}

// Rollback the database to previous version
func (Migration) RollbackMilestonesTable() (e error) {
	db := database.NewPGInstance()
	defer db.Close()
	_, e = db.Exec("DROP TABLE IF EXISTS public.milestones CASCADE;")

	return
}
