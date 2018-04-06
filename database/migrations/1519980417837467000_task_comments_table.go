package migrations

import "devin/database"

// Migrate the database to a new version
func (Migration) MigrateTaskCommentsTable() (e error) {
	db := database.NewGORMInstance()
	defer db.Close()
	e = db.Exec(`CREATE TABLE IF NOT EXISTS public.task_comments (
    id bigserial NOT NULL,
    task_id bigint NOT NULL,
    reply_to_id bigint,
    comment text NOT NULL,
    attachment_path varchar(255),
    created_by_id bigint NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    deleted_at timestamp with time zone,

    CONSTRAINT task_comments_pkey PRIMARY KEY(id),
    CONSTRAINT task_comments_task_id_tasks_id FOREIGN KEY (task_id)
        REFERENCES public.tasks (id) MATCH SIMPLE
        ON DELETE CASCADE
        ON UPDATE CASCADE,
    CONSTRAINT task_comments_reply_to_id_tasks_id FOREIGN KEY (reply_to_id)
        REFERENCES public.task_comments (id) MATCH SIMPLE
        ON DELETE CASCADE
        ON UPDATE CASCADE,
    CONSTRAINT task_comments_created_by_id_users_id FOREIGN KEY (created_by_id)
        REFERENCES public.users (id) MATCH SIMPLE
        ON DELETE CASCADE
        ON UPDATE CASCADE
    )`).Error

	return
}

// Rollback the database to previous version
func (Migration) RollbackTaskCommentsTable() (e error) {
	db := database.NewGORMInstance()
	defer db.Close()
	e = db.Exec("DROP TABLE IF EXISTS public.task_comments;").Error

	return
}
