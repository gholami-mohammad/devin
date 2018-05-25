package migrations

import "devin/database"

// Migrate the database to a new version
func (Migration) MigratePasswordReset() (e error) {
	db := database.NewGORMInstance()
	defer db.Close()
	e = db.Exec(`CREATE TABLE IF NOT EXISTS public.password_resets (
        id bigserial NOT NULL,
        user_id bigint NOT NULL,
        token varchar(512) NOT NULL,
        used_for_reset boolean DEFAULT false,
        created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
        updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
        expires_at timestamp with time zone,

        CONSTRAINT password_reset_pkey PRIMARY KEY (id),
        CONSTRAINT password_resets_user_id_users_id FOREIGN KEY (user_id)
            REFERENCES public.users (id)
            ON DELETE CASCADE
            ON UPDATE CASCADE

    )`).Error

	return
}

// Rollback the database to previous version
func (Migration) RollbackPasswordReset() (e error) {
	db := database.NewGORMInstance()
	defer db.Close()
	e = db.Exec(`DROP TABLE IF EXISTS public.password_resets;`).Error

	return
}
