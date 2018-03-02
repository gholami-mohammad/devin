package migrations

import "gogit/database"

// Migrate the database to a new version
func (Migration) MigrateAddressProvincesTable() (e error) {
	db := database.NewPGInstance()
	defer db.Close()
	_, e = db.Exec(`CREATE TABLE IF NOT EXISTS public.address_provinces (
    id serial NOT NULL,
    name varchar(255),
    country_id integer NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    deleted_at timestamp with time zone,

    CONSTRAINT address_provinces_pkey PRIMARY KEY (id),
    CONSTRAINT address_provinces_country_id_address_countries_id FOREIGN KEY (country_id)
        REFERENCES public.address_countries (id) MATCH SIMPLE
        ON DELETE CASCADE
        ON UPDATE CASCADE

    );`)

	return
}

// Rollback the database to previous version
func (Migration) RollbackAddressProvincesTable() (e error) {
	db := database.NewPGInstance()
	defer db.Close()
	_, e = db.Exec("DROP TABLE IF EXISTS public.address_provinces CASCADE;")

	return
}
