package logs

import (
	"time"

	"github.com/gofrs/uuid"
)

/*
Logs is a completed action
*/
type Logs struct {
	ID          uuid.UUID
	TaskID      uuid.UUID
	CompletedOn time.Time
}
