package migrations

import "gogit/database"

// Migrate the database to a new version
func (Migration) MigrateAddressCitiesTable() (e error) {
	db := database.NewPGInstance()
	defer db.Close()
	_, e = db.Exec(`CREATE TABLE IF NOT EXISTS public.address_cities(
    id serial NOT NULL,
    name varchar(200) NOT NULL,
    province_id integer NOT NULL,
    country_id integer NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    deleted_at timestamp with time zone,

    CONSTRAINT address_cities_pkey PRIMARY KEY (id),
    CONSTRAINT address_cities_province_id_address_provinces_id FOREIGN KEY (province_id)
        REFERENCES public.address_provinces (id)
        ON DELETE CASCADE
        ON UPDATE CASCADE,
    CONSTRAINT address_cities_country_id_address_countries_id FOREIGN KEY (country_id)
        REFERENCES public.address_countries (id)
        ON DELETE CASCADE
        ON UPDATE CASCADE
    );`)

	return
}

// Rollback the database to previous version
func (Migration) RollbackAddressCitiesTable() (e error) {
	db := database.NewPGInstance()
	defer db.Close()
	_, e = db.Exec("DROP TABLE IF EXISTS public.address_cities;")

	return
}
