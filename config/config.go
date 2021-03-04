package config

import "github.com/gofrs/uuid"

type Config struct {
	ID     uuid.UUID
	TaskID uuid.UUID
	Pin1   int
	Pin2   int
}
