package database

import (
	// "github.com/go-pg/pg"
	"ems_ms/database"
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	// _ "github.com/lib/pq"
)

// NewPGInstance opens new connection to Postgres database
// func NewPGInstance() *pg.DB {
// 	return pg.Connect(&pg.Options{
// 		User:     "mgh",
// 		Password: "mgh_ua6872",
// 		Database: "gogit",
// 		Network:  "tcp",
// 		Addr:     "127.0.0.1:5432",
// 	})
// }

// NewGORMInstance create DB instance using gorm
func NewGORMInstance() *gorm.DB {
	con, e := gorm.Open("postgres", "user=mgh password=mgh_ua6872 host=127.0.0.1 dbname=gogit port=5432 sslmode=disable")
	if e != nil {
		return nil
	}

	return con
}

func AddSeqID(tableName string) {
	db := database.NewInstance()
	defer db.Close()
	db.Exec(fmt.Sprintf("SELECT setval('%v_id_seq', max(id)) FROM %v;", tableName, tableName))
}
