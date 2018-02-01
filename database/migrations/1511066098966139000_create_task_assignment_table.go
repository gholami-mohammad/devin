package migrations

import "log"
import "gogit/database"

func init() {
	// Migrations["CreateTaskAssignmentTable"] = CreateTaskAssignmentTable{}
}

// CreateTaskAssignmentTable Migration Struct
type CreateTaskAssignmentTable struct{}

// Migrate the database to a new version
func (CreateTaskAssignmentTable) Migrate() {
	db := database.NewPGInstance()
	defer db.Close()
	_, e := db.Exec(`CREATE TABLE IF NOT EXISTS ?TableName (
    id bigserial NOT NULL,
    task_id bigint NOT NULL,
    assign_to_user_id bigint NOT NULL,
    created_by_id bigint NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    deleted_at timestamp with time zone,
    CONSTRAINT task_assignments_pkey PRIMARY KEY (id),
    CONSTRAINT ... FOREIGN KEY () MATCH SIMPLE
        REFERENCES ... ()
        ON DELETE RESTRICT
        ON UPDATE CASCADE,
    CONSTRAINT ... FOREIGN KEY () MATCH SIMPLE
        REFERENCES ... ()
        ON DELETE RESTRICT
        ON UPDATE CASCADE,
    CONSTRAINT ... FOREIGN KEY () MATCH SIMPLE
        REFERENCES ... ()
        ON DELETE RESTRICT
        ON UPDATE CASCADE,
    )`)
	if e != nil {
		log.Println(e)
	}

}

// Rollback the database to previous version
func (CreateTaskAssignmentTable) Rollback() {
	db := database.NewPGInstance()
	defer db.Close()
	_, e := db.Exec(`DROP TABLE IF EXISTS ?TableName;`)
	if e != nil {
		log.Println(e)
	}

}
