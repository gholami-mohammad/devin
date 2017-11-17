package main

import (
	"fmt"
	"os"
	"strings"

	"gogit/database/migrations"
	"gogit/database/seeders"
)

func main() {
	var cmd string
	if len(os.Args) > 1 {
		cmd = os.Args[1]
	} else {
		fmt.Println("Enter your command: ")
		fmt.Scanln(&cmd)
	}

	fmt.Println("loading ", cmd)
	cmd = strings.ToLower(cmd)
	switch cmd {
	case "migrate":
		{
			migrations.Migrate()
			fmt.Println("Database migrated")
		}
	case "db:seed":
		{
			seeders.Seed()
			fmt.Println("Seed finished")
		}
	case "migrate:rollback":
		{
			// migrations.Rollback()
			// fmt.Println("Database Rolled back")
		}
	}

}
