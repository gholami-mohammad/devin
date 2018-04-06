package migrations

import "devin/database"

// Migrate the database to a new version
func (Migration) MigrateObjectPermissionsTable() (e error) {
	db := database.NewGORMInstance()
	defer db.Close()
	e = db.Exec(`CREATE TABLE IF NOT EXISTS public.object_permissions (
    id bigserial NOT NULL ,
    user_id bigint NOT NULL,
    object_id bigint NOT NULL,
    module_id bigint NOT NULL,
    can_read bool DEFAULT false,
    can_update bool DEFAULT false,
    can_delete bool DEFAULT false,
    created_by_id bigint NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT object_permissions_pkey PRIMARY KEY (id),
    CONSTRAINT object_permissions_composit_unique UNIQUE (user_id, object_id, module_id),
    CONSTRAINT object_permissions_user_id_users_id FOREIGN KEY (user_id)
        REFERENCES public.users (id)
        ON DELETE CASCADE
        ON UPDATE CASCADE,
    CONSTRAINT object_permissions_created_by_id_users_id FOREIGN KEY (created_by_id)
        REFERENCES public.users (id)
        ON DELETE CASCADE
        ON UPDATE CASCADE
    )`).Error

	return
}

// Rollback the database to previous version
func (Migration) RollbackObjectPermissionsTable() (e error) {
	db := database.NewGORMInstance()
	defer db.Close()
	e = db.Exec("DROP TABLE IF EXISTS public.object_permissions;").Error

	return
}
