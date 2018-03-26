package database

import (
	"github.com/go-pg/pg"
	_ "github.com/lib/pq"
)

// NewPGInstance opens new connection to Postgres database
func NewPGInstance() *pg.DB {
	return pg.Connect(&pg.Options{
		User:     "DBUser",
		Password: "DBPassword",
		Database: "gogit",
		Network:  "tcp",
		Addr:     "127.0.0.1:5432",
	})
}
