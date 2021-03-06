package helpers

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"reflect"
	"strings"
	"time"

	"github.com/iancoleman/strcase"

	"devin/database"
	"devin/database/migrations"
)

func init() {
	log.SetFlags(log.Lshortfile)
}

type migration struct {
	ID    uint
	Name  string
	Batch uint
}

// MakeMigration create new migration file
func MakeMigration(create *string) {
	timestamp := time.Now().UnixNano()
	filename := fmt.Sprintf("./database/migrations/%v_%v.go", timestamp, strcase.ToSnake(*create))
	content := fmt.Sprintf(`package migrations

import "devin/database"

// Migrate the database to a new version
func (Migration) Migrate%v() (e error) {
    db := database.NewGORMInstance()
    defer db.Close()
    e = db.Exec("").Error

    return
}

// Rollback the database to previous version
func (Migration) Rollback%v() (e error) {
    db := database.NewGORMInstance()
    defer db.Close()
    e = db.Exec("").Error

    return
}`, strcase.ToCamel(*create), strcase.ToCamel(*create))
	f, e := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, 0777)
	if e != nil {
		fmt.Println(e)
		os.Exit(1)
	}
	defer f.Close()
	f.WriteString(content)
}

// checkMigrationTable will check existance of migrations table and create it if it dosen't exist.
func checkMigrationTable() {
	db := database.NewGORMInstance()
	defer db.Close()

	db.Exec(`CREATE SCHEMA IF NOT EXISTS public;`)
	db.Exec(`CREATE TABLE IF NOT EXISTS migrations (
    id serial NOT NULL,
    name varchar(255) NOT NULL,
    batch integer,
    CONSTRAINT migrations_pkey PRIMARY KEY (id)
    )`)
}

// Migrate call Migrate function of files and save batches in database
func Migrate() {
	checkMigrationTable()

	db := database.NewGORMInstance()
	defer db.Close()

	files, e := ioutil.ReadDir("./database/migrations/")
	if e != nil {
		Printer{}.Error("Error on loading migrations directory")
	}

	var batch struct {
		LastBatch uint `sql:"last_batch"`
	}
	db.Raw(`select max(batch) as last_batch from migrations;`).Scan(&batch)
	batch.LastBatch++

	for _, v := range files {
		if v.IsDir() {
			continue
		}
		var mg migration
		filename := strings.TrimSuffix(v.Name(), ".go")
		undescoreIndex := strings.Index(filename, "_")
		filename = filename[undescoreIndex+1:]
		name := strcase.ToCamel(filename)
		if strings.EqualFold(name, "Migration") {
			continue
		}
		db.Where("name LIKE ?", filename).First(&mg)

		if mg.ID != 0 {
			//This file already migrated
			continue
		}
		m := migrations.Migration{}
		val := reflect.ValueOf(m)
		funcName := "Migrate" + name
		f := val.MethodByName(funcName)
		if !f.IsValid() {
			Printer{}.Error("Invalid migration funciton name: Migrate" + name)
			continue
		}

		if f.Type().NumOut() == 0 {
			Printer{}.Error("Function must return at least one value")
			continue
		}
		lastReturn := f.Type().Out(f.Type().NumOut() - 1).String()
		if !strings.EqualFold(lastReturn, "error") {
			Printer{}.Error("Last return value of function must be of type `error`. ", "Value of type `", lastReturn, "` returned")
			continue
		}

		rets := f.Call(nil)
		ee := rets[f.Type().NumOut()-1].Interface()

		if rets[f.Type().NumOut()-1].Interface() != nil {
			Printer{}.Error(funcName, " exited with error: ", ee)
			continue
		}

		// save migrated file to DB
		mg.Name = filename
		mg.Batch = batch.LastBatch
		db.Create(&mg)

		Printer{}.Success(name, " Migrated")
	}

	db.Exec("update migrations set batch=coalesce((select max(batch) from migrations) , 0)+1 where batch is null;")

}

// Rollback will rollback database using batch number
func Rollback() {
	checkMigrationTable()

	db := database.NewGORMInstance()
	defer db.Close()

	var rollbacks []migration
	db.Where("batch = (select max(batch) from migrations)").Find(&rollbacks)
	for _, v := range rollbacks {
		name := strcase.ToCamel(v.Name)
		Printer{}.Warning("Rollback ", name)
		m := migrations.Migration{}
		val := reflect.ValueOf(m)
		f := val.MethodByName("Rollback" + name)
		if !f.IsValid() {
			Printer{}.Error("Invalid rollback funciton name: Rollback" + name)
			continue
		}

		if f.Type().NumOut() == 0 {
			Printer{}.Error("Function must return at least one value")
			continue
		}
		lastReturn := f.Type().Out(f.Type().NumOut() - 1).String()
		if !strings.EqualFold(lastReturn, "error") {
			Printer{}.Error("Last return value of function must be of type `error`. ", "Value of type `", lastReturn, "` returned")
			continue
		}

		rets := f.Call(nil)
		if len(rets) == 0 {
			continue
		}

		if rets[0].Interface() != nil {
			Printer{}.Error(rets[0].Interface())
			continue
		}

		db.Delete(&v)
	}

}
