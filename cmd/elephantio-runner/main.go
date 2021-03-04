package main

import (
	"log"

	"github.com/pjheden/elephantio/database"
)

func main() {

	dbPath := "./db.sqlite"

	db, err := database.New(dbPath)
	if err != nil {
		log.Fatal(err)
	}

	fts, err := db.FullTasks()
	if err != nil {
		log.Fatal(err)
	}

	for _, ft := range fts {
		log.Printf("got fts %+v\n", ft)

		// TODO: setup pi connections
	}

}
