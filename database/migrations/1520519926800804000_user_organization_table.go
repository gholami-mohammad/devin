package migrations

import "devin/database"

// Migrate the database to a new version
func (Migration) MigrateUserOrganizationTable() (e error) {
	db := database.NewGORMInstance()
	defer db.Close()
	e = db.Exec(`CREATE TABLE IF NOT EXISTS public.user_organization (
    id bigserial NOT NULL,
    user_id bigint NOT NULL,
    organization_id bigint NOT NULL,
    is_admin_of_organization bool DEFAULT false,
    can_create_project bool DEFAULT false,
    can_add_user_to_organization bool DEFAULT false,
    created_by_id bigint NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    deleted_at timestamp with time zone,

    CONSTRAINT user_organization_pkey PRIMARY KEY (id),
    CONSTRAINT user_organization_user_id_users_id FOREIGN KEY (user_id)
        REFERENCES public.users (id)
        ON DELETE CASCADE
        ON UPDATE CASCADE,
    CONSTRAINT user_organization_organization_id_users_id FOREIGN KEY (organization_id)
        REFERENCES public.users (id)
        ON DELETE CASCADE
        ON UPDATE CASCADE,
    CONSTRAINT user_organization_created_by_id_users_id FOREIGN KEY (created_by_id)
        REFERENCES public.users (id)
        ON DELETE CASCADE
        ON UPDATE CASCADE
    )`).Error

	return
}

// Rollback the database to previous version
func (Migration) RollbackUserOrganizationTable() (e error) {
	db := database.NewGORMInstance()
	defer db.Close()
	e = db.Exec("DROP TABLE IF EXISTS public.user_organization;").Error

	return
}
