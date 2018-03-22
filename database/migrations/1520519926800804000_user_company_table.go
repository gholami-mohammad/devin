package migrations

import "devin/database"

// Migrate the database to a new version
func (Migration) MigrateUserCompanyTable() (e error) {
	db := database.NewPGInstance()
	defer db.Close()
	_, e = db.Exec(`CREATE TABLE IF NOT EXISTS public.user_company (
    id bigserial NOT NULL,
    user_id bigint NOT NULL,
    company_id bigint NOT NULL,
    is_admin_of_company bool DEFAULT false,
    can_create_project bool DEFAULT false,
    can_add_user_to_company bool DEFAULT false,
    created_by_id bigint NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    deleted_at timestamp with time zone,

    CONSTRAINT user_company_pkey PRIMARY KEY (id),
    CONSTRAINT user_company_user_id_users_id FOREIGN KEY (user_id)
        REFERENCES public.users (id)
        ON DELETE CASCADE
        ON UPDATE CASCADE,
    CONSTRAINT user_company_company_id_users_id FOREIGN KEY (company_id)
        REFERENCES public.users (id)
        ON DELETE CASCADE
        ON UPDATE CASCADE,
    CONSTRAINT user_company_created_by_id_users_id FOREIGN KEY (created_by_id)
        REFERENCES public.users (id)
        ON DELETE CASCADE
        ON UPDATE CASCADE
    )`)

	return
}

// Rollback the database to previous version
func (Migration) RollbackUserCompanyTable() (e error) {
	db := database.NewPGInstance()
	defer db.Close()
	_, e = db.Exec("DROP TABLE IF EXISTS public.user_company;")

	return
}
