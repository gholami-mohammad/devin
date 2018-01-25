package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"gogit/database/seeders"

	"gogit/cmd/helpers"
)

var (
	create  *string
	command *string
)

// Basic flag declarations are available for string, integer, and boolean options.
func init() {
	command = flag.String("run", "", "The command")
	create = flag.String("create", "migration", "The file name to create")
}

func main() {

	flag.Parse()

	if strings.EqualFold(*command, "") {
		fmt.Println("No command specified")
		os.Exit(1)
	}
	fmt.Println("loading ", *command)

	switch *command {

	case "migrate":
		{
			helpers.Migrate()
			fmt.Println("Database migrated")
		}
	case "db:seed":
		{
			seeders.Seed()
			fmt.Println("Seed finished")
		}
	case "migrate:rollback":
		{
			fmt.Println("Database Rollbacked")
		}
	case "make:migration":
		{
			helpers.MakeMigration(create)
		}
	default:
		{
			fmt.Println("Command not found :( ")
		}
	}

}
