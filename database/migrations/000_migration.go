package migrations

import (
	"log"
)

// Migrate , Load all migrations functions
func Migrate() {
	log.SetFlags(log.Lshortfile)

	CreatePublicSchema{}.Migrate()
	CreateProjectManagementSchema{}.Migrate()
	CreateUserTable{}.Migrate()

	// e = db.CreateTable(&models.TaskStatus{}, &opt)
	// if e != nil {
	// 	log.Println(e)
	// }
	// e = db.CreateTable(&models.TaskType{}, &opt)
	// if e != nil {
	// 	log.Println(e)
	// }
	// e = db.CreateTable(&models.Task{}, &opt)
	// if e != nil {
	// 	log.Println(e)
	// }
	// e = db.CreateTable(&models.TaskAssignment{}, &opt)
	// if e != nil {
	// 	log.Println(e)
	// }
	// e = db.CreateTable(&models.TaskAttachment{}, &opt)
	// if e != nil {
	// 	log.Println(e)
	// }
	// e = db.CreateTable(&models.TaskComment{}, &opt)
	// if e != nil {
	// 	log.Println(e)
	// }
	// e = db.CreateTable(&models.TaskDiscussion{}, &opt)
	// if e != nil {
	// 	log.Println(e)
	// }
	// e = db.CreateTable(&models.TaskPredecessor{}, &opt)
	// if e != nil {
	// 	log.Println(e)
	// }
	// e = db.CreateTable(&models.TaskReminder{}, &opt)
	// if e != nil {
	// 	log.Println(e)
	// }
}

func Rollback() {
	CreateUserTable{}.Rollback()

	CreateProjectManagementSchema{}.Rollback()
	CreatePublicSchema{}.Rollback()

}
