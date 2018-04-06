package migrations

import "devin/database"

// Migrate the database to a new version
func (Migration) MigrateTaggedObjectsTable() (e error) {
	db := database.NewGORMInstance()
	defer db.Close()
	e = db.Exec(`CREATE TABLE IF NOT EXISTS public.tagged_objects(
    id bigserial NOT NULL,
    tag_id bigint NOT NULL,
    object_id bigint NOT NULL,
    module_id bigint NOT NULL,
    created_by_id bigint NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT tagged_objects_pkey PRIMARY KEY (id),
    CONSTRAINT tagged_objects_composit_unique UNIQUE (tag_id, object_id, module_id),
    CONSTRAINT tagged_objects_tag_id_tags_id FOREIGN KEY (tag_id)
        REFERENCES public.tags (id)
        ON DELETE CASCADE
        ON UPDATE CASCADE,
    CONSTRAINT tagged_objects_created_by_id_users_id FOREIGN KEY (created_by_id)
        REFERENCES public.users (id)
        ON DELETE CASCADE
        ON UPDATE CASCADE
    )`).Error

	return
}

// Rollback the database to previous version
func (Migration) RollbackTaggedObjectsTable() (e error) {
	db := database.NewGORMInstance()
	defer db.Close()
	e = db.Exec("DROP TABLE IF EXISTS public.tagged_objects;").Error

	return
}
