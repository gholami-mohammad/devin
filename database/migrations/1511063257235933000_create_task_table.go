package migrations

import "log"
import "gogit/database"
import "gogit/models"

// CreateTaskTable Migration Struct
type CreateTaskTable struct{}

// Migrate the database to a new version
func (CreateTaskTable) Migrate() {
	db := database.NewPGInstance()
	defer db.Close()
	_, e := db.Model(&models.Task{}).Exec(`CREATE TABLE IF NOT EXISTS ?TableName(
    id bigserial NOT NULL,
    status_id integer,
    type_id integer,
    created_by_id bigint,
    title varchar(512),
    description text,
    duration interval,
    scheduled_start_date timestamp with time zone,
    scheduled_completion_date timestamp with time zone ,
    started_on timestamp with time zone,
    completed_on timestamp with time zone ,
    completion_percentage numeric(3,3) ,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    deleted_at timestamp with time zone,

    CHECK (scheduled_start_date <= scheduled_completion_date),
    CHECK (started_on <= completed_on),
    CHECK (completion_percentage > 0 AND completion_percentage <= 100),

    CONSTRAINT tasks_pkey PRIMARY KEY (id),
    CONSTRAINT tasks_status_id_statuses_id FOREIGN KEY (status_id)
        REFERENCES project_management.task_statuses (id) MATCH SIMPLE
        ON DELETE RESTRICT
        ON UPDATE CASCADE,
    CONSTRAINT tasks_type_id_task_types_id FOREIGN KEY (type_id)
        REFERENCES project_management.task_types (id) MATCH SIMPLE
        ON DELETE RESTRICT
        ON UPDATE CASCADE,
    CONSTRAINT tasks_created_by_id_users_id FOREIGN KEY (created_by_id)
        REFERENCES public.users (id) MATCH SIMPLE
        ON DELETE RESTRICT
        ON UPDATE CASCADE
    )`)
	if e != nil {
		log.Println(e)
	}

}

// Rollback the database to previous version
func (CreateTaskTable) Rollback() {
	db := database.NewPGInstance()
	defer db.Close()
	_, e := db.Model(&models.Task{}).Exec("DROP TABLE IF EXISTS ?TableName CASCADE;")
	if e != nil {
		log.Println(e)
	}

}
