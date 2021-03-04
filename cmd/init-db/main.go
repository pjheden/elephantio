package main

import (
	"log"

	"github.com/pjheden/elephantio/database"
)

func main() {
	dbPath := "./db.sqlite"

	err := database.InitDB(dbPath)
	if err != nil {
		log.Fatal(err)
	}

	_, err = database.New(dbPath)
	if err != nil {
		log.Fatal(err)
	}

}
