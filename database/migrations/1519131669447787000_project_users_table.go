package migrations

import "devin/database"

// Migrate the database to a new version
func (Migration) MigrateProjectUsersTable() (e error) {
	db := database.NewPGInstance()
	defer db.Close()
	_, e = db.Exec(`CREATE TABLE IF NOT EXISTS public.project_users(
    id                                  bigserial NOT NULL,
    user_id                             bigint NOT NULL,
    project_id                          bigint NOT NULL,
    created_by_id                       bigint NOT NULL,
    id_admin                            bool,
    can_update_prject_profile           bool,
    can_add_user_to_project             bool,
    can_create_milestone                bool,
    can_create_task_list                bool,
    can_create_task                     bool,
    can_create_issue                    bool,
    can_create_repository               bool,
    can_create_tag                      bool,
    can_create_board                    bool,
    can_create_reminder                 bool,
    can_create_time_log                 bool,
    can_list_all_time_logs              bool,
    can_create_wiki                     bool,
    created_at                          timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at                          timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    deleted_at                          timestamp with time zone,

    CONSTRAINT project_users_pkey PRIMARY KEY (id),

    CONSTRAINT project_users_user_id_users_id FOREIGN KEY (user_id)
        REFERENCES public.users (id) MATCH SIMPLE
        ON DELETE CASCADE
        ON UPDATE CASCADE,

    CONSTRAINT project_users_project_id_users_id FOREIGN KEY (project_id)
        REFERENCES public.projects (id) MATCH SIMPLE
        ON DELETE CASCADE
        ON UPDATE CASCADE,
        
    CONSTRAINT project_users_created_by_id_users_id FOREIGN KEY (created_by_id)
        REFERENCES public.users (id) MATCH SIMPLE
        ON DELETE CASCADE
        ON UPDATE CASCADE
    )`)

	return
}

// Rollback the database to previous version
func (Migration) RollbackProjectUsersTable() (e error) {
	db := database.NewPGInstance()
	defer db.Close()
	_, e = db.Exec("DROP TABLE IF EXISTS public.project_users;")

	return
}
