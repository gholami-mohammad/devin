package migrations

import "devin/database"

// Migrate the database to a new version
func (Migration) MigrateUserOrganizationInvitations() (e error) {
	db := database.NewGORMInstance()
	defer db.Close()
	e = db.Exec(`CREATE TABLE IF NOT EXISTS public.user_organization_invitations (
        id                      bigserial NOT NULL,
        user_id                 bigint,
        email                   varchar(255),
        organization_id         bigint NOT NULL,
        accepted                bool,
        accepted_rejected_at    timestamp with time zone,
        created_by_id           bigint NOT NULL,
        created_at              timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
        updated_at              timestamp with time zone DEFAULT CURRENT_TIMESTAMP,

        CONSTRAINT user_organization_invitations_pkey PRIMARY KEY (id),
        CONSTRAINT user_organization_invitations_organization_id_users_id FOREIGN KEY (organization_id)
            REFERENCES public.users (id)
            ON UPDATE CASCADE
            ON DELETE CASCADE,
        CONSTRAINT user_organization_invitations_created_by_id_users_id FOREIGN KEY (created_by_id)
            REFERENCES public.users (id)
            ON DELETE CASCADE
            ON UPDATE CASCADE
    );`).Error

	return
}

// Rollback the database to previous version
func (Migration) RollbackUserOrganizationInvitations() (e error) {
	db := database.NewGORMInstance()
	defer db.Close()
	e = db.Exec(`DROP TABLE IF EXISTS public.user_organization_invitations;`).Error

	return
}
