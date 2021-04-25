package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/gofrs/uuid"
	"github.com/pjheden/elephantio/config"
	"github.com/pjheden/elephantio/logs"
	"github.com/pjheden/elephantio/module"
	"github.com/pjheden/elephantio/task"

	// Necessary to get specifig driver that works with RPi
	_ "modernc.org/sqlite"
)

type Database struct {
	dbPath string
}

/*
	New makes one.
*/
func New(dbPath string) (*Database, error) {
	if _, err := os.Stat(dbPath); err != nil {
		return nil, fmt.Errorf("dbPath %q: %v", dbPath, err)
	}
	return &Database{
		dbPath: dbPath,
	}, nil
}

/*
Open uses the connection string to open the DB. Don't forget to close it!
*/
func (d *Database) Open() (*sql.DB, error) {
	return sql.Open("sqlite", d.dbPath)
}

/*
Modules returns the complete module including logs, config and task
*/
func (d *Database) Modules() ([]*module.Module, error) {
	// TODO: seperate which modules are active and which is not. so that we can have
	// old historical tasks saved in db
	query := sq.Select(
		"tasks.id",
		"task",
		"interval",
		"completed_on",
		"config.id",
		"pin_button",
		"pin_light",
	).
		From("tasks").
		Join("logs ON tasks.id = logs.task_id").
		Join("config on tasks.id = config.task_id")

	conn, err := d.Open()

	defer conn.Close()
	if err != nil {
		return nil, fmt.Errorf("connecting to db: %v", err)
	}
	rows, err := query.PlaceholderFormat(sq.Dollar).RunWith(conn).Query()
	if err != nil {
		return nil, fmt.Errorf("getting rows: %v", err)
	}

	ms := []*module.Module{}
	for rows.Next() {
		m := &module.Module{
			Task:   &task.Task{},
			Config: &config.Config{},
		}
		var timestamp string

		err := rows.Scan(
			&m.Task.ID,
			&m.Task.Name,
			&m.Task.Interval,
			&timestamp,
			&m.Config.ID,
			&m.Config.ButtonPin,
			&m.Config.LEDPin,
		)
		if err != nil {
			return nil, fmt.Errorf("scanning issues: %v", err)
		}

		// convert timestamp to time.Time
		log.Println("trying to parse ", timestamp)
		layout := "2006-01-02 04:05:06"
		t, err := time.Parse(layout, timestamp)

		if err != nil {
			return nil, err
		}
		m.CompletedOn = t

		ms = append(ms, m)
	}

	return ms, nil
}

/*
Tasks retrieves all tasks defined in table tasks and returns them in a slice of *task.Task
*/
func (d *Database) Tasks() ([]*task.Task, error) {
	query := sq.Select(
		"id",
		"task",
		"interval",
	).
		From("tasks")

	conn, err := d.Open()

	defer conn.Close()
	if err != nil {
		return nil, fmt.Errorf("connecting to db: %v", err)
	}
	rows, err := query.PlaceholderFormat(sq.Dollar).RunWith(conn).Query()
	if err != nil {
		return nil, fmt.Errorf("getting rows: %v", err)
	}

	allTasks := []*task.Task{}
	for rows.Next() {
		t := &task.Task{}
		err := rows.Scan(
			&(t.ID),
			&(t.Name),
			&(t.Interval),
		)
		if err != nil {
			return nil, fmt.Errorf("scanning issues: %v", err)
		}
		allTasks = append(allTasks, t)
	}

	return allTasks, nil
}

/*
Logs retrieves all Logs defined in table Logs and returns them in a slice of *logs.Task
*/
func (d *Database) Logs() ([]*logs.Logs, error) {
	query := sq.Select(
		"id",
		"task_id",
		"completed_on",
	).
		From("logs")

	conn, err := d.Open()

	defer conn.Close()
	if err != nil {
		return nil, fmt.Errorf("connecting to db: %v", err)
	}
	rows, err := query.PlaceholderFormat(sq.Dollar).RunWith(conn).Query()
	if err != nil {
		return nil, fmt.Errorf("getting rows: %v", err)
	}

	allLogs := []*logs.Logs{}
	for rows.Next() {
		l := &logs.Logs{}
		err := rows.Scan(
			&(l.ID),
			&(l.TaskID),
			&(l.CompletedOn),
		)
		if err != nil {
			return nil, fmt.Errorf("scanning issues: %v", err)
		}
		allLogs = append(allLogs, l)
	}

	return allLogs, nil
}

/*
Confisg retrieves all Confisg defined in table Config and returns them in a slice of *config.Config
*/
func (d *Database) Configs() ([]*config.Config, error) {
	query := sq.Select(
		"id",
		"task_id",
		"pin_button",
		"pin_light",
	).
		From("config")

	conn, err := d.Open()

	defer conn.Close()
	if err != nil {
		return nil, fmt.Errorf("connecting to db: %v", err)
	}
	rows, err := query.PlaceholderFormat(sq.Dollar).RunWith(conn).Query()
	if err != nil {
		return nil, fmt.Errorf("getting rows: %v", err)
	}

	allConfig := []*config.Config{}
	for rows.Next() {
		c := &config.Config{}
		err := rows.Scan(
			&(c.ID),
			&(c.TaskID),
			&(c.ButtonPin),
			&(c.LEDPin),
		)
		if err != nil {
			return nil, fmt.Errorf("scanning issues: %v", err)
		}
		allConfig = append(allConfig, c)
	}

	return allConfig, nil
}

// TODO Get all tasks in combination with logs and task

// TODO add completed task to logs

/*
AddLogs adds a row to logs
*/
func (d *Database) AddLogs(t *task.Task) error {
	// Create a Version 4 UUID.
	u2, err := uuid.NewV4()
	if err != nil {
		return err
	}

	query := sq.Insert(
		"logs",
	).Columns(
		"id",
		"task_id",
	).Values(
		u2,
		t.ID,
	)

	conn, err := d.Open()

	defer conn.Close()
	if err != nil {
		return fmt.Errorf("connecting to db: %v", err)
	}
	_, err = query.PlaceholderFormat(sq.Dollar).RunWith(conn).Query()
	if err != nil {
		return fmt.Errorf("inserting row: %v", err)
	}

	log.Println("Inserted %+v into logs", t)
	return nil
}
