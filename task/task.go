package task

import (
	"github.com/gofrs/uuid"
)

/*
Task is an action that should be done on set interveals defined by Interval
*/
type Task struct {
	ID       uuid.UUID
	Name     string
	Interval float64
}
