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

func taskStatusSeeder(db *pg.DB) {
	_, e = db.Model(&models.TaskStatus{}).Exec(`INSERT INTO ?TableName (id, title, created_at, updated_at, deleted_at) VALUES (1, 'Not started', now() , now() , null) ON CONFLICT (id) DO NOTHING`)
	if e != nil {
		log.Println(e)
	}
	_, e = db.Model(&models.TaskStatus{}).Exec(`INSERT INTO ?TableName (id, title, created_at, updated_at, deleted_at) VALUES (1, 'In Progress', now() , now() , null) ON CONFLICT (id) DO NOTHING`)
	if e != nil {
		log.Println(e)
	}
	_, e = db.Model(&models.TaskStatus{}).Exec(`INSERT INTO ?TableName (id, title, created_at, updated_at, deleted_at) VALUES (1, 'Completed', now() , now() , null) ON CONFLICT (id) DO NOTHING`)
	if e != nil {
		log.Println(e)
	}

}