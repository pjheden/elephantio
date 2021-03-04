package task

import (
	"time"

	"github.com/gofrs/uuid"
)

/*
FullTask holds all relevant information regarding a task
*/
type FullTask struct {
	Task        string
	Interval    float64
	CompletedOn time.Time
	Pin1        int
	Pin2        int
}

/*
Task is an action that should be done on set interveals defined by Interval
*/
type Task struct {
	ID       uuid.UUID
	Task     string
	Interval float64
}
