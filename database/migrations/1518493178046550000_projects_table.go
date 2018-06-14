package migrations

import "devin/database"

// Migrate the database to a new version
func (Migration) MigrateProjectsTable() (e error) {
	db := database.NewGORMInstance()
	defer db.Close()
	e = db.Exec(`CREATE TABLE IF NOT EXISTS public.projects (
    id                              bigserial NOT NULL,
    name                            varchar(200) NOT NULL,
    title                           varchar(200),
    description                     text,
    scheduled_start_date            timestamp with time zone,
    start_date                      timestamp with time zone,
    scheduled_completion_date       timestamp with time zone,
    completion_date                 timestamp with time zone,
    default_task_view_id            integer DEFAULT 1,
    status_id                       integer,
    owner_user_id                   bigint NOT NULL,
    owner_organization_id           bigint ,
    project_manager_id              bigint,
    created_by_id                   bigint NOT NULL,

    enable_wiki_module              bool,
    allow_public_wiki               bool,
    enable_tasks_module             bool,
    enale_milestones_module         bool,
    enable_files_module             bool,
    enable_messages_module          bool,
    enable_time_logs_module         bool,
    enable_notebooks_module         bool,
    enable_risks_module             bool,
    enable_links_module             bool,
    enable_billing_module           bool,
    enable_git_module               bool,
    enable_issue_tracker_module     bool,
    allow_public_issues             bool,
    enable_bug_tracker_module       bool,
    allow_public_bugs               bool,
    enable_project_comments         bool,

    created_at                      timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at                      timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    deleted_at                      timestamp with time zone,

    CHECK (scheduled_start_date <= scheduled_completion_date),
    CHECK (start_date <= completion_date),
    CHECK (default_task_view_id = 1 OR default_task_view_id = 2),

    CONSTRAINT projects_pkey PRIMARY KEY (id),
    CONSTRAINT projects_key UNIQUE (name),

    CONSTRAINT projects_status_id_project_statuses_id FOREIGN KEY (status_id)
        REFERENCES public.project_statuses (id) MATCH SIMPLE
        ON DELETE CASCADE
        ON UPDATE CASCADE,

    CONSTRAINT projects_owner_user_id_users_id FOREIGN KEY (owner_user_id)
        REFERENCES public.users (id) MATCH SIMPLE
        ON DELETE CASCADE
        ON UPDATE CASCADE,

    CONSTRAINT projects_owner_organization_id_users_id FOREIGN KEY (owner_organization_id)
        REFERENCES public.users (id) MATCH SIMPLE
        ON DELETE CASCADE
        ON UPDATE CASCADE,

    CONSTRAINT projects_project_manager_id_users_id FOREIGN KEY (project_manager_id)
        REFERENCES public.users (id) MATCH SIMPLE
        ON DELETE CASCADE
        ON UPDATE CASCADE,

    CONSTRAINT projects_created_by_id_users_id FOREIGN KEY (created_by_id)
        REFERENCES public.users (id) MATCH SIMPLE
        ON DELETE CASCADE
        ON UPDATE CASCADE
    );`).Error

	return
}

// Rollback the database to previous version
func (Migration) RollbackProjectsTable() (e error) {
	db := database.NewGORMInstance()
	defer db.Close()
	e = db.Exec("DROP TABLE IF EXISTS public.projects CASCADE;").Error

	return
}
