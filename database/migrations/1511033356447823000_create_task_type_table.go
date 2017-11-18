package migrations

import "log"
import "gogit/database"
import "gogit/models"

// CreateTaskTypeTable Migration Struct
type CreateTaskTypeTable struct{}

// Migrate the database to a new version
func (CreateTaskTypeTable) Migrate() {
	db := database.NewPGInstance()
	defer db.Close()
	_, e := db.Model(&models.TaskType{}).Exec(`CREATE TABLE IF NOT EXISTS ?TableName (
    id serial NOT NULL,
    name varchar(255),
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    deleted_at timestamp with time zone,
    CONSTRAINT task_types_pkey PRIMARY KEY (id)
    )`)
	if e != nil {
		log.Println(e)
	}

}

// Rollback the database to previous version
func (CreateTaskTypeTable) Rollback() {
	db := database.NewPGInstance()
	defer db.Close()
	_, e := db.Model(&models.TaskType{}).Exec("DROP TABLE IF EXISTS ?TableName")
	if e != nil {
		log.Println(e)
	}

}
