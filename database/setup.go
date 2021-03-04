package database

import (
	"database/sql"
	"errors"
	"log"
	"os"
)

/*
 fileExists checks if a file exists and is not a directory before we
 try using it to prevent further errors.
*/
func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

/*
InitDB sets up a new sqlite db with the given path
It can then be used by sqlite.New()
*/
func InitDB(dbPath string) error {
	if fileExists(dbPath) {
		return errors.New("DB already exists, remove existing sqlite db first")
	}
	stmts := [...]string{
		`create table tasks (
			id uuid primary key,
			task text not null unique,
			interval real DEFAULT (24.0)
		)`,
		`create table logs (
			id uuid primary key,
			task_id uuid not null,
			completed_on datetime DEFAULT (datetime('now','utc'))
		)`,
		`create table config (
			id uuid primary key,
			task_id uuid not null,
			pin_1 int not null,
			pin_2 int not null
		)`,

		// Insert some initial data

		// knee
		`insert into tasks (id, task) VALUES (
			"0503ec8c-ab35-44ed-b2af-30cc772196e7",
			"knee"
		)`,
		`insert into logs (id, task_id) VALUES (
			"2054845f-4bc7-429b-bdc8-723ae74a5c41",
			"0503ec8c-ab35-44ed-b2af-30cc772196e7"
		)`,
		`insert into config (id, task_id, pin_1, pin_2) VALUES (
			"1a7abd5f-233c-4efe-b238-a2ff329efb78",
			"0503ec8c-ab35-44ed-b2af-30cc772196e7",
			6,
			7
		)`,
	}

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return err
	}
	defer db.Close()

	for _, stmt := range stmts {
		_, err := db.Exec(stmt)
		if err != nil {
			log.Printf("%q: %s\n", err, stmt)
			return err
		}
	}

	return nil
}
