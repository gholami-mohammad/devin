package migrations

import "devin/database"

// Migrate the database to a new version
func (Migration) MigrateAddressCountry() (e error) {
	db := database.NewPGInstance()
	defer db.Close()
	_, e = db.Exec(`CREATE TABLE IF NOT EXISTS public.address_countries (
    id serial NOT NULL,
    name varchar(100) NOT NULL,
    phone_Prefix varchar(3),
    alpha2_code varchar(2),
    alpha3_code varchar(3),
    flag text,
    locale_code varchar(6),
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    deleted_at timestamp with time zone,
    CONSTRAINT address_countries_pkey PRIMARY KEY (id)
    );`)
	return
}

// Rollback the database to previous version
func (Migration) RollbackAddressCountry() (e error) {
	db := database.NewPGInstance()
	defer db.Close()
	_, e = db.Exec("DROP TABLE IF EXISTS public.address_countries CASCADE;")

	return
}
