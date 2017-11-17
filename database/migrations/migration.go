package migrations

import (
	"log"

	"github.com/go-pg/pg/orm"

	"gogit/database"
	"gogit/models"
)

func Migrate() {
	log.SetFlags(log.Lshortfile)

	db := database.NewPGInstance()
	defer db.Close()
	var e error
	opt := orm.CreateTableOptions{
		FKConstraints: true,
		IfNotExists:   true,
	}

	db.Exec(`DROP SCHEMA public CASCADE;`)
	db.Exec(`CREATE SCHEMA public;`)

	db.Exec(`DROP SCHEMA project_management CASCADE;`)
	db.Exec(`CREATE SCHEMA project_management;`)

	e = db.CreateTable(&models.User{}, &opt)
	if e != nil {
		log.Println(e)
	}

	e = db.CreateTable(&models.TaskStatus{}, &opt)
	if e != nil {
		log.Println(e)
	}
	e = db.CreateTable(&models.TaskType{}, &opt)
	if e != nil {
		log.Println(e)
	}
	e = db.CreateTable(&models.Task{}, &opt)
	if e != nil {
		log.Println(e)
	}
	e = db.CreateTable(&models.TaskAssignment{}, &opt)
	if e != nil {
		log.Println(e)
	}
	e = db.CreateTable(&models.TaskAttachment{}, &opt)
	if e != nil {
		log.Println(e)
	}
	e = db.CreateTable(&models.TaskComment{}, &opt)
	if e != nil {
		log.Println(e)
	}
	e = db.CreateTable(&models.TaskDiscussion{}, &opt)
	if e != nil {
		log.Println(e)
	}
	e = db.CreateTable(&models.TaskPredecessor{}, &opt)
	if e != nil {
		log.Println(e)
	}
	e = db.CreateTable(&models.TaskReminder{}, &opt)
	if e != nil {
		log.Println(e)
	}
}
