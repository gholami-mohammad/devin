package migrations

import (
	"log"
)

// Migrate database
func Migrate() {
	log.SetFlags(log.Lshortfile)

	// TODO: DELETE THESE 2 Lines
	CreatePublicSchema{}.Rollback()
	CreateProjectManagementSchema{}.Rollback()

	CreatePublicSchema{}.Migrate()
	CreateProjectManagementSchema{}.Migrate()
	CreateUserTable{}.Migrate()
	CreateTaskStatusTable{}.Migrate()
	CreateTaskTypeTable{}.Migrate()
	CreateTaskTable{}.Migrate()

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

// Rollback all migrations
func Rollback() {
	CreateTaskTable{}.Rollback()
	CreateTaskTypeTable{}.Rollback()
	CreateTaskStatusTable{}.Rollback()
	CreateUserTable{}.Rollback()

	CreateProjectManagementSchema{}.Rollback()
	CreatePublicSchema{}.Rollback()

}
