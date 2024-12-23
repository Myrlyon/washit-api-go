package main

import (
	"log"
	"os"
	"washit-api/pkg/configs"
	"washit-api/pkg/db/dbs"
	"washit-api/pkg/utils"
)

func main() {
	db, err := dbs.NewDatabase(configs.Envs.URI)
	if err != nil {
		log.Fatal("Failed to connect to the database", err)
	}

	cmd := os.Args[len(os.Args)-1]
	if cmd == "up" {
		db.Migrate(utils.ModelList...)
	}
	if cmd == "down" {
		db.DropTable(utils.ModelList...)
	}
}
