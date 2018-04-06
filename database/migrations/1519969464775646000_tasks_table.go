package migrations

import "devin/database"

// Migrate the database to a new version
func (Migration) MigrateTasksTable() (e error) {
	db := database.NewGORMInstance()
	defer db.Close()
	e = db.Exec(`CREATE TABLE IF NOT EXISTS public.tasks (
    id bigserial NOT NULL,
    title varchar (255),
    order_id integer DEFAULT 1,
    description text,
    scheduled_start_date timestamp with time zone,
    scheduled_completion_date timestamp with time zone,
    start_date timestamp with time zone,
    completion_date timestamp with time zone,
    completed_successfully bool DEFAULT true,
    priority_id integer,
    font_color varchar(25),
    background_color varchar(25),
    progress smallint,
    estimated_time interval,
    task_board_id bigint,
    created_by_id bigint NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    deleted_at timestamp with time zone,

    CHECK (progress >= 0 AND progress <= 100),

    CONSTRAINT tasks_pkey PRIMARY KEY (id),
    CONSTRAINT tasks_priority_id_task_priorities_id FOREIGN KEY (priority_id)
        REFERENCES public.task_priorities (id) MATCH SIMPLE
        ON DELETE CASCADE
        ON UPDATE CASCADE,
    CONSTRAINT tasks_task_board_id_task_boards_id FOREIGN KEY (task_board_id)
        REFERENCES public.task_boards (id) MATCH SIMPLE
        ON DELETE CASCADE
        ON UPDATE CASCADE,
    CONSTRAINT tasks_created_by_id_users_id FOREIGN KEY (created_by_id)
        REFERENCES public.users (id) MATCH SIMPLE
        ON DELETE CASCADE
        ON UPDATE CASCADE
    )`).Error

	return
}

// Rollback the database to previous version
func (Migration) RollbackTasksTable() (e error) {
	db := database.NewGORMInstance()
	defer db.Close()
	e = db.Exec("DROP TABLE IF EXISTS public.tasks CASCADE;").Error

	return
}
