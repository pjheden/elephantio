package database

import (
	"database/sql"
	"fmt"
	"os"

	sq "github.com/Masterminds/squirrel"
	_ "github.com/mattn/go-sqlite3"
	"github.com/pjheden/elephantio/config"
	"github.com/pjheden/elephantio/logs"
	"github.com/pjheden/elephantio/task"
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
	return sql.Open("sqlite3", d.dbPath)
}

/*
FullTasks... TODO:
*/
func (d *Database) FullTasks() ([]*task.FullTask, error) {
	query := sq.Select(
		"tasks.task",
		"tasks.interval",
		"logs.completed_on",
		"config.pin_1",
		"config.pin_2",
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

	fts := []*task.FullTask{}
	for rows.Next() {
		ft := &task.FullTask{}
		err := rows.Scan(
			&(ft.Task),
			&(ft.Interval),
			&(ft.CompletedOn),
			&(ft.Pin1),
			&(ft.Pin2),
		)
		if err != nil {
			return nil, fmt.Errorf("scanning issues: %v", err)
		}
		fts = append(fts, ft)
	}

	return fts, nil
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
			&(t.Task),
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
		"pin_1",
		"pin_2",
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
			&(c.Pin1),
			&(c.Pin2),
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
