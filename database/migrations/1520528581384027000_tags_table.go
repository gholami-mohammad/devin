package migrations

import "devin/database"

// Migrate the database to a new version
func (Migration) MigrateTagsTable() (e error) {
	db := database.NewPGInstance()
	defer db.Close()
	_, e = db.Exec(`CREATE TABLE IF NOT EXISTS public.tags (
    id bigserial NOT NULL,
    company_id bigint NOT NULL,
    title varchar(255),
    created_by_id bigint NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    deleted_at timestamp with time zone,

    CONSTRAINT tags_pkey PRIMARY KEY (id),
    CONSTRAINT tags_company_id_users_id FOREIGN KEY (company_id)
        REFERENCES public.users (id)
        ON DELETE CASCADE
        ON UPDATE CASCADE,
    CONSTRAINT tags_created_by_id_users_id FOREIGN KEY (created_by_id)
        REFERENCES public.users (id)
        ON DELETE CASCADE
        ON UPDATE CASCADE
    )`)

	return
}

// Rollback the database to previous version
func (Migration) RollbackTagsTable() (e error) {
	db := database.NewPGInstance()
	defer db.Close()
	_, e = db.Exec("DROP TABLE IF EXISTS public.tags CASCADE;")

	return
}
