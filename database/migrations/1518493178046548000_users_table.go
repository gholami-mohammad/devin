package migrations

import "gogit/database"

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
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    deleted_at timestamp with time zone,

    CONSTRAINT users_pkey PRIMARY KEY (id)
    CONSTRAINT users_username_key UNIQUE (username)
    CONSTRAINT users_email_key UNIQUE (email)
    CONSTRAINT users_owner_id_users_id FOREIGN KEY (owner_id)
        REFERENCES public.users (id)
        ON DELETE CASCADE
        ON UPDATE CASCADE
    );`)

	return
}

// Rollback the database to previous version
func (Migration) RollbackUsersTable() (e error) {
	db := database.NewPGInstance()
	defer db.Close()
	_, e = db.Exec("DROP TABLE IF EXISTS public.users;")

	return
}
