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
			completed_on TIMESTAMP 
			DEFAULT CURRENT_TIMESTAMP
		)`,
		`create table config (
			id uuid primary key,
			task_id uuid not null,
			pin_button int not null,
			pin_light int not null
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
		`insert into config (id, task_id, pin_button, pin_light) VALUES (
			"1a7abd5f-233c-4efe-b238-a2ff329efb78",
			"0503ec8c-ab35-44ed-b2af-30cc772196e7",
			20,
			21
		)`,

		// water plants
		`insert into tasks (id, task) VALUES (
			"d0851ce6-807f-4e15-b84c-62ccae79d5bd",
			"water plants"
		)`,
		`insert into logs (id, task_id) VALUES (
			"c1cabb55-a56c-4264-8f52-90727971dd46",
			"d0851ce6-807f-4e15-b84c-62ccae79d5bd"
		)`,
		`insert into config (id, task_id, pin_button, pin_light) VALUES (
			"0acfab8b-924c-408d-8ca4-e6e951ffa8e1",
			"d0851ce6-807f-4e15-b84c-62ccae79d5bd",
			5,
			27
		)`,

		// d-vitamin
		`insert into tasks (id, task) VALUES (
			"d7e94418-6dcb-40d7-86e2-ae7b14472e53",
			"d-vitamin"
		)`,
		`insert into logs (id, task_id) VALUES (
			"211643c1-3310-45dd-bbbd-347ea78a4e14",
			"d7e94418-6dcb-40d7-86e2-ae7b14472e53"
		)`,
		`insert into config (id, task_id, pin_button, pin_light) VALUES (
			"c78fedc4-8d58-444f-a075-574ce9ff9056",
			"d7e94418-6dcb-40d7-86e2-ae7b14472e53",
			3,
			2
		)`,
	}

	db, err := sql.Open("sqlite", dbPath)
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
