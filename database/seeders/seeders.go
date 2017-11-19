package seeders

import (
	"log"

	"github.com/go-pg/pg"

	"gogit/database"
	"gogit/models"
)

var e error

func Seed() {
	db := database.NewPGInstance()
	defer db.Close()
	db.OnQueryProcessed(func(ev *pg.QueryProcessedEvent) {
		// log.Println(ev.FormattedQuery())
	})

	taskStatusSeeder(db)
}

func setSeqID(model interface{}) {
	db := database.NewPGInstance()
	defer db.Close()
	db.OnQueryProcessed(func(ev *pg.QueryProcessedEvent) {
		// log.Println(ev.FormattedQuery())
	})
	db.Model(model).Exec("SELECT setval('?TableName' || '_id_seq', max(id)) from ?TableName;")
}

func taskStatusSeeder(db *pg.DB) {
	_, e = db.Model(&models.TaskStatus{}).Exec(`INSERT INTO ?TableName (id, title, created_at, updated_at, deleted_at) VALUES (1, 'Not started', now() , now() , null) ON CONFLICT (id) DO NOTHING`)
	if e != nil {
		log.Println(e)
	}
	_, e = db.Model(&models.TaskStatus{}).Exec(`INSERT INTO ?TableName (id, title, created_at, updated_at, deleted_at) VALUES (2, 'In Progress', now() , now() , null) ON CONFLICT (id) DO NOTHING`)
	if e != nil {
		log.Println(e)
	}
	_, e = db.Model(&models.TaskStatus{}).Exec(`INSERT INTO ?TableName (id, title, created_at, updated_at, deleted_at) VALUES (3, 'Completed', now() , now() , null) ON CONFLICT (id) DO NOTHING`)
	if e != nil {
		log.Println(e)
	}
	db.Model(&models.TaskStatus{}).Exec(`SELECT set_val ?TableName (id, title, created_at, updated_at, deleted_at) VALUES (1, 'Completed', now() , now() , null) ON CONFLICT (id) DO NOTHING`)

	setSeqID(&models.TaskStatus{})

}

func taskTypeSeeder(db *pg.DB) {
	_, e = db.Model(&models.TaskType{}).Exec(`INSERT INTO ?TableName (id, name, created_at, updated_at, deleted_at VALUES (1, 'Project' , now() , now() , null)) ON CONFLICT (id) DO NOTHING`)
	if e != nil {
		log.Println(e)
	}
	_, e = db.Model(&models.TaskType{}).Exec(`INSERT INTO ?TableName (id, name, created_at, updated_at, deleted_at VALUES (2, 'Section OR Module' , now() , now() , null)) ON CONFLICT (id) DO NOTHING`)
	if e != nil {
		log.Println(e)
	}
	_, e = db.Model(&models.TaskType{}).Exec(`INSERT INTO ?TableName (id, name, created_at, updated_at, deleted_at VALUES (3, 'Mileston' , now() , now() , null)) ON CONFLICT (id) DO NOTHING`)
	if e != nil {
		log.Println(e)
	}
	_, e = db.Model(&models.TaskType{}).Exec(`INSERT INTO ?TableName (id, name, created_at, updated_at, deleted_at VALUES (4, 'Task' , now() , now() , null)) ON CONFLICT (id) DO NOTHING`)
	if e != nil {
		log.Println(e)
	}
	_, e = db.Model(&models.TaskType{}).Exec(`INSERT INTO ?TableName (id, name, created_at, updated_at, deleted_at VALUES (, 'Sub task' , now() , now() , null)) ON CONFLICT (id) DO NOTHING`)
	if e != nil {
		log.Println(e)
	}

	setSeqID(&models.TaskType{})

}
