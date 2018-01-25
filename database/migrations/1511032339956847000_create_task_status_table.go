package migrations

import "log"
import "gogit/database"
import "gogit/models"

func init() {
	// Migrations["CreateTaskStatusTable"] = CreateTaskStatusTable{}
}

// CreateTaskStatusTable Migration Struct
type CreateTaskStatusTable struct{}

// Migrate the database to a new version
func (CreateTaskStatusTable) Migrate() {
	db := database.NewPGInstance()
	defer db.Close()
	_, e := db.Model(&models.TaskStatus{}).Exec(`CREATE TABLE IF NOT EXISTS ?TableName(
    id serial NOT NULL,
    title varchar(255),
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    deleted_at timestamp with time zone,
    CONSTRAINT task_statuses_pkey PRIMARY KEY (id)
    )`)
	if e != nil {
		log.Println(e)
	}

}

// Rollback the database to previous version
func (CreateTaskStatusTable) Rollback() {
	db := database.NewPGInstance()
	defer db.Close()
	_, e := db.Model(&models.TaskStatus{}).Exec("DROP TABLE IF EXISTS ?TableName CASCADE;")
	if e != nil {
		log.Println(e)
	}

}
