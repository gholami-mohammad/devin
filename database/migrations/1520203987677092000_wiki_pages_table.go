package migrations

import "gogit/database"

// Migrate the database to a new version
func (Migration) MigrateWikiPagesTable() (e error) {
	db := database.NewPGInstance()
	defer db.Close()
	_, e = db.Exec(`CREATE TABLE IF NOT EXISTS public.wiki_pages (
    id bigserial NOT NULL,
    wiki_id bigint NOT NULL,
    created_by_id bigint NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    deleted_at timestamp with time zone,

    CONSTRAINT wiki_pages_pkey PRIMARY KEY (id),
    CONSTRAINT wiki_pages_wiki_id_wikis_id FOREIGN KEY (wiki_id)
        REFERENCES public.wikis (id)
        ON DELETE CASCADE
        ON UPDATE CASCADE,
    CONSTRAINT wiki_pages_created_by_id_users_id FOREIGN KEY (created_by_id)
        REFERENCES public.users (id) MATCH SIMPLE
        ON DELETE CASCADE
        ON UPDATE CASCADE
    )`)

	return
}

// Rollback the database to previous version
func (Migration) RollbackWikiPagesTable() (e error) {
	db := database.NewPGInstance()
	defer db.Close()
	_, e = db.Exec("DROP TABLE IF EXISTS public.wiki_pages;")

	return
}
