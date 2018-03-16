package migrations

import "devin/database"

// Migrate the database to a new version
func (Migration) MigrateUsersTable() (e error) {
	db := database.NewPGInstance()
	defer db.Close()
	_, e = db.Exec(`CREATE TABLE IF NOT EXISTS public.users (
    id bigserial NOT NULL,
    username varchar(200) NOT NULL,
    email varchar(300) NOT NULL,
    first_name varchar(150),
    last_name varchar(150),
    user_type integer,
    avatar varchar(200),
    owner_id bigint,
    password varchar(512),

    job_title varchar(150),
    localization_language_id integer,
    date_format varchar(100),
    time_format varchar(100),
    calendar_system_id integer,
    office_phone_country_code_id integer,
    home_phone_country_code_id integer,
    cell_phone_country_code_id integer,
    country_id integer,
    province_id integer,
    city_id integer,
    twitter varchar(300),
    linkedin varchar(300),
    google_plus varchar(300),
    facebook varchar(300),
    telegram varchar(300),
    website varchar(300),

    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    deleted_at timestamp with time zone,

    CONSTRAINT users_pkey PRIMARY KEY (id),
    CONSTRAINT users_username_key UNIQUE (username),
    CONSTRAINT users_email_key UNIQUE (email),
    CONSTRAINT users_owner_id_users_id FOREIGN KEY (owner_id)
        REFERENCES public.users (id)
        ON DELETE CASCADE
        ON UPDATE CASCADE,
    CONSTRAINT users_localization_language_id_address_countries_id FOREIGN KEY (localization_language_id)
        REFERENCES public.address_countries (id)
        ON DELETE CASCADE
        ON UPDATE CASCADE,
    CONSTRAINT users_calendar_system_id_calendar_systems_id FOREIGN KEY (calendar_system_id)
        REFERENCES public.calendar_systems (id)
        ON DELETE CASCADE
        ON UPDATE CASCADE,
    CONSTRAINT users_office_phone_country_code_id_address_countries_id FOREIGN KEY (office_phone_country_code_id)
        REFERENCES public.address_countries (id)
        ON DELETE CASCADE
        ON UPDATE CASCADE,
    CONSTRAINT users_home_phone_country_code_id_address_countries_id FOREIGN KEY (home_phone_country_code_id)
        REFERENCES public.address_countries (id)
        ON DELETE CASCADE
        ON UPDATE CASCADE,
    CONSTRAINT users_cell_phone_country_code_id_address_countries_id FOREIGN KEY (cell_phone_country_code_id)
        REFERENCES public.address_countries (id)
        ON DELETE CASCADE
        ON UPDATE CASCADE,
    CONSTRAINT users_country_id_address_countries_id FOREIGN KEY (country_id)
        REFERENCES public.address_countries (id)
        ON DELETE CASCADE
        ON UPDATE CASCADE,
    CONSTRAINT users_province_id_address_provinces_id FOREIGN KEY (province_id)
        REFERENCES public.address_provinces (id)
        ON DELETE CASCADE
        ON UPDATE CASCADE,
    CONSTRAINT users_city_id_address_cities_id FOREIGN KEY (city_id)
        REFERENCES public.address_cities (id)
        ON DELETE CASCADE
        ON UPDATE CASCADE
    );`)

	return
}

// Rollback the database to previous version
func (Migration) RollbackUsersTable() (e error) {
	db := database.NewPGInstance()
	defer db.Close()
	_, e = db.Exec("DROP TABLE IF EXISTS public.users CASCADE;")

	return
}
