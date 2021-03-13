package config

import "github.com/gofrs/uuid"

type Config struct {
	ID        uuid.UUID
	TaskID    uuid.UUID
	ButtonPin int
	LEDPin    int
}
