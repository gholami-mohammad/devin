package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"devin/cmd/helpers"
)

func main() {

	flag.Parse()
	command := flag.Arg(0)

	if strings.EqualFold(command, "") {
		fmt.Println("No command specified")
		os.Exit(1)
	}
	fmt.Println("loading ", command)

	switch command {

	case "migrate":
		{
			helpers.Migrate()
			fmt.Println("Database migrated")
		}
	case "db:seed":
		{
			fmt.Println("Seed finished")
		}
	case "migrate:rollback":
		{
			helpers.Rollback()
		}
	case "make:migration":
		{
			set := flag.NewFlagSet("make:migrate", flag.ContinueOnError)
			set.Parse(os.Args[2:])
			name := set.Arg(0)

			if strings.EqualFold(name, "") {
				os.Exit(1)
				return
			}
			helpers.MakeMigration(&name)
		}
	default:
		{
			fmt.Println("Command not found :( ")
		}
	}

}
